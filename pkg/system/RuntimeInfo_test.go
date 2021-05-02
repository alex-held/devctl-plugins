package system_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alex-held/devctl-plugins/pkg/system"
)

func TestRuntimeInfoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RuntimeInfo Suite")
}

var _ = Describe("RuntimeInfo", func() {

	Describe("Format", func() {

		const os = "darwin"
		const arch = "amd64"
		var sut system.RuntimeInfo

		BeforeEach(func() {
			sut = system.RuntimeInfo{
				OS:   os,
				Arch: arch,
			}
		})

		Context("pattern contains just runtime info templates ", func() {
			When("pattern contains [os]", func() {
				It("[os] gets replaced", func() {
					Expect(sut.Format("/some/filename.1.32.3[os]")).Should(Equal("/some/filename.1.32.3" + os))
				})
			})
		})

	})
})
