package plugins

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

func TestPlugins(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugins Suite")
}

var _ = Describe("Plugin", func() {

	Describe("Cache", func() {

		When("FsMode == 0 ", func() {
			Specify("DEFAULT", func() {})

			BeforeEach(func() {
				fmt.Printf("DEFAULT FSMode VALUE = %v\n", FsMode)
				flagSet := pflag.NewFlagSet("PluginsCacheFlagSet", pflag.ContinueOnError)
				err := flagSet.Parse(os.Args[0:])
				if err != nil {
					fmt.Printf("ERROR PARSING FLAGSET ERR=%v\n", err)
				}
				fmt.Printf("PARSED FSMode VALUE = %v\n", FsMode)
			})

			It("Uses os.Fs", func() {
				var fs = CacheFs()
				Expect(fs).ShouldNot(BeNil())
				Expect(fs).Should(BeAssignableToTypeOf(afero.NewOsFs()))
			})
		})

		When("FsMode > 0", func() {

			BeforeEach(func() {
				ParseFsModeFlag()
			})

			It("uses memory.Fs", func() {
				FsMode = 1
				var fs = CacheFs()
				Expect(fs).ShouldNot(BeNil())
				Expect(fs).Should(BeAssignableToTypeOf(afero.NewMemMapFs()))
			})
		})

	})
})
