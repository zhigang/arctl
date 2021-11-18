package delete

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdDeleteAppConfig is a command to delete application's config.
func newCmdDeleteAppConfig(o *DeleteOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "config",
		Short: "Delete a application's deploy config.",
		Run: func(cmd *cobra.Command, args []string) {

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			id, err := strconv.Atoi(o.ID)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			as := f.AppService()

			resp, err := as.DeleteDeployConfig(id)

			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			baseInfo := fmt.Sprintf("delete application's deploy config %s", o.ID)
			if resp.Code == 0 && resp.Result.Success {
				printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}

	return cmd
}
