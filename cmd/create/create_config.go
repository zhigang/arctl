package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// CreateAppConfigOptions is returned by newCreateAppConfigOptions
type CreateAppConfigOptions struct {
	AppId         int
	EnvType       int
	Name          string
	CodePath      string
	Deployment    string
	StatefulSet   string
	CronJob       string
	FilePath      string
	FileType      string
	ConfigMapList []string
	ConfigType    string
	printers.IOStreams
}

// newCreateAppConfigOptions returns an initialized CreateAppConfigOptions instance
func newCreateAppConfigOptions(streams printers.IOStreams) *CreateAppConfigOptions {
	return &CreateAppConfigOptions{
		IOStreams: streams,
		EnvType:   0,
		FileType:  "deploy",
	}
}

// newCmdCreateAppConfig is a command to create application's config.
func newCmdCreateAppConfig(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newCreateAppConfigOptions(streams)

	var cmd = &cobra.Command{
		Use:   "config",
		Short: "Create a new deploy config for application.",
		Run: func(cmd *cobra.Command, args []string) {

			envType := "test"
			if o.EnvType == 1 {
				envType = "online"
			}

			if o.FilePath != "" {
				cfg, err := readFileBase64(o.FilePath)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				switch o.FileType {
				case "stateful":
					o.StatefulSet = cfg
				case "cronjob":
					o.CronJob = cfg
				default:
					o.Deployment = cfg
				}
			}

			if len(o.ConfigMapList) > 0 && o.ConfigType == "path" {
				for i, v := range o.ConfigMapList {
					cfg, err := readFileBase64(v)
					if err != nil {
						util.CheckErr(o.ErrOut, err)
					}
					o.ConfigMapList[i] = cfg
				}
			}

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			as := f.AppService()

			resp, err := as.CreateDeployConfig(o.AppId, envType, o.Name, o.CodePath, o.Deployment, o.StatefulSet, o.CronJob, o.ConfigMapList)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			baseInfo := fmt.Sprintf("create application's deploy config %s", o.Name)
			if resp.Code == 0 {
				printers.PrintValue(w, fmt.Sprintf("%s succeed, config schema id: %d\n", baseInfo, resp.Result.SchemaId))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMessage)
			}
		},
	}
	cmd.Flags().IntVar(&o.AppId, "appID", o.AppId, "ID of application. (required)")
	cmd.Flags().IntVarP(&o.EnvType, "type", "t", o.EnvType, "Type of environment. Value of 0 [test], 1 [production] (required)")
	cmd.Flags().StringVar(&o.Name, "name", o.Name, "Name of application's deploy config. (required)")

	cmd.Flags().StringVarP(&o.FilePath, "file", "f", o.FilePath, "File path of application's yaml. Type of Deployment, StatefulSet, CronJob")
	cmd.Flags().StringVar(&o.FileType, "fileType", o.FileType, "Value of 'deploy', 'stateful', 'cronjob'. Useful with flag 'file'")

	cmd.Flags().StringVar(&o.CodePath, "codePath", o.CodePath, "Code path of application's deploy config.")
	cmd.Flags().StringVar(&o.Deployment, "deploy", o.Deployment, "Deployment base64 string of application's deploy config.")
	cmd.Flags().StringVar(&o.StatefulSet, "stateful", o.StatefulSet, "StatefulSet base64 string of application's deploy config.")
	cmd.Flags().StringVar(&o.CronJob, "cronjob", o.CronJob, "CronJob base64 string of application's deploy config.")

	cmd.Flags().StringSliceVar(&o.ConfigMapList, "configMap", o.ConfigMapList, "Config map base64 string list of application's deploy config.")
	cmd.Flags().StringVar(&o.ConfigType, "configType", o.ConfigType, "Value of 'path', 'content'. Useful with flag 'configMap'")

	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")

	return cmd
}
