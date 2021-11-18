package get

import (
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetPods creates a command object for get pods from a specific application.
func newCmdGetPods(o *GetOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:                   "pods [APPLICATION_NAME]",
		Short:                 "Get pods of a specific application",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()

			ns := f.AppService()
			var app retailcloud.AppDetail
			if o.ID == "" {
				if len(args) == 0 {
					util.ExitWithMsg(o.ErrOut, "[APPLICATION_NAME] or --id not set")
				}

				application := strings.ToLower(args[0])
				resp, err := ns.GetAppList(o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				for _, dtl := range resp.Data {
					if strings.ToLower(dtl.Title) == application {
						app = dtl
						break
					}
				}

			} else {
				id, err := strconv.Atoi(o.ID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				resp, err := ns.GetApp(id)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				res := resp.Result
				if res.AppId > 0 {
					app = retailcloud.AppDetail{
						AppId:           res.AppId,
						Title:           res.Title,
						BizName:         res.BizName,
						BizTitle:        res.BizTitle,
						DeployType:      res.DeployType,
						Language:        res.Language,
						AppStateType:    res.AppStateType,
						OperatingSystem: res.OperatingSystem,
						ServiceType:     res.ServiceType,
						Description:     res.Description,
					}
				}
			}

			if app.AppId == 0 {
				util.ExitWithMsg(o.ErrOut, "application not found")
			}

			headers := []string{"App ID", "App Title", "Env ID", "Env Name", "Env Type", "Instance ID", "Pod IP", "Host IP", "Restart", "Health", "Requests", "Limits", "Create Time"}
			w.SetHeader(headers)

			appInfo := []string{strconv.FormatInt(app.AppId, 10), app.Title}

			es := f.EnvironmentService()
			envs, err := es.GetEnvList(int(app.AppId), o.PageNumber, o.PageSize, o.Type, "")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			for _, e := range envs.Data {
				envInfo := []string{strconv.FormatInt(e.EnvId, 10), e.EnvName, strconv.Itoa(e.EnvType)}
				ed, err := ns.GetAppInstance(int(app.AppId), int(e.EnvId), o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				for _, v := range ed.Data {
					one := []string{v.AppInstanceId, v.PodIp, v.HostIp, strconv.Itoa(v.RestartCount), v.Health, v.Requests, v.Limits, v.CreateTime}
					all := append(appInfo, envInfo...)
					all = append(all, one...)
					w.Append(all)
				}
			}
			footer := make([]string, len(headers))
			footer[len(footer)-2] = "Total"
			footer[len(footer)-1] = strconv.Itoa(w.NumLines())
			w.SetFooter(footer)
		},
	}
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	return cmd
}
