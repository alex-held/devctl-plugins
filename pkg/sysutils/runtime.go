package sysutils

import (
	"fmt"
	"runtime"
	"strings"
)

type RuntimeInfo struct {
	OS, Arch string
}

// Format formats pattern string using the provided args and the runtime info
// [os] = OS (darwin, linux)
// [OS] = OS (DARWIN, LINUX)
// [arch] = amd64, arm64
// [ARCH] = AMD64, ARM64.
func (info RuntimeInfo) Format(pattern string, args ...interface{}) string {
	formatted := fmt.Sprintf(pattern, args...)
	formatted = strings.ReplaceAll(formatted, "[os]", strings.ToLower(info.OS))
	formatted = strings.ReplaceAll(formatted, "[OS]", strings.ToUpper(info.OS))
	formatted = strings.ReplaceAll(formatted, "[arch]", strings.ToLower(info.Arch))
	formatted = strings.ReplaceAll(formatted, "[ARCH]", strings.ToUpper(info.Arch))

	return formatted
}

func (info RuntimeInfo) IsLinux() bool   { return IsLinux(info) }
func IsLinux(info RuntimeInfo) bool      { return strings.ContainsAny(info.OS, "linux") }
func (info RuntimeInfo) IsWindows() bool { return IsWindows(info) }
func IsWindows(info RuntimeInfo) bool    { return strings.ContainsAny(info.OS, "windows") }
func (info RuntimeInfo) IsDarwin() bool  { return IsDarwin(info) }
func IsDarwin(info RuntimeInfo) bool     { return strings.ContainsAny(info.OS, "darwin") }

type OSRuntimeInfoGetter struct{}
type DefaultRuntimeInfoGetter struct {
	GOOS, GOARCH string
}
type RuntimeInfoGetter interface {
	Get() (info RuntimeInfo)
}

// Format formats pattern string using the provided args and the runtime info
// [os] = OS (darwin, linux)
// [OS] = OS (DARWIN, LINUX)
// [arch] = amd64, arm64
// [ARCH] = AMD64, ARM64.
func (d *DefaultRuntimeInfoGetter) Format(pattern string, args ...interface{}) string {
	return d.Get().Format(pattern, args...)
}

func (d *DefaultRuntimeInfoGetter) Get() RuntimeInfo {
	if d == nil {
		return OSRuntimeInfoGetter{}.Get()
	}

	return RuntimeInfo{
		OS:   d.GOOS,
		Arch: d.GOARCH,
	}
}

// Format formats pattern string using the provided args and the runtime info
// [os] = OS (darwin, linux)
// [OS] = OS (DARWIN, LINUX)
// [arch] = amd64, arm64
// [ARCH] = AMD64, ARM64.
func (g *OSRuntimeInfoGetter) Format(pattern string, args ...interface{}) string {
	info := g.Get()
	return info.Format(pattern, args...)
}

func (OSRuntimeInfoGetter) Get() (info RuntimeInfo) {
	if archID := runtime.GOARCH; archID == "arm" {
		return RuntimeInfo{
			OS:   runtime.GOOS,
			Arch: "arm64",
		}
	}

	return RuntimeInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
