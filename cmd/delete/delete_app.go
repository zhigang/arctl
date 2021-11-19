package delete

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdDeleteApp is a command to delete application.
func newCmdDeleteApp(o *deleteOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "app",
		Short: "Delete a application.",
		Run: func(cmd *cobra.Command, args []string) {

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			id, err := strconv.Atoi(o.ID)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			es := f.AppService()

			resp, err := es.DeleteApp(id, o.Force)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			baseInfo := fmt.Sprintf("delete application %s", o.ID)
			if resp.Code == 0 && resp.Result.Success {
				printers.PrintValue(w, fmt.Sprintf("%s succeed\n", baseInfo))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}

	return cmd
}
