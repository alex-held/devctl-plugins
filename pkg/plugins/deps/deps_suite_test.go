package deps_test

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"testing"

	"github.com/gobuffalo/meta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl-plugins/pkg/devctlpath"
	"github.com/alex-held/devctl-plugins/pkg/plugins/deps"
)

//go:embed testdata/with-multiple-plugins.golden.toml
var golden_with_multiple_plugins []byte

//go:embed testdata/with-go-sdk-plugin.golden.toml
var golden_with_go_sdk_plugin []byte

func TestDeps(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deps Suite")
}

var _ = Describe("PlugDeps", func() {

	Describe("Plugins", func() {

		Context("with default value", func() {
			var sut *deps.Plugins
			var expected string
			var actualBuf *bytes.Buffer

			BeforeEach(func() {
				actualBuf = &bytes.Buffer{}
				sut = deps.New()
			})

			It("should not error", func() {
				Expect(sut.Encode(actualBuf)).Should(Succeed())
			})

			It("encodes into empty string", func() {
				_ = sut.Encode(actualBuf)
				Expect(actualBuf.String()).Should(Equal(expected))
			})
		})

		Context("with go-sdk plugin", func() {
			var sut *deps.Plugins
			var expected string
			var actualBuf *bytes.Buffer

			BeforeEach(func() {
				actualBuf = &bytes.Buffer{}
				sut = deps.New()
				expected = string(golden_with_go_sdk_plugin)
				pather := devctlpath.NewPather()
				sut.Add(deps.Plugin{
					Binary: "devctl-go-sdk-plugin",
					GoGet:  "go get -u github.com/alex-held/devctl-go-sdk-plugin",
					Local:  "." + strings.TrimPrefix(pather.Plugins("sdk", "devctl-go-sdk-plugin"), pather.ConfigRoot()),
					Commands: []deps.Command{
						{Name: "download"},
						{Name: "list"},
						{Name: "install"},
						{Name: "link"},
						{Name: "use"},
					},
					Tags: meta.BuildTags{"devctl-plugin", "sdk", "v1.0.0"},
				})
			})

			It("should not error", func() {
				Expect(sut.Encode(actualBuf)).Should(Succeed())
			})

			It("encodes into expected toml file", func() {
				_ = sut.Encode(actualBuf)
				actualToml := actualBuf.String()
				Expect(actualToml).Should(Equal(expected))
				Expect(actualBuf.Bytes()).Should(HaveLen(412))
			})
		})
	})

	Context("with multiple plugins", func() {
		var sut *deps.Plugins
		var expected *deps.Plugins
		var fileReader io.Reader

		BeforeEach(func() {
			fileReader = strings.NewReader(string(golden_with_multiple_plugins))
			sut = deps.New()
			expected = deps.New()

			expected.Add(deps.Plugin{
				Binary: "buffalo-heroku",
				GoGet:  "github.com/gobuffalo/buffalo-heroku",
			})
			expected.Add(deps.Plugin{
				Binary: "buffalo-pop",
				GoGet:  "github.com/gobuffalo/buffalo-pop",
			})
			expected.Add(deps.Plugin{
				Binary: "buffalo-trash",
				GoGet:  "github.com/markbates/buffalo-trash",
			})
		})

		It("should not error", func() {
			Expect(sut.Decode(fileReader)).Should(Succeed())
		})

		It("decodes from toml file into expected values", func() {
			_ = sut.Decode(fileReader)
			Expect(sut).Should(Equal(expected))
			Expect(sut.List()).Should(HaveLen(3))
			Expect(sut.List()[0].Binary).Should(Equal("buffalo-heroku"))
			Expect(sut.List()[1].Binary).Should(Equal("buffalo-pop"))
			Expect(sut.List()[2].GoGet).Should(Equal("github.com/markbates/buffalo-trash"))
		})
	})
})
