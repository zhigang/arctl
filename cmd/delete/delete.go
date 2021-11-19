package delete

import (
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// deleteOptions is returned by newDeleteOptions
type deleteOptions struct {
	ID            string
	Force         bool
	LabelSelector string
	PageNumber    int
	PageSize      int
	printers.IOStreams
}

// newDeleteOptions returns an initialized deleteOptions instance
func newDeleteOptions(streams printers.IOStreams) *deleteOptions {
	return &deleteOptions{
		IOStreams:  streams,
		PageNumber: 1,
		PageSize:   999,
	}
}

// NewCmdDelete returns new initialized instance of delete sub command
func NewCmdDelete(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newDeleteOptions(streams)

	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a resources",
	}

	cmd.PersistentFlags().StringVarP(&o.ID, "id", "i", o.ID, "Instance ID of resource. (required)")
	cmd.PersistentFlags().BoolVarP(&o.Force, "force", "f", o.Force, "Force delete resource.")
	cmd.MarkPersistentFlagRequired("id")

	cmd.AddCommand(newCmdDeleteApp(o, f))
	cmd.AddCommand(newCmdDeleteAppEnv(o, f))
	cmd.AddCommand(newCmdDeleteAppConfig(o, f))
	cmd.AddCommand(newCmdDeleteLabel(o, f))
	return cmd
}
