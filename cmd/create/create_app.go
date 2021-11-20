package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// CreateAppOptions is returned by newCreateAppOptions
type createAppOptions struct {
	Title       string
	BizCode     string
	Language    string
	OS          string
	ServiceType string
	Owner       string
	BizTitle    string
	Description string
	StateType   int
	Namespace   string
	printers.IOStreams
}

// newCreateAppOptions returns an initialized createAppOptions instance
func newCreateAppOptions(streams printers.IOStreams) *createAppOptions {
	return &createAppOptions{
		IOStreams:   streams,
		BizCode:     "JST",
		Language:    "Java",
		OS:          "Linux",
		StateType:   1, // 应用状态类型 （默认为无状态应用，1： 无状态应用，2：有状态应用， 3：守护进程集， 5：定时任务）
		ServiceType: "StoreApplication",
	}
}

// newCmdCreateApp is a command to create application.
func newCmdCreateApp(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newCreateAppOptions(streams)

	var cmd = &cobra.Command{
		Use:   "app",
		Short: "Create a new application.",
		Run: func(cmd *cobra.Command, args []string) {

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			as := f.AppService()

			resp, err := as.CreateApp(o.Title, o.BizCode, o.BizTitle, o.Description, o.Owner, o.Language, o.OS, o.ServiceType, o.Namespace, o.StateType)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			baseInfo := fmt.Sprintf("create application %s", o.Title)
			if resp.Code == 0 {
				printers.PrintValue(w, fmt.Sprintf("%s succeed, application id: %d\n", baseInfo, resp.Result.AppId))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}

	cmd.Flags().StringVar(&o.Title, "title", o.Title, "Title of application. (required)")
	cmd.Flags().StringVar(&o.BizCode, "bizCode", o.BizCode, "Value of 'JST','NEW_RETAIL','MINI_APP','SUPPLY','MESSAGE'")
	cmd.Flags().StringVar(&o.ServiceType, "serviceType", o.ServiceType, "Value of 'GeneralApplication','StoreApplication','TaoHuDongApplication'")
	cmd.Flags().StringVar(&o.Language, "lang", o.Language, "Language of application. Value of 'Java','PHP','C#','Python'")
	cmd.Flags().StringVar(&o.OS, "os", o.OS, "value of 'Linux','Windows'")
	cmd.Flags().StringVar(&o.Owner, "owner", o.Owner, "Owner user id of application. (required)")
	cmd.Flags().StringVar(&o.BizTitle, "bizTitle", o.BizTitle, "Short description of application. (required)")
	cmd.Flags().StringVar(&o.Description, "desc", o.Description, "Description of application.")
	cmd.Flags().StringVar(&o.Namespace, "ns", o.Namespace, "Namespace of application.")
	cmd.Flags().IntVar(&o.StateType, "stateType", o.StateType, "StateType of application. Value Of 1(Deployment), 2(StatefulSet), 5(CronJob)")
	util.MarkFlagRequired(cmd, o.ErrOut, "title")
	util.MarkFlagRequired(cmd, o.ErrOut, "owner")
	util.MarkFlagRequired(cmd, o.ErrOut, "bizTitle")
	return cmd
}
