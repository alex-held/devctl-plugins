// Package devctlpath calculates filesystem paths to devctl's configuration, cache and data.
package devctlpath

const devctlConfigFileName = "config.yaml"
const PluginConfigFileName = "plugins.toml"

type finder struct {
	GetUserHomeFn   UserHomePathFinder
	GetCachePathFn  CachePathFinder
	GetConfigRootFn ConfigRootFinder
}

func (f *finder) UserHomePathFinder() string { return f.GetUserHomeFn() }
func (f *finder) CachePath() string          { return f.GetCachePathFn() }
func (f *finder) ConfigRoot() string         { return f.GetConfigRootFn() }

type UserHomePathFinder func() string
type CachePathFinder func() string
type ConfigRootFinder func() string

// ConfigRoot the path where Helm stores configuration.
func (f *lazypathFinder) ConfigRoot(elem ...string) string { return f.configRoot(elem...) }

// ConfigFilePath  path where Helm stores configuration.
func (f *lazypathFinder) ConfigFilePath() string {
	return f.configRoot(f.cfgName)
}

// PluginConfigFilePath resolves the plugin configuration
// The default is $DEVCTL_ROOT/plugins.toml.
func (f *lazypathFinder) PluginConfigFilePath() string {
	return f.configRoot(PluginConfigFileName)
}

// Config returns the path where various application configurations are stored.
func (f *lazypathFinder) Config(elem ...string) string { return f.resolveSubDir("config", elem...) }

// Bin returns the path where executable are stored.
func (f *lazypathFinder) Bin(elem ...string) string { return f.resolveSubDir("bin", elem...) }

// Download returns the path where downloads are stored.
func (f *lazypathFinder) Download(elem ...string) string {
	return f.resolveSubDir("downloads", elem...)
}

// Plugins the path where plugins are stored.
func (f *lazypathFinder) Plugins(elem ...string) string { return f.resolveSubDir("plugins", elem...) }

// SDK returns the path where sdk installations are stored & managed.
func (f *lazypathFinder) SDK(elem ...string) string { return f.resolveSubDir("sdks", elem...) }

// Cache returns the path where transient information are cached.
func (f *lazypathFinder) Cache(elem ...string) string { return f.cachePath(elem...) }
