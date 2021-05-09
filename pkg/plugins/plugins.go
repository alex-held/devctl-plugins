package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/hashicorp/go-hclog"
	"github.com/karrick/godirwalk"
	"github.com/markbates/errx"
	"github.com/markbates/oncer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/alex-held/devctl-plugins/pkg/devctlog"
	"github.com/alex-held/devctl-plugins/pkg/devctlpath"
	"github.com/alex-held/devctl-plugins/pkg/plugins/deps"
	"github.com/alex-held/devctl-plugins/pkg/sysutils"
)

const timeoutEnv = "DEVCTL_PLUGIN_TIMEOUT"

var t = time.Second * 2

func timeout() time.Duration {
	oncer.Do("plugins.timeout", func() {
		rawTimeout, err := envy.MustGet(timeoutEnv)
		if err == nil {
			if parsed, err := time.ParseDuration(rawTimeout); err == nil {
				t = parsed
			} else {
				logrus.Errorf("%q value is malformed assuming default %q: %v", timeoutEnv, t, err)
			}
		} else {
			logrus.Debugf("%q not set, assuming default of %v", timeoutEnv, t)
		}
	})
	return t
}

// List is a map of  devctl command to a slice of Command.
type List map[string]Commands

var _list List

// Available plugins for the `devctl` command.
// It will look in $GOPATH/bin and the `./plugins` directory.
// This can be changed by setting the $DEVCTL_PLUGIN_ROOT
// environment variable.
//
// Requirements:
// * file/command must be executable
// * file/command must match regex `devctl-\w+?-plugin-.*`
// * file/command must respond to `available` and return JSON of
//	 plugins.Commands{}
//
// Limit full path scan with direct plugin path
//
// If a file/command doesn't respond to being invoked with `available`
// within one second, buffalo will assume that it is unable to load. This
// can be changed by setting the $DEVCTL_PLUGIN_TIMEOUT environment
// variable. It must be set to a duration that `time.ParseDuration` can
// process.
func Available() (List, error) {
	var err error
	log := devctlog.New("plugins-loader")
	oncer.Do("plugins.Available", func() {
		defer func() {
			if err := saveCache(); err != nil {
				log.Error(err.Error())
			}
		}()

		app := New(Options{
			Name:     "devctl",
			Env:      "development",
			LogLevel: hclog.Debug,
			Prefix:   ".devctl",
		})

		if exist, _ := afero.Exists(app.Context.GetFs(), app.Context.GetPather().PluginConfigFilePath()); exist {
			app.Logger.Info("PluginConfigFilePath exists!")
			_list, err = listPlugDeps(app)
			return
		}

		paths := []string{"plugins"}

		from, err := envy.MustGet("DEVCTL_PLUGIN_ROOT")
		if err != nil {
			from, err = envy.MustGet("GOPATH")
			if err != nil {
				return
			}
			from = filepath.Join(from, "bin")
		}

		paths = append(paths, strings.Split(from, string(os.PathListSeparator))...)

		list := List{}
		for _, p := range paths {
			// todo: find out usecase
			// if ignorePath(p) {
			// 	continue
			// }
			if _, err := os.Stat(p); err != nil {
				continue
			}

			err := godirwalk.Walk(p, &godirwalk.Options{
				FollowSymbolicLinks: true,
				Callback: func(path string, info *godirwalk.Dirent) error {
					if err != nil {
						// May indicate a permissions problem with the path, skip it
						return nil
					}
					if info.IsDir() {
						return nil
					}
					base := filepath.Base(path)
					if strings.HasPrefix(base, "devctl-") && !strings.HasPrefix(base, "buffalo-plugins") {
						// todo: app.HostContext probably does not make sense...
						_, cancel := context.WithTimeout(app.Context, t)
						commands := askBin(app.Context, path)
						cancel()
						for _, c := range commands {
							bc := c.DevctlCommand
							if _, ok := list[bc]; !ok {
								list[bc] = Commands{}
							}
							c.Binary = path
							list[bc] = append(list[bc], c)
						}
					}
					return nil
				},
			})

			if err != nil {
				return
			}
		}
		_list = list
	})
	return _list, err
}

func askBin(ctx HostContext, path string) Commands {
	start := time.Now()
	defer func() {
		ctx.GetLogger().Debug("askBin %s=%.4f s", path, time.Since(start).Seconds())
	}()

	commands := Commands{}
	if cp, ok := findInCache(path); ok {
		s := sum(path)
		if s == cp.CheckSum {
			ctx.GetLogger().Debug("cache hit: %s", path)
			commands = cp.Commands
			return commands
		}
	}
	ctx.GetLogger().Debug("cache miss: %s", path)
	if strings.HasPrefix(filepath.Base(path), "devctl-no-sqlite") {
		return commands
	}

	cmd := exec.CommandContext(ctx, path, "available")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb
	err := cmd.Run()
	if err != nil {
		ctx.GetLogger().Error("[PLUGIN] error loading plugin %s: %s\n", path, err)
		return commands
	}

	msg := bb.String()
	for len(msg) > 0 {
		err = json.NewDecoder(strings.NewReader(msg)).Decode(&commands)
		if err == nil {
			addToCache(path, cachedPlugin{
				Commands: commands,
			})
			return commands
		}
		msg = msg[1:]
	}
	ctx.GetLogger().Error("[PLUGIN] error decoding plugin %s: %s\n%s\n", path, err, msg)
	return commands
}

func ignorePath(p string) bool {
	p = strings.ToLower(p)
	for _, x := range []string{`c:\windows`, `c:\program`} {
		if strings.HasPrefix(p, x) {
			return true
		}
	}
	return false
}

func listPlugDeps(app *App) (List, error) {
	list := List{}
	plugs, err := deps.List(app.Context.GetLogger(), app.Context.GetPather())

	if err != nil {
		return list, err
	}

	for _, p := range plugs.List() {
		_, cancel := context.WithTimeout(context.Background(), timeout())
		defer cancel()
		bin := p.Binary
		if len(p.Local) != 0 {
			bin = p.Local
		}
		bin, err := LookPath(bin)
		if err != nil {
			if errx.Unwrap(err) != ErrPluginMissing {
				return list, err
			}
			continue
		}
		commands := askBin(app.Context, bin)
		cancel()
		for _, c := range commands {
			bc := c.DevctlCommand
			if _, ok := list[bc]; !ok {
				list[bc] = Commands{}
			}
			c.Binary = p.Binary
			for _, pc := range p.Commands {
				if c.Name == pc.Name {
					c.Flags = pc.Flags
					break
				}
			}
			list[bc] = append(list[bc], c)
		}
	}
	return list, nil
}

// Options are used to configure and define how your application should run.
type Options struct {
	// The name and id of plugin host app
	Name string `json:"name"`

	// Env is the "environment" in which the App is running. Default is "development".
	Env string `json:"env"`

	// LogLevel defaults to hclog.Level (INFO).
	LogLevel hclog.Level `json:"log_level"`

	// Logger to be used with the application.
	// A default one is provided.
	Logger devctlog.Logger `json:"-"`

	Prefix  string      `json:"prefix"`
	Context HostContext `json:"-"`

	cancel context.CancelFunc
}

// HostContext holds on to information as you
// pass it down through middleware, Handlers,
// templates, etc... It strives to make your
// life a happier one.
type HostContext interface {
	context.Context
	GetRuntimeInfo() sysutils.RuntimeInfo
	GetLogger() devctlog.Logger
	GetFs() afero.Fs
	GetPather() devctlpath.Pather
	DataStore() map[string]interface{}
}

// assert that DefaultContext is implementing HostContext.
var _ HostContext = &DefaultContext{}
var _ context.Context = &DefaultContext{}

type DefaultContext struct {
	context.Context
	pather            devctlpath.Pather
	runtimeInfoGetter *sysutils.DefaultRuntimeInfoGetter
	logger            devctlog.Logger
	datastore         *sync.Map
	fs                afero.Fs
}

func (d *DefaultContext) GetFs() afero.Fs                      { return d.fs }
func (d *DefaultContext) GetRuntimeInfo() sysutils.RuntimeInfo { return d.runtimeInfoGetter.Get() }
func (d *DefaultContext) GetLogger() devctlog.Logger           { return d.logger }
func (d *DefaultContext) GetPather() devctlpath.Pather         { return d.pather }

func New(options Options) *App {
	app := &App{
		Options: options,
	}
	defaultCtx := &DefaultContext{
		pather:    devctlpath.NewPather(devctlpath.WithAppPrefix(options.Prefix)),
		logger:    app.Logger,
		fs:        afero.NewOsFs(),
		datastore: &sync.Map{},
	}
	defaultCtx.Context, app.cancel = context.WithCancel(context.Background())
	app.Context = defaultCtx

	return app
}

func (d *DefaultContext) String() string {
	data := d.DataStore()
	bb := make([]string, 0, len(data))

	for k, v := range data {
		bb = append(bb, fmt.Sprintf("%s: %s", k, v))
	}
	sort.Strings(bb)
	return strings.Join(bb, "\n\n")
}

// DataStore contains all the values set through Get/Set.
func (d *DefaultContext) DataStore() map[string]interface{} {
	m := map[string]interface{}{}
	d.datastore.Range(func(k, v interface{}) bool {
		s, ok := k.(string)
		if !ok {
			return false
		}
		m[s] = v
		return true
	})
	return m
}

// MarshalJSON implements json marshaling for the context.
func (d *DefaultContext) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	data := d.DataStore()
	for k, v := range data {
		// don't try and marshal ourself
		if _, ok := v.(*DefaultContext); ok {
			continue
		}
		if _, err := json.Marshal(v); err == nil {
			// it can be marshaled, so add it:
			m[k] = v
		}
	}

	d.GetLogger().Trace("staring to marshall DefaultContext", "context", *d, "datastore", data)
	bytes, err := json.Marshal(m)
	if err != nil {
		d.GetLogger().Error("failed marshaling of DefaultContext", "context", *d, "err", err)
		return bytes, err
	}
	jsonContext := string(bytes)
	d.GetLogger().Debug("marshaled DefaultContext", "context", *d, "json", jsonContext)
	return bytes, err
}
