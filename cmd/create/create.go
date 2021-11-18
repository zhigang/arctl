package create

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCreate(streams printers.IOStreams, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create one or many resources",
	}

	cmd.AddCommand(newCmdCreateApp(streams, f))
	cmd.AddCommand(newCmdCreateEnv(streams, f))
	cmd.AddCommand(newCmdCreateAppConfig(streams, f))

	return cmd
}

func readFileBase64(path string) (string, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	enc := base64.StdEncoding.EncodeToString(c)
	return enc, nil
}
