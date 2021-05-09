package plugins

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	"github.com/alex-held/devctl-plugins/pkg/devctlpath"
	"github.com/alex-held/devctl-plugins/pkg/sysutils"
)

var FsMode int

func ParseFsModeFlag() {
	fmt.Printf("DEFAULT FSMode VALUE = %v\n", FsMode)
	flagSet := pflag.NewFlagSet("PluginsCacheFlagSet", pflag.ContinueOnError)
	flagSet.IntVar(&FsMode, "fs-mode", 0, "-fs-mode=1")
	err := flagSet.Parse(os.Args[0:])
	if err != nil {
		fmt.Printf("ERROR PARSING FLAGSET ERR=%v\n", err)
	}
	fmt.Printf("PARSED FSMode VALUE = %v\n", FsMode)
}

func init() {
	ParseFsModeFlag()
}

type cachedPlugin struct {
	Commands Commands `json:"commands"`
	CheckSum string   `json:"check_sum"`
}

type cachedPlugins map[string]cachedPlugin

// CachePath returns the path to the plugins cache
// Defaults to DEVCTL_ROOT/plugins/plugin.cache.
var CachePath = func() string {
	pather := devctlpath.NewPather()
	return pather.Plugins("plugin.cache")
}()

const (
	PLUGIN_CACHE_KEY = "DEVCTL_PLUGIN_CACHE"
	PLUGIN_CACHE_ON  = "on"
)

var cacheMoot sync.RWMutex
var cacheOn = sysutils.GetEnvOrDefault(PLUGIN_CACHE_KEY, PLUGIN_CACHE_ON)

var CacheFs = func() afero.Fs {
	switch FsMode {
	case 0:
		return afero.NewOsFs()
	default:
		return afero.NewMemMapFs()
	}
}

var fs = CacheFs()

var cache = func() cachedPlugins {
	m := cachedPlugins{}
	if cacheOn != PLUGIN_CACHE_ON {
		return m
	}

	f, err := fs.Open(CachePath)
	if err != nil {
		return m
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		f.Close()
		fs.Remove(f.Name())
	}
	return m
}()

func findInCache(path string) (cachedPlugin, bool) {
	cacheMoot.RLock()
	defer cacheMoot.RUnlock()
	cp, ok := cache[path]
	return cp, ok
}

func saveCache() error {
	if cacheOn != PLUGIN_CACHE_ON {
		return nil
	}
	cacheMoot.Lock()
	defer cacheMoot.Unlock()
	_ = fs.MkdirAll(filepath.Dir(CachePath), 0744)
	f, err := fs.Create(CachePath)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(cache)
}

func sum(path string) string {
	f, err := fs.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return ""
	}
	sum := hash.Sum(nil)

	s := fmt.Sprintf("%x", sum)
	return s
}

func addToCache(path string, cp cachedPlugin) {
	if cp.CheckSum == "" {
		cp.CheckSum = sum(path)
	}
	cacheMoot.Lock()
	defer cacheMoot.Unlock()
	cache[path] = cp
}
