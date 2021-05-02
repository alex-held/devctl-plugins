package devctlpath

type Option func(*lazypathFinder) *lazypathFinder

func WithAppPrefix(prefix string) Option {
	return func(l *lazypathFinder) *lazypathFinder {
		l.lp = lazypath(prefix)
		return l
	}
}
func WithConfigFile(cfgName string) Option {
	return func(l *lazypathFinder) *lazypathFinder {
		l.cfgName = cfgName
		return l
	}
}
func WithUserHomeFn(userHomeFn UserHomePathFinder) Option {
	return func(l *lazypathFinder) *lazypathFinder {
		l.finder.GetUserHomeFn = userHomeFn
		return l
	}
}
func WithConfigRootFn(cfgRootFn ConfigRootFinder) Option {
	return func(l *lazypathFinder) *lazypathFinder {
		l.finder.GetConfigRootFn = cfgRootFn
		return l
	}
}
func WithCachePathFn(cacheFn CachePathFinder) Option {
	return func(l *lazypathFinder) *lazypathFinder {
		l.finder.GetCachePathFn = cacheFn
		return l
	}
}

func defaults() []Option {
	return []Option{WithAppPrefix("devctl"),
		WithCachePathFn(nil),
		WithUserHomeFn(nil),
		WithConfigRootFn(nil),
		WithConfigFile(devctlConfigFileName),
	}
}
