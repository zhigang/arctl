package describe

import (
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// describeOptions is returned by newOptions
type describeOptions struct {
	ID         string
	Name       string
	Type       int
	PageNumber int
	PageSize   int
	printers.IOStreams
}

// newOptions returns a describeOptions with default page size 999.
func newOptions(streams printers.IOStreams) *describeOptions {
	return &describeOptions{
		IOStreams:  streams,
		Type:       0,
		PageNumber: 1,
		PageSize:   999,
	}
}

// NewCmdDescribe returns new initialized instance of describe sub command
func NewCmdDescribe(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newOptions(streams)

	var cmd = &cobra.Command{
		Use:   "describe",
		Short: "describe one or many resources",
	}

	cmd.AddCommand(newCmdDescribeApp(o, f))
	cmd.AddCommand(newCmdDescribeAppConfig(o, f))

	cmd.PersistentFlags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")

	return cmd
}
