package delete

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdDeleteAppEnv is a command to delete application's environment.
func newCmdDeleteAppEnv(o *DeleteOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "env",
		Short: "Delete a application's environment.",
		Run: func(cmd *cobra.Command, args []string) {

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			id, err := strconv.Atoi(o.ID)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			appID, err := cmd.Flags().GetInt("appID")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			es := f.EnvironmentService()

			resp, err := es.DeleteEnv(appID, id, o.Force)

			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			baseInfo := fmt.Sprintf("delete application's environment %s", o.ID)
			if resp.Code == 0 && resp.Result.Success {
				printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}

	cmd.PersistentFlags().Int("appID", 0, "ID of application.")
	cmd.MarkFlagRequired("appID")

	return cmd
}
