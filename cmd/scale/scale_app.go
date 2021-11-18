package scale

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// ScaleOptions is returned by newScaleOptions
type ScaleOptions struct {
	AppID      string
	AppName    string
	EnvID      string
	EnvName    string
	Replicas   int
	Type       int
	PageNumber int
	PageSize   int
	printers.IOStreams
}

// newScaleOptions returns an initialized ScaleOptions instance
func newScaleOptions(streams printers.IOStreams) *ScaleOptions {
	return &ScaleOptions{
		IOStreams:  streams,
		Type:       0,
		PageNumber: 1,
		PageSize:   999,
	}
}

// NewCmdScale returns new initialized instance of scale command
func NewCmdScale(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newScaleOptions(streams)

	var cmd = &cobra.Command{
		Use:   "scale [APPLICATION_NAME] [ENVIRONMENT_NAME]",
		Short: "Set a new size for a application.",
		Run: func(cmd *cobra.Command, args []string) {

			if o.Replicas < 0 {
				err := fmt.Errorf("the --replicas=COUNT flag is required, and COUNT must be greater than or equal to 0, current is %d", o.Replicas)
				util.CheckErr(o.ErrOut, err)
			}

			if len(args) > 0 {
				o.AppName = args[0]
				if len(args) > 1 {
					o.EnvName = args[1]
				}
			}

			if o.AppID == "" && o.AppName == "" {
				util.ExitWithMsg(o.ErrOut, "[APPLICATION_NAME] or --appid not set")
			}

			if o.EnvID == "" && o.EnvName == "" {
				util.ExitWithMsg(o.ErrOut, "[ENVIRONMENT_NAME] or --envid not set")
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
				envType := -1
				envs, err := es.GetEnvList(int(app.AppId), o.PageNumber, o.PageSize, envType, "")
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

			baseInfo := fmt.Sprintf("scale app (%d)%s env (%d)%s size %d to %d", app.AppId, app.Title, env.EnvId, env.EnvName, env.Replicas, o.Replicas)
			if o.Replicas != env.Replicas {
				resp, err := ns.ScaleApp(int(env.EnvId), o.Replicas)
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
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, same replicas.\n", baseInfo)
			}
		},
	}
	cmd.Flags().StringVar(&o.AppID, "appid", o.AppID, "ID of application. useful with envid.")
	cmd.Flags().StringVar(&o.EnvID, "envid", o.EnvID, "ID of environment.")
	cmd.Flags().IntVar(&o.Replicas, "replicas", -1, "The new desired number of replicas. Required.")
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	cmd.MarkFlagRequired("replicas")
	cmd.MarkFlagRequired("type")
	return cmd
}
