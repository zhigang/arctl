package get

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/pool"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetApps creates a command object for get applications
func newCmdGetApps(o *getOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "apps",
		Short: "List apps",
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()

			showAll, err := cmd.Flags().GetBool("all")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			showEnvs, err := cmd.Flags().GetBool("show-envs")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			as := f.AppService()
			var apps []retailcloud.AppDetail
			if o.ID == "" {
				resp, err := as.GetAppList(o.PageNumber, o.PageSize)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				apps = append(apps, resp.Data...)
			} else {
				id, err := strconv.Atoi(o.ID)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				resp, err := as.GetApp(id)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				app := resp.Result
				if app.AppId > 0 {
					apps = append(apps, retailcloud.AppDetail{
						AppId: app.AppId, Title: app.Title,
						AppStateType: app.AppStateType, BizName: app.BizName,
						BizTitle: app.BizTitle, DeployType: app.DeployType,
						Language: app.Language, OperatingSystem: app.OperatingSystem,
					})
				}
			}

			if len(apps) == 0 {
				util.ExitWithMsg(o.ErrOut, "applications not found")
			}

			headers := []string{"App ID", "App Title", "STATE TYPE"}
			sorts := [][]string{}
			if showEnvs {
				envHeaders := []string{"Env ID", "Env Name", "Env Type", "Replicas", "Config ID", "Region"}
				headers = append(headers, envHeaders...)
				c, err := f.Config()
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				p, err := pool.NewPool(c.Pool.Size)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				for _, v := range apps {
					w := &pool.Worker{
						Process: func(param ...interface{}) interface{} {
							app := param[0].(retailcloud.AppDetail)
							opt := param[1].(*getOptions)
							fct := param[2].(util.Factory)
							all := param[3].(bool)
							es := fct.EnvironmentService()
							bases := []string{strconv.FormatInt(app.AppId, 10), app.Title, app.AppStateType}
							envs, err := es.GetEnvList(int(app.AppId), opt.PageNumber, opt.PageSize, opt.Type, "")
							if err != nil {
								util.CheckErr(o.ErrOut, err)
							}
							infors := [][]string{}
							for _, env := range envs.Data {
								dtl, err := es.GetEnvDetail(int(app.AppId), int(env.EnvId))
								if err != nil {
									util.CheckErr(o.ErrOut, err)
								}
								e := dtl.Result
								ed, err := as.GetAppInstance(int(app.AppId), int(env.EnvId), o.PageNumber, o.PageSize)
								if err != nil {
									util.CheckErr(o.ErrOut, err)
								}
								runningCnt := 0
								for _, v := range ed.Data {
									if strings.ToLower(v.Health) == "running" {
										runningCnt++
									}
								}
								rows := append(bases, strconv.FormatInt(e.EnvId, 10), e.EnvName, e.EnvTypeName, fmt.Sprintf("%d/%d", runningCnt, e.Replicas), strconv.FormatInt(e.AppSchemaId, 10), e.Region)
								infors = append(infors, rows)
							}
							if len(envs.Data) == 0 && all {
								rows := append(bases, "-", "-", "-", "-", "-", "-")
								infors = append(infors, rows)
							}
							return infors
						},
						Param: []interface{}{v, o, f, showAll},
					}
					p.Put(w)
				}

				for d := range p.Run() {
					if d.Result != nil {
						data := d.Result.([][]string)
						sorts = append(sorts, data...)
					}
				}
			} else {
				subHeaders := []string{"Biz Name", "Biz Title", "Deploy Type", "Language", "OS"}
				headers = append(headers, subHeaders...)
				for _, app := range apps {
					data := []string{strconv.FormatInt(app.AppId, 10), app.Title, app.AppStateType, app.BizName, app.BizTitle, app.DeployType, app.Language, app.OperatingSystem}
					sorts = append(sorts, data)
				}
			}

			// sort by app id
			sort.Slice(sorts, func(i, j int) bool {
				a, _ := strconv.Atoi(sorts[i][0])
				b, _ := strconv.Atoi(sorts[j][0])
				return a > b
			})

			w.SetHeader(headers)
			footer := make([]string, len(headers))
			footer[len(footer)-2] = "Total"
			footer[len(footer)-1] = strconv.Itoa(len(sorts))
			w.SetFooter(footer)

			w.AppendBulk(sorts)
		},
	}
	cmd.Flags().Bool("show-envs", false, "show application when it's environment found")
	cmd.Flags().BoolP("all", "a", false, "show all applications when flag 'show-envs' is used")
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	return cmd
}
