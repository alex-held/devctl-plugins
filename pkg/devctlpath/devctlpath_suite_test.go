package devctlpath_test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl-plugins/pkg/devctlpath"
)

func TestDevCtlPatherSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pather Suite")
}

const (
	prefix    = ".devctl"
	goosLinux = "linux"
)

var _ = Describe("Pather", func() {

	var sut devctlpath.Pather
	var osDefault, expected string

	BeforeEach(func() {
		sut = devctlpath.NewPather()
	})

	Context("with default value", func() {

		Context("ConfigRoot", func() {
			Specify("Linux -> attempts to get XDG_CACHE_HOME/.devctl", func() {})
			Specify("Darwin -> ~/.devctl", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				sut = devctlpath.NewPather()
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix)
			})

			It("ConfigRoot returns os default ConfigRoot", func() {
				Expect(sut.ConfigRoot()).Should(Equal(expected))
			})

			It("ConfigRoot has default '.devctl' prefix", func() {
				Expect(sut.ConfigRoot()).Should(HaveSuffix(prefix))
			})
		})

		Context("Cache", func() {
			Specify("Linux -> attempts to get XDG_CACHE_HOME/io.alexheld.devctl", func() {})
			Specify("Darwin -> ~/Library/Caches/io.alexheld.devctl", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CACHE_HOME", "/my/cache")
				}
				osDefault, _ = os.UserCacheDir()
				expected = filepath.Join(osDefault, "io.alexheld.devctl")
			})

			It("returns os default cache directory", func() {
				Expect(sut.Cache()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.Cache("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})

		Context("Bin()", func() {
			Specify("Linux -> attempts to get XDG_CONFIG_HOME/.devctl/bin", func() {})
			Specify("Darwin -> ~/.devctl/bin", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix, "bin")
			})

			It("returns os default devctl bin directory", func() {
				Expect(sut.Bin()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.Bin("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})

		Context("SDK()", func() {
			Specify("Linux -> attempts to get XDG_CONFIG_HOME/.devctl/sdks", func() {})
			Specify("Darwin -> ~/.devctl/sdks", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix, "sdks")
			})

			It("returns os default sdks directory", func() {
				Expect(sut.SDK()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.SDK("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})

		Context("Config()", func() {
			Specify("Linux -> attempts to get XDG_CONFIG_HOME/devctl/config", func() {})
			Specify("Darwin -> ~/.devctl/config", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix, "config")
			})

			It("returns os default config directory", func() {
				Expect(sut.Config()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.Config("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})

		Context("Download()", func() {
			Specify("Linux -> attempts to get XDG_CONFIG_HOME/devctl/downloads", func() {})
			Specify("Darwin -> ~/.devctl/downloads", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix, "downloads")
			})

			It("returns os default config directory", func() {
				Expect(sut.Download()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.Download("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})

		Context("Plugins()", func() {
			Specify("Linux -> attempts to get XDG_CONFIG_HOME/.devctl/plugins", func() {})
			Specify("Darwin -> ~/.devctl/plugins", func() {})

			BeforeEach(func() {
				if runtime.GOOS == goosLinux {
					// emulate freedesktop environment variables
					_ = os.Setenv("XDG_CONFIG_HOME", "/my/home")
				}
				osDefault, _ = os.UserHomeDir()
				expected = path.Join(osDefault, prefix, "plugins")
			})

			It("returns os default plugins directory", func() {
				Expect(sut.Plugins()).To(Equal(expected))
			})

			When("providing additional path elements", func() {
				It("creates a valid subdirectory path", func() {
					Expect(sut.Plugins("some", "sub", "dir")).Should(Equal(path.Join(expected, "some", "sub", "dir")))
				})
			})
		})
	})
})
