package get

import (
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// getOptions is returned by newGetOptions
type getOptions struct {
	ID            string
	LabelSelector string
	Type          int
	PageNumber    int
	PageSize      int
	printers.IOStreams
}

// newGetOptions returns a getOptions with default page size 999.
func newGetOptions(streams printers.IOStreams) *getOptions {
	return &getOptions{
		IOStreams:  streams,
		Type:       0,
		PageNumber: 1,
		PageSize:   999,
	}
}

// NewCmdGet creates a command object for the generic "get" action, which
// retrieves one or more resources from a server.
func NewCmdGet(streams printers.IOStreams, f util.Factory) *cobra.Command {

	o := newGetOptions(streams)

	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
	}

	cmd.PersistentFlags().StringVarP(&o.ID, "id", "i", o.ID, "ID of resource.")

	cmd.AddCommand(newCmdGetClusters(o, f))
	cmd.AddCommand(newCmdGetNodes(o, f))
	cmd.AddCommand(newCmdGetPods(o, f))
	cmd.AddCommand(newCmdGetApps(o, f))
	cmd.AddCommand(newCmdGetAppConfig(o, f))
	cmd.AddCommand(newCmdGetUsers(o, f))

	return cmd
}
