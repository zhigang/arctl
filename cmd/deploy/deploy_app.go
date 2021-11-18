package deploy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// DeployOptions is returned by newDeployOptions
type DeployOptions struct {
	Name            string
	AppID           string
	AppName         string
	EnvID           string
	EnvName         string
	Image           string
	Type            int
	TotalPartitions int
	PageNumber      int
	PageSize        int
	printers.IOStreams
}

// newDeployOptions returns an initialized DeployOptions instance
func newDeployOptions(streams printers.IOStreams) *DeployOptions {
	return &DeployOptions{
		IOStreams:       streams,
		Type:            0,
		TotalPartitions: 2,
		PageNumber:      1,
		PageSize:        999,
	}
}

// NewCmdDeploy returns new initialized instance of deploy sub command
func NewCmdDeploy(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newDeployOptions(streams)

	var cmd = &cobra.Command{
		Use:   "deploy [APPLICATION_NAME] [ENVIRONMENT_NAME]",
		Short: "deploy a new images for a application.",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				o.AppName = args[0]
				if len(args) > 1 {
					o.EnvName = args[1]
				}
			}

			if o.AppID == "" && o.AppName == "" {
				util.ExitWithMsg(o.ErrOut, " [APPLICATION_NAME] or --appid not set")
			}

			if o.EnvID == "" && o.EnvName == "" {
				util.ExitWithMsg(o.ErrOut, " [ENVIRONMENT_NAME] or --envid not set")
			}

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			ns := f.AppService()
			var app retailcloud.AppDetail
			if o.AppID == "" {
				resp, err := ns.GetAppList(o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				for _, v := range resp.Data {
					if strings.EqualFold(v.Title, o.AppName) {
						app = v
						break
					}
				}
			} else {
				id, err := strconv.Atoi(o.AppID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				resp, err := ns.GetApp(id)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				res := resp.Result
				if res.AppId > 0 {
					app = retailcloud.AppDetail{AppId: res.AppId, Title: res.Title}
				}
			}

			if app.AppId == 0 {
				util.ExitWithMsg(o.ErrOut, "application %s not found", o.AppID)
			}

			es := f.EnvironmentService()
			var env retailcloud.Result
			if o.EnvID == "" {
				envs, err := es.GetEnvList(int(app.AppId), o.PageNumber, o.PageSize, o.Type, "")
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				for _, e := range envs.Data {

					if strings.EqualFold(e.EnvName, o.EnvName) {
						resp, err := es.GetEnvDetail(int(app.AppId), int(e.EnvId))
						if err != nil {
							util.CheckErr(o.ErrOut, err)
						}
						env = resp.Result
						break
					}
				}
			} else {
				envID, err := strconv.Atoi(o.EnvID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				resp, err := es.GetEnvDetail(int(app.AppId), envID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				env = resp.Result
			}

			if env.EnvId == 0 {
				util.ExitWithMsg(o.ErrOut, "environment %s not found", o.EnvID)
			}

			baseInfo := fmt.Sprintf("deploy app (%d)%s env (%d)%s image to %s", app.AppId, app.Title, env.EnvId, env.EnvName, o.Image)

			resp, err := ns.DeployApp(int(env.EnvId), o.TotalPartitions, o.Name, o.Image)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			if resp.Success {
				if resp.Result.Admitted {
					printers.PrintValue(w, fmt.Sprintf("%s succeed, deploy order id: %d\n", baseInfo, resp.Result.DeployOrderId))
				} else {
					msg := util.GetBusinessMessage(resp.Result.BusinessCode)
					util.ExitWithMsg(o.ErrOut, "%s not admitted, business message: %s\n", baseInfo, msg)
				}
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}

	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Name of deploy. (required)")
	cmd.Flags().StringVar(&o.AppID, "appid", o.AppID, "ID of application. useful with envid.")
	cmd.Flags().StringVar(&o.EnvID, "envid", o.EnvID, "ID of environment.")
	cmd.Flags().StringVar(&o.Image, "image", o.Image, "The new desired number of replicas. (required)")
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production) (required)")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("image")
	cmd.MarkFlagRequired("type")

	return cmd
}
