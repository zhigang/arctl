package describe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdDescribeApp show details of a specific application.
func newCmdDescribeApp(o *DescribeOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:                   "app [APPLICATION_NAME]",
		Short:                 "Show details of a specific application",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			ns := f.AppService()

			var app retailcloud.AppDetail
			if o.ID == "" {
				if len(args) == 0 {
					util.ExitWithMsg(o.ErrOut, "[APPLICATION_NAME] or --id not set")
				}

				o.Name = strings.ToLower(args[0])
				resp, err := ns.GetAppList(o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				for _, dtl := range resp.Data {
					if strings.ToLower(dtl.Title) == o.Name {
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
				util.ExitWithMsg(o.ErrOut, "application %s not found", o.ID)
			}

			headers := []string{"App ID", "App Title", "Biz Name", "Biz Title", "Deploy Type", "Language", "State Type", "OS", "ServiceType", "Description"}
			appInfo := []string{strconv.FormatInt(app.AppId, 10), app.Title, app.BizName, app.BizTitle, app.DeployType, app.Language, app.AppStateType, app.OperatingSystem, app.ServiceType, app.Description}
			printers.PrintTable(o.Out, headers, appInfo)

			es := f.EnvironmentService()
			envs, err := es.GetEnvList(int(app.AppId), o.PageNumber, o.PageSize, o.Type, "")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			for _, e := range envs.Data {
				resp, _ := es.GetEnvDetail(int(app.AppId), int(e.EnvId))
				envDtl := resp.Result

				ed, err := ns.GetAppInstance(int(app.AppId), int(e.EnvId), o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				runningCnt := 0
				instances := [][]string{}
				for _, v := range ed.Data {
					if strings.ToLower(v.Health) == "running" {
						runningCnt += 1
					}
					instances = append(instances, []string{v.AppInstanceId, v.PodIp, v.HostIp, strconv.Itoa(v.RestartCount), v.Health, v.Requests, v.Limits, v.CreateTime})
				}

				headers := []string{"Env ID", "Env Name", "Env Type", "Replicas", "Config ID", "Region"}
				envInfo := []string{strconv.FormatInt(e.EnvId, 10), e.EnvName, e.EnvTypeName, fmt.Sprintf("%d/%d", runningCnt, envDtl.Replicas), strconv.FormatInt(e.AppSchemaId, 10), e.Region}
				printers.PrintTable(o.Out, headers, envInfo)

				if len(instances) > 0 {
					iHeaders := []string{"Instance ID", "Pod IP", "Host IP", "Restart", "Health", "Requests", "Limits", "Create Time"}
					printers.PrintTable(o.Out, iHeaders, instances...)
				}
			}
		},
	}
	cmd.Flags().StringVarP(&o.ID, "id", "i", o.ID, "ID of application.")
	return cmd
}
