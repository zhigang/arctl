package label

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// LabelOptions is returned by newLabelOptions
type LabelOptions struct {
	ClusterID     string
	InstanceIDs   []string
	LabelSelector string
	Add           bool
	Remove        bool
	printers.IOStreams
}

// newLabelOptions returns an initialized LabelOptions instance
func newLabelOptions(streams printers.IOStreams) *LabelOptions {
	return &LabelOptions{
		IOStreams: streams,
	}
}

// NewCmdLabel returns new initialized instance of label command
func NewCmdLabel(streams printers.IOStreams, f util.Factory) *cobra.Command {

	o := newLabelOptions(streams)

	var cmd = &cobra.Command{
		Use:   "label",
		Short: "label nodes",
		Run: func(cmd *cobra.Command, args []string) {

			ns := f.NodeService()

			if o.LabelSelector == "" {
				err := fmt.Errorf("the --label flag is required")
				util.CheckErr(o.ErrOut, err)
			}

			labels, _, err := util.ParseLabels(o.LabelSelector)

			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			w := printers.GetNewTabWriter(o.Out)
			for k, v := range labels {
				for _, i := range o.InstanceIDs {
					if o.Add {
						resp, err := ns.BindNodeLabel(o.ClusterID, i, k, v)
						if err != nil {
							util.CheckErr(o.ErrOut, err)
						}
						baseInfo := fmt.Sprintf("add label %s=%s on cluster %s instance %s", k, v, o.ClusterID, i)
						if resp.Code == 0 && resp.Success {
							printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
						} else {
							printers.PrintValue(w, fmt.Sprintf("%s failed, error msg: %s\n", baseInfo, resp.ErrMsg))
						}
					} else if o.Remove {
						resp, err := ns.UnbindNodeLabel(o.ClusterID, i, k, v)
						if err != nil {
							util.CheckErr(o.ErrOut, err)
						}
						baseInfo := fmt.Sprintf("remove label %s=%s on cluster %s instance %s", k, v, o.ClusterID, i)
						if resp.Code == 0 && resp.Success {
							printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
						} else {
							printers.PrintValue(w, fmt.Sprintf("%s failed, error msg: %s\n", baseInfo, resp.ErrMsg))
						}
					}
					w.Flush()
				}
			}
		},
	}

	cmd.Flags().StringSliceVarP(&o.InstanceIDs, "instance", "i", o.InstanceIDs, "ID of instances. (required)")
	cmd.Flags().StringVarP(&o.ClusterID, "cluster", "c", o.ClusterID, "ID of Cluster. (required)")
	cmd.Flags().BoolVarP(&o.Add, "add", "a", false, "Add labels on a instance.")
	cmd.Flags().BoolVarP(&o.Remove, "remove", "r", false, "Remove labels on a instance.")
	cmd.Flags().StringVarP(&o.LabelSelector, "label", "l", o.LabelSelector, "Labels for instance, supports '='.(e.g. -l key1=value1,key2=value2)")

	cmd.MarkFlagRequired("instance")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("label")

	return cmd
}
