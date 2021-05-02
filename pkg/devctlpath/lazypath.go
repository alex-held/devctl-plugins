package devctlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl-plugins/pkg/system"
)

type lazypath string
type lazypathFinder struct {
	cfgName string
	lp      lazypath
	finder  finder
}

//goland:noinspection GoUnusedConst
const (
	DevctlRootKey       = "DEVCTL_ROOT"
	DevctlConfigHomeKey = "DEVCTL_CONFIG_HOME"
	DevctlCacheHomeKey  = "DEVCTL_CACHE_HOME_KEY"
)

// Pather resolves different paths related to the CLI itself.
type Pather interface {

	// ConfigFilePath returns the path of the config.yaml used to configure the app itself
	ConfigFilePath() string

	// ConfigRoot returns the root path of the CLI configuration.
	ConfigRoot(elem ...string) string

	// Config returns a path to store configuration.
	Config(elem ...string) string

	// Bin returns a path to store executable binaries.
	Bin(elem ...string) string

	// Download returns the path where to save downloads
	Download(elem ...string) string

	// SDK returns the path where sdks are installed
	SDK(elem ...string) string

	// Cache returns the path where to cache files
	Cache(elem ...string) string

	// Plugins returns the path where plugins are stored
	Plugins(elem ...string) string
}

// NewPather creates and configures a Pather using the default Option's and then applies provided opts Option's.
func NewPather(opts ...Option) Pather {
	lpFinder := &lazypathFinder{}
	for _, opt := range defaults() {
		lpFinder = opt(lpFinder)
	}

	for _, opt := range opts {
		lpFinder = opt(lpFinder)
	}

	return lpFinder
}

func (f lazypathFinder) resolveSubDir(sub string, elem ...string) string {
	subConfig := f.configRoot(sub)
	return filepath.Join(subConfig, filepath.Join(elem...))
}

// There is an order to checking for a path.
// 1. GetConfigRootFn has been provided
// 1. GetUserHomeFn + AppPrefix has been provided
// 2. See if a devctl specific environment variable has been set.
// 2. Check if an XDG environment variable is set
// 3. Fall back to a default.
func (f lazypathFinder) configRoot(elem ...string) string {
	runInfo := system.OSRuntimeInfoGetter{}.Get()

	if f.finder.GetConfigRootFn != nil {
		p := f.finder.ConfigRoot()
		return filepath.Join(p, filepath.Join(elem...))
	}

	if f.finder.GetUserHomeFn != nil {
		p := f.finder.GetUserHomeFn()

		switch {
		case runInfo.IsLinux():
			p = filepath.Join(p, ".config", f.lp.getAppPrefix())
		case runInfo.IsDarwin():
			p = filepath.Join(p, f.lp.getAppPrefix())
		default:
			p = filepath.Join(p, f.lp.getAppPrefix())
		}

		return filepath.Join(p, filepath.Join(elem...))
	}

	base := os.Getenv(DevctlRootKey)
	if base != "" {
		return filepath.Join(base, filepath.Join(elem...))
	}

	if base != "" {
		confRoot := filepath.Join(base, f.lp.getAppPrefix())
		return filepath.Join(confRoot, filepath.Join(elem...))
	}

	base = configHome()(f.lp)

	return filepath.Join(base, filepath.Join(elem...))
}

// cachePath resolves the path where devctl will cache data
// There is an order to checking for a path.
// 1. GetCachePathFn has been provided
// 2. See if a devctl specific environment variable has been set.
// 2. Check if an XDG environment variable is set
// 3. Fall back to a default.
func (f lazypathFinder) cachePath(elem ...string) string {
	fqrdn := fmt.Sprintf("io.alexheld%s", f.lp.getAppPrefix())

	if f.finder.GetCachePathFn != nil {
		p := f.finder.CachePath()
		p = filepath.Join(p, fqrdn)

		return filepath.Join(p, filepath.Join(elem...))
	}

	p := os.Getenv(DevctlCacheHomeKey)
	if p != "" {
		p = filepath.Join(p, fqrdn)
		return filepath.Join(p, filepath.Join(elem...))
	}

	p = os.Getenv(DevctlCacheHomeKey)
	if p != "" {
		p = filepath.Join(p, fqrdn)
		return filepath.Join(p, filepath.Join(elem...))
	}

	p = cacheHome()
	p = filepath.Join(p, fqrdn)

	return filepath.Join(p, filepath.Join(elem...))
}

func (l lazypath) getAppPrefix() (prefix string) {
	prefix = strings.ToLower(fmt.Sprintf(".%s", strings.TrimPrefix(string(l), ".")))
	return prefix
}
