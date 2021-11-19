package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// CreateEnvOptions is returned by newCreateEnvOptions
type createEnvOptions struct {
	AppID     int
	SchemaID  int
	Replicas  int
	EnvType   int
	EnvName   string
	Region    string
	ClusterID string
	printers.IOStreams
}

// newCreateEnvOptions returns an initialized createEnvOptions instance
func newCreateEnvOptions(streams printers.IOStreams) *createEnvOptions {
	return &createEnvOptions{
		IOStreams: streams,
	}
}

// newCmdCreateEnv is a command to create application's environment.
func newCmdCreateEnv(streams printers.IOStreams, f util.Factory) *cobra.Command {
	o := newCreateEnvOptions(streams)

	var cmd = &cobra.Command{
		Use:   "env",
		Short: "Create a new environment.",
		Run: func(cmd *cobra.Command, args []string) {

			w := printers.GetNewTabWriter(o.Out)
			defer w.Flush()

			es := f.EnvironmentService()

			resp, err := es.CreateEnv(o.AppID, o.SchemaID, o.Replicas, o.EnvType, o.EnvName, o.Region, o.ClusterID)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			baseInfo := fmt.Sprintf("create environment %s", o.EnvName)
			if resp.Code == 0 {
				printers.PrintValue(w, fmt.Sprintf("%s succeed, environment id: %d\n", baseInfo, resp.Result.AppEnvId))
			} else {
				util.ExitWithMsg(o.ErrOut, "%s failed, error msg: %s\n", baseInfo, resp.ErrMsg)
			}
		},
	}
	cmd.Flags().StringVar(&o.EnvName, "name", o.EnvName, "Name of environment.")
	cmd.Flags().IntVarP(&o.EnvType, "type", "t", o.EnvType, "Type of environment. Value of 0 [test], 1 [production] ")
	cmd.Flags().IntVar(&o.AppID, "appID", o.AppID, "ID of application.")
	cmd.Flags().IntVar(&o.SchemaID, "schema", o.SchemaID, "Schema ID is application's deploy config ID.")
	cmd.Flags().IntVar(&o.Replicas, "replicas", o.Replicas, "Replicas of application in environment.")
	cmd.Flags().StringVar(&o.Region, "region", o.Region, "Region of environment.")
	cmd.Flags().StringVar(&o.ClusterID, "cluster", o.ClusterID, "Cluster ID of environment. (required)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("appID")
	cmd.MarkFlagRequired("schema")
	cmd.MarkFlagRequired("replicas")
	cmd.MarkFlagRequired("region")
	cmd.MarkFlagRequired("cluster")
	return cmd
}
