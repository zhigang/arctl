package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// NewCmdVersion returns new initialized instance of version command
func NewCmdVersion(streams printers.IOStreams, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of arctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("arctl v0.1")
		},
	}

	return cmd
}
