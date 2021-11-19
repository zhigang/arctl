package get

import (
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetAppConfig creates a command object for get application's configs
func newCmdGetAppConfig(o *getOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "configs [APPLICATION_NAME]",
		Short: "List deploy config of application",
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()

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

			as := f.AppService()
			var app retailcloud.AppDetail
			if appID > 0 {
				resp, err := as.GetApp(appID)
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
				resp, err := as.GetAppList(o.PageNumber, o.PageSize)
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
				util.ExitWithMsg(o.ErrOut, "applicatiuon not found")
			}

			headers := []string{"App ID", "App Title", "Config ID", "Config Name", "Config Type", "Container Yaml"}
			appInfo := []string{strconv.FormatInt(app.AppId, 10), app.Title}

			resp, err := as.GetDeployConfig(int(app.AppId), id, "", envType)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			for _, c := range resp.Data {
				resType := "Deployment"
				if c.ContainerCodePath.CodePath != "" {
					resType = "Code Path"
				} else if c.ContainerYamlConf.Deployment != "" {
					resType = "Deployment"
				} else if c.ContainerYamlConf.StatefulSet != "" {
					resType = "StatefulSet"
				} else if c.ContainerYamlConf.CronJob != "" {
					resType = "CronJob"
				}
				one := []string{strconv.FormatInt(c.Id, 10), c.Name, c.EnvType, resType}
				all := append(appInfo, one...)
				w.Append(all)
			}

			w.SetHeader(headers)
			footer := make([]string, len(headers))
			footer[len(footer)-2] = "Total"
			footer[len(footer)-1] = strconv.Itoa(w.NumLines())
			w.SetFooter(footer)
		},
	}
	cmd.Flags().Int("appID", 0, "ID of application.")
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	return cmd
}
