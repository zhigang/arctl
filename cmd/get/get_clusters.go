package get

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetClusters creates a command object for get clusters
func newCmdGetClusters(o *getOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "clusters",
		Short: "List clusters",
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()

			envType := ""
			switch o.Type {
			case 0:
				envType = "TEST"
			case 1:
				envType = "PRO"
			default:
				envType = ""
			}

			cs := f.ClusterService()

			resp, err := cs.GetClusterList(envType, o.PageNumber, o.PageSize)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}
			if len(resp.Data) == 0 {
				util.ExitWithMsg(o.ErrOut, "clusters not found")
			}
			headers := []string{"Index", "ID", "Title", "Business Code", "Status", "Nodes", "Net", "Pod CIDR", "Service CIDR", "Env Type", "CPU", "Memery", "VPC", "Region"}
			for i, c := range resp.Data {
				if o.ID == "" {
					w.Append([]string{strconv.Itoa(i), c.InstanceId, c.ClusterTitle, c.BusinessCode, c.Status, strconv.Itoa(len(c.EcsIds)), c.NetPlug, c.PodCIDR, c.ServiceCIDR, c.EnvType, getPersent(c.WorkLoadCpu), getPersent(c.WorkLoadMem), c.VpcId, c.RegionName})
				} else if strings.ToLower(o.ID) == c.InstanceId {
					w.Append([]string{strconv.Itoa(i), c.InstanceId, c.ClusterTitle, c.BusinessCode, c.Status, strconv.Itoa(len(c.EcsIds)), c.NetPlug, c.PodCIDR, c.ServiceCIDR, c.EnvType, getPersent(c.WorkLoadCpu), getPersent(c.WorkLoadMem), c.VpcId, c.RegionName})
					break
				}
			}
			w.SetHeader(headers)
			footer := make([]string, len(headers))
			footer[len(footer)-2] = "Total"
			footer[len(footer)-1] = strconv.Itoa(w.NumLines())
			w.SetFooter(footer)
		},
	}
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	return cmd
}

func getPersent(value string) string {
	v, _ := strconv.ParseFloat(value, 32)
	return fmt.Sprintf("%.1f%%", v*100)
}
