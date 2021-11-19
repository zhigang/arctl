package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdDeleteLabel is a command to delete labels from a cluster.
func newCmdDeleteLabel(o *deleteOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "label",
		Short: "Delete labels on a .",
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
				resp, err := ns.DeleteNodeLabel(o.ID, k, v, o.Force)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				baseInfo := fmt.Sprintf("delete label %s=%s on cluster %s", k, v, o.ID)
				if resp.Code == 0 && resp.Success {
					printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
				} else {
					printers.PrintValue(w, fmt.Sprintf("%s failed, error msg: %s\n", baseInfo, resp.ErrMsg))
				}
				w.Flush()
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&o.LabelSelector, "label", "l", o.LabelSelector, "Labels for instance, supports '='.(e.g. -l key1=value1,key2=value2)")
	cmd.MarkFlagRequired("label")

	return cmd
}
