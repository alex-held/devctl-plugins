package cmds_test

import (
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl-plugins/pkg/plugins/cmds"
)

func TestCmds(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmds Suite")
}

var _ = Describe("Available", func() {

	Context("Encode", func() {
		const exp = `[{"name":"foo","use_command":"foo","buffalo_command":"generate","description":"generates foo","aliases":["f"]}]`
		bb := &bytes.Buffer{}
		a := cmds.NewAvailable()
		Expect(a.Add("generate", &cobra.Command{
			Use:     "foo",
			Short:   "generates foo",
			Aliases: []string{"f"},
		})).Should(Succeed())
		Expect(a.Encode(bb)).Should(Succeed())
		//		Expect(strings.TrimSpace(bb.String())).Should(Equal(exp))
	})
})
