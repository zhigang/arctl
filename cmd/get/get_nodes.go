package get

import (
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/spf13/cobra"
	"github.com/zhigang/arctl/pool"
	"github.com/zhigang/arctl/printers"
	"github.com/zhigang/arctl/util"
)

// newCmdGetNodes creates a command object for get nodes from a cluster.
func newCmdGetNodes(o *getOptions, f util.Factory) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "nodes",
		Short: "List nodes",
		Run: func(cmd *cobra.Command, args []string) {
			w := printers.GetNewTableWriter(o.Out)
			defer w.Render()
			showLabels, err := cmd.Flags().GetBool("show-labels")
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			cs := f.ClusterService()

			envType := ""
			if o.ID == "" {
				switch o.Type {
				case 0:
					envType = "TEST"
				case 1:
					envType = "PRO"
				default:
					envType = "TEST"
				}
			}

			respCluster, err := cs.GetClusterList(envType, o.PageNumber, o.PageSize)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			if len(respCluster.Data) == 0 {
				util.ExitWithMsg(o.ErrOut, "clusters not found")
			}

			var cluster retailcloud.ClusterInfo
			if o.ID == "" {
				index := 0
				fIndex, err := cmd.Flags().GetInt("index")
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				if fIndex > 0 {
					index = fIndex
				}
				if index >= len(respCluster.Data) {
					util.ExitWithMsg(o.ErrOut, "index out of range with length %d", len(respCluster.Data))
				}
				cluster = respCluster.Data[index]
			} else {
				for _, c := range respCluster.Data {
					if strings.ToLower(o.ID) == c.InstanceId {
						cluster = c
						break
					}
				}
			}

			if cluster.InstanceId == "" {
				util.ExitWithMsg(o.ErrOut, "cluster not found")
			}

			o.ID = cluster.InstanceId
			cHeaders := []string{"Cluster ID", "Cluster Title", "Env Type", "Nodes", "REGION"}
			printers.Print(o.Out, cHeaders, []string{cluster.InstanceId, cluster.ClusterTitle, cluster.EnvType, strconv.Itoa(len(cluster.EcsIds)), cluster.RegionId})

			ns := f.NodeService()
			resp, err := ns.GetClusterNodes(o.ID, o.PageNumber, o.PageSize)
			if err != nil {
				util.CheckErr(o.ErrOut, err)
			}

			headers := []string{"Instance ID", "Instance Name", "Private IP", "CPU", "Memory", "OS", "Region", "Expired Time"}
			if showLabels || o.LabelSelector != "" {
				headers = append(headers, "Labels")
				c, err := f.Config()
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}
				p, err := pool.NewPool(c.Pool.Size)
				if err != nil {
					util.CheckErr(o.ErrOut, err)
				}

				for _, v := range resp.Data {
					w := &pool.Worker{
						Process: func(param ...interface{}) interface{} {
							node := param[0].(retailcloud.ClusterNodeInfo)
							opt := param[1].(*getOptions)
							fct := param[2].(util.Factory)
							matched, labels, err := getNodeLabels(node, opt, fct)
							if err != nil {
								util.CheckErr(o.ErrOut, err)
							}
							if matched {
								infors := []string{node.InstanceId, node.InstanceName, node.EcsPrivateIp, node.EcsCpu, convertMBToGB(node.EcsMemory), node.EcsOsType, node.RegionId, node.EcsExpiredTime, labels}
								return infors
							}
							// time.Sleep(time.Duration(10) * time.Microsecond)
							return nil
						},
						Param: []interface{}{v, o, f},
					}
					if err := p.Put(w); err != nil {
						util.CheckErr(o.ErrOut, err)
					}
				}

				for d := range p.Run() {
					if d.Result != nil {
						data := d.Result.([]string)
						w.Append(data)
					}
				}
			} else {
				for _, node := range resp.Data {
					w.Append([]string{node.InstanceId, node.InstanceName, node.EcsPrivateIp, node.EcsCpu, convertMBToGB(node.EcsMemory), node.EcsOsType, node.RegionId, node.EcsExpiredTime})
				}
			}

			if w.NumLines() > 0 {
				w.SetHeader(headers)
				footer := make([]string, len(headers))
				footer[len(footer)-2] = "Total"
				footer[len(footer)-1] = strconv.Itoa(w.NumLines())
				w.SetFooter(footer)
			}
		},
	}
	cmd.Flags().Bool("show-labels", false, "show labels for resources")
	cmd.Flags().Int("index", 0, "index of cluser, default 0. if set id then index is ignore.")
	cmd.Flags().IntVarP(&o.Type, "type", "t", o.Type, "Selector (type query) to filter on type.(e.g. -t 0 or 1, Type: 0: test; 1: production)")
	cmd.Flags().StringVarP(&o.LabelSelector, "selector", "l", o.LabelSelector, "Selector (label query) to filter on, supports '='.(e.g. -l key1=value1,key2=value2)")
	return cmd
}

func convertMBToGB(mem string) string {
	mb, err := strconv.Atoi(mem)
	if err != nil {
		return "-1"
	}
	return strconv.Itoa(mb/1024) + "GB"
}

func getNodeLabels(v retailcloud.ClusterNodeInfo, o *getOptions, f util.Factory) (bool, string, error) {
	ns := f.NodeService()
	matchLabels := false
	if o.LabelSelector == "" {
		matchLabels = true
	}
	labels := ""
	labResp, err := ns.GetNodeLabels(o.ID, v.InstanceId, "", "", o.PageNumber, o.PageSize)
	if err != nil {
		labels = "load failed"
	}
	if labResp != nil && len(labResp.Data) > 0 {
		labelSelector, _, _ := util.ParseLabels(o.LabelSelector)
		for _, v := range labResp.Data {
			if !matchLabels && labelSelector[v.LabelKey] == v.LabelValue {
				matchLabels = true
			}
			labels += v.LabelKey + "=" + v.LabelValue + ","
		}
		labels = labels[0 : len(labels)-1]
	}
	return matchLabels, labels, err
}
