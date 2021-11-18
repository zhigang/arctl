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

// newCmdDescribeAppConfig show details of a specific application's deploy config.
func newCmdDescribeAppConfig(o *DescribeOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "config [APPLICATION_NAME]",
		Short: "Show details of application's deploy config",
		Run: func(cmd *cobra.Command, args []string) {

			appID, err := cmd.Flags().GetInt("appID")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			id := 0
			if o.ID != "" {
				id, err = strconv.Atoi(o.ID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
			}

			envType := ""
			switch o.Type {
			case 0:
				envType = "test"
			case 1:
				envType = "online"
			default:
				envType = ""
			}

			ns := f.AppService()
			var app retailcloud.AppDetail
			if appID > 0 {
				resp, err := ns.GetApp(appID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				r := resp.Result
				if r.AppId > 0 {
					app = retailcloud.AppDetail{
						AppId: r.AppId, Title: r.Title,
						AppStateType: r.AppStateType, BizName: r.BizName,
						BizTitle: r.BizTitle, DeployType: r.DeployType,
						Language: r.Language, OperatingSystem: r.OperatingSystem,
					}
				}
			} else {
				if len(args) == 0 {
					util.ExitWithMsg(o.ErrOut, "[APPLICATION_NAME] or --appID not set")
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
			}

			if app.AppId == 0 {
				util.ExitWithMsg(o.ErrOut, "application not found")
			}

			w := printers.GetNewTabWriter(o.Out)

			printers.PrintObj(w, "App ID", strconv.FormatInt(app.AppId, 10))
			printers.PrintObj(w, "App Title", app.Title)
			printers.PrintObj(w, "Biz Name", app.BizName)
			printers.PrintObj(w, "Biz Title", app.Title)
			printers.PrintObj(w, "Deploy Type", app.DeployType)
			printers.PrintObj(w, "Language", app.Language)
			printers.PrintObj(w, "State Type", app.AppStateType)
			fmt.Fprintln(w)

			resp, err := ns.GetDeployConfig(int(app.AppId), id, o.Name, envType)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			for _, c := range resp.Data {
				printers.PrintObj(w, "Config ID", strconv.FormatInt(c.Id, 10))
				printers.PrintObj(w, "Config Name", c.Name)
				printers.PrintObj(w, "Config Type", c.EnvType)
				if c.ContainerCodePath.CodePath != "" {
					printers.PrintObj(w, "Code Path", c.ContainerCodePath.CodePath)
				} else if c.ContainerYamlConf.Deployment != "" {
					printers.PrintObj(w, "Deployment", c.ContainerYamlConf.Deployment)
				} else if c.ContainerYamlConf.StatefulSet != "" {
					printers.PrintObj(w, "StatefulSet", c.ContainerYamlConf.StatefulSet)
				} else if c.ContainerYamlConf.CronJob != "" {
					printers.PrintObj(w, "CronJob", c.ContainerYamlConf.CronJob)
				}
				fmt.Fprintln(w)
				if len(c.ContainerYamlConf.ConfigMapList) > 0 {
					for i, v := range c.ContainerYamlConf.ConfigMapList {
						printers.PrintObj(w, "Config_"+strconv.Itoa(i+1), v)
						fmt.Fprintln(w)
					}
				}
			}
			w.Flush()
		},
	}

	cmd.Flags().Int("appID", 0, "ID of application.")
	cmd.Flags().StringVarP(&o.ID, "id", "i", o.ID, "ID of config. If set id then ignore name of config.")
	cmd.Flags().StringVarP(&o.Name, "name", "n", o.ID, "Name of config.")

	return cmd
}
