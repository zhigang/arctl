package get

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetUsers creates a command object for get users.
func newCmdGetUsers(o *getOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "users",
		Short: "List users",
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()

			cs := f.UserService()

			resp, err := cs.GetUserList(o.PageNumber, o.PageSize)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			if len(resp.Data) == 0 {
				util.ExitWithMsg(o.ErrOut, "users not found")
			}
			headers := []string{"User ID", "User Name", "User Type"}
			for _, c := range resp.Data {
				if o.ID == "" {
					w.Append([]string{c.UserId, c.RealName, c.UserType})
				} else if o.ID == c.UserId {
					w.Append([]string{c.UserId, c.RealName, c.UserType})
					break
				}
			}
			w.SetHeader(headers)
			footer := make([]string, len(headers))
			footer[len(footer)-2] = "Total"
			footer[len(footer)-1] = strconv.Itoa(w.NumLines())
			w.SetFooter(footer)
		},
	}
	return cmd
}
