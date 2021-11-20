package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/cmd/create"
	"github.com/zhigang/arctl/cmd/delete"
	"github.com/zhigang/arctl/cmd/deploy"
	"github.com/zhigang/arctl/cmd/describe"
	"github.com/zhigang/arctl/cmd/get"
	"github.com/zhigang/arctl/cmd/label"
	"github.com/zhigang/arctl/cmd/scale"
	"github.com/zhigang/arctl/cmd/version"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/services"
	"github.com/zhigang/arctl/util"
)

var (
	cfgFile   string
	f         util.Factory
	ioStreams = printers.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
)

func main() {

	f = services.NewFactory()

	cobra.OnInitialize(func() {
		services.SetConfigFilePath(cfgFile)
	})

	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "arctl",
		Short: "arctl controls the Aliyun retail cloud",
		Run:   runHelp,
	}

	cmds.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.arctl.yml)")

	cmds.AddCommand(version.NewCmdVersion(ioStreams, f))
	cmds.AddCommand(get.NewCmdGet(ioStreams, f))
	cmds.AddCommand(describe.NewCmdDescribe(ioStreams, f))
	cmds.AddCommand(scale.NewCmdScale(ioStreams, f))
	cmds.AddCommand(create.NewCmdCreate(ioStreams, f))
	cmds.AddCommand(delete.NewCmdDelete(ioStreams, f))
	cmds.AddCommand(deploy.NewCmdDeploy(ioStreams, f))
	cmds.AddCommand(label.NewCmdLabel(ioStreams, f))

	if err := cmds.Execute(); err != nil {
		util.CheckErr(ioStreams.ErrOut, err)
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		util.CheckErr(ioStreams.ErrOut, err)
	}
}
