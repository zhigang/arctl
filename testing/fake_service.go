package testing

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

var (
	ErrorNotFound = errors.New("resource not found")
)

const (
	APP1_ID      int64  = 12345
	APP1_ENV1_ID int64  = 1234511
	APP1_ENV2_ID int64  = 1234522
	APP2_ID      int64  = 54321
	APP2_ENV1_ID int64  = 1234533
	CLUSTER1_ID  string = "c-id-778899"
	CLUSTER2_ID  string = "c-id-990099"
)

type appServiceImpl struct {
	util.AppService
	factory util.Factory
}

func newAppService(factory util.Factory) *appServiceImpl {
	return &appServiceImpl{
		factory: factory,
	}
}

func (s *appServiceImpl) GetApp(appID int) (*retailcloud.DescribeAppDetailResponse, error) {
	resp := &retailcloud.DescribeAppDetailResponse{}
	if appID == 12345 {
		resp.Result = retailcloud.Result{
			AppId:           APP1_ID,
			Title:           "test1",
			AppStateType:    "stateless",
			BizName:         "JST",
			BizTitle:        "测试1",
			DeployType:      "ContainerDeployment",
			Language:        "Java",
			OperatingSystem: "Linux",
		}
	}
	return resp, nil
}

func (s *appServiceImpl) GetAppList(pageNumber, pageSize int) (*retailcloud.ListAppResponse, error) {
	resp := &retailcloud.ListAppResponse{}
	resp.Data = []retailcloud.AppDetail{}
	app1 := makeFakeApp(APP1_ID, "test1", "测试1")
	app2 := makeFakeApp(APP2_ID, "test2", "测试2")
	resp.Data = append(resp.Data, app1, app2)
	return resp, nil
}

func (s *appServiceImpl) GetDeployConfig(appID, id int, name, envType string) (*retailcloud.ListDeployConfigResponse, error) {
	resp := &retailcloud.ListDeployConfigResponse{}
	resp.Data = []retailcloud.DeployConfigInstance{}

	conf1 := retailcloud.DeployConfigInstance{
		Id:      1,
		Name:    "conf-test-1",
		AppId:   APP1_ID,
		EnvType: "test",
		ContainerYamlConf: retailcloud.ContainerYamlConf{
			Deployment: "yaml...content",
		},
	}
	conf2 := retailcloud.DeployConfigInstance{
		Id:      2,
		Name:    "conf-test-2",
		AppId:   APP1_ID,
		EnvType: "online",
		ContainerYamlConf: retailcloud.ContainerYamlConf{
			Deployment: "yaml...content",
		},
	}
	if appID == int(APP1_ID) {
		if envType == "test" {
			resp.Data = append(resp.Data, conf1)
		} else if envType == "online" {
			resp.Data = append(resp.Data, conf2)
		} else {
			resp.Data = append(resp.Data, conf1, conf2)
		}
	}
	return resp, nil
}

func (s *appServiceImpl) GetAppInstance(appID, envID, pageNumber, pageSize int) (*retailcloud.ListAppInstanceResponse, error) {
	resp := &retailcloud.ListAppInstanceResponse{}
	resp.Data = []retailcloud.AppInstanceDetail{}
	i1 := retailcloud.AppInstanceDetail{
		AppInstanceId: "jck-23638-24399-998880-585fb6c46-fzs2l",
		PodIp:         "172.20.1.147",
		HostIp:        "172.26.124.74",
		RestartCount:  0,
		Health:        "RUNNING",
		Requests:      "2vcpu 4GB",
		Limits:        "4vcpu 8GB",
		CreateTime:    "2021-11-18T08:30:20",
	}
	i2 := retailcloud.AppInstanceDetail{
		AppInstanceId: "jck-23638-24399-998880-585fb6c46-dfrec",
		PodIp:         "172.20.1.148",
		HostIp:        "172.26.124.74",
		RestartCount:  0,
		Health:        "RUNNING",
		Requests:      "2vcpu 4GB",
		Limits:        "4vcpu 8GB",
		CreateTime:    "2021-11-18T08:30:20",
	}
	if appID == int(APP1_ID) {
		if envID == int(APP1_ENV1_ID) {
			resp.Data = append(resp.Data, i1)
		} else if envID == int(APP1_ENV2_ID) {
			resp.Data = append(resp.Data, i2)
		} else {
			resp.Data = append(resp.Data, i1, i2)
		}
	}
	return resp, nil
}

func (s *appServiceImpl) DeployApp(envID, totalPartitions int, name, image string) (*retailcloud.DeployAppResponse, error) {
	resp := &retailcloud.DeployAppResponse{}
	resp.Success = true
	resp.Result.AppSchemaId = 1234
	return resp, nil
}

type clusterServiceImpl struct {
	util.ClusterService
	factory util.Factory
}

func newClusterService(factory util.Factory) *clusterServiceImpl {
	return &clusterServiceImpl{
		factory: factory,
	}
}

func (s *clusterServiceImpl) GetClusterList(envType string, pageNumber, pageSize int) (*retailcloud.ListClusterResponse, error) {
	resp := &retailcloud.ListClusterResponse{}
	resp.Data = []retailcloud.ClusterInfo{}
	test := retailcloud.ClusterInfo{
		InstanceId:   CLUSTER1_ID,
		ClusterTitle: "Test集群",
		BusinessCode: "JST",
		Status:       "running",
		NetPlug:      "terway",
		PodCIDR:      "172.20.0.0/16",
		ServiceCIDR:  "172.21.0.0/16",
		EnvType:      "0",
		WorkLoadCpu:  "0.65",
		WorkLoadMem:  "0.76",
		VpcId:        "vpc-test-id-1",
		RegionName:   "cn-zhangjiakou",
		EcsIds: []string{
			"ecs-test-id-1",
			"ecs-test-id-2",
		},
	}

	test2 := retailcloud.ClusterInfo{
		InstanceId:   CLUSTER2_ID,
		ClusterTitle: "PRO集群",
		BusinessCode: "JST",
		Status:       "running",
		NetPlug:      "terway",
		PodCIDR:      "172.21.0.0/16",
		ServiceCIDR:  "172.22.0.0/16",
		EnvType:      "1",
		WorkLoadCpu:  "0.81",
		WorkLoadMem:  "0.78",
		VpcId:        "vpc-test-id-2",
		RegionName:   "cn-zhangjiakou",
		EcsIds: []string{
			"ecs-test-id-3",
			"ecs-test-id-4",
			"ecs-test-id-5",
		},
	}

	if envType == "TEST" {
		resp.Data = append(resp.Data, test)
	} else if envType == "PRO" {
		resp.Data = append(resp.Data, test2)
	} else {
		resp.Data = append(resp.Data, test, test2)
	}
	return resp, nil
}

type environmentServiceImpl struct {
	util.EnvironmentService
	factory util.Factory
}

func newEnvironmentService(factory util.Factory) *environmentServiceImpl {
	return &environmentServiceImpl{
		factory: factory,
	}
}

func (s *environmentServiceImpl) GetEnvList(appID, pageNumber, pageSize, envType int, envName string) (*retailcloud.ListAppEnvironmentResponse, error) {
	resp := &retailcloud.ListAppEnvironmentResponse{}
	resp.Data = []retailcloud.AppEnvironmentResponse{}
	env1 := retailcloud.AppEnvironmentResponse{
		AppId:   APP1_ID,
		EnvId:   APP1_ENV1_ID,
		EnvName: "Test",
		EnvType: 0,
	}
	env2 := retailcloud.AppEnvironmentResponse{
		AppId:   APP1_ID,
		EnvId:   APP1_ENV2_ID,
		EnvName: "Production",
		EnvType: 1,
	}
	env3 := retailcloud.AppEnvironmentResponse{
		AppId:   APP2_ID,
		EnvId:   APP2_ENV1_ID,
		EnvName: "Test",
		EnvType: 0,
	}
	if appID == int(APP1_ID) {
		if envType == 0 {
			resp.Data = append(resp.Data, env1)
		} else if envType == 1 {
			resp.Data = append(resp.Data, env2)
		} else {
			resp.Data = append(resp.Data, env1, env2)
		}
	} else if appID == int(APP2_ID) && envType == 0 {
		resp.Data = append(resp.Data, env3)
	}
	return resp, nil
}

type nodeServiceImpl struct {
	util.NodeService
	factory util.Factory
}

func newNodeService(factory util.Factory) *nodeServiceImpl {
	return &nodeServiceImpl{
		factory: factory,
	}
}

func (s *nodeServiceImpl) GetClusterNodes(clusterID string, pageNumber, pageSize int) (*retailcloud.ListClusterNodeResponse, error) {
	resp := &retailcloud.ListClusterNodeResponse{}
	resp.Data = []retailcloud.ClusterNodeInfo{}

	n1 := makeFakeNode("ecs-test-id-1", "ecs-test-1", "172.26.124.74")
	n2 := makeFakeNode("ecs-test-id-2", "ecs-test-2", "172.26.124.75")
	n3 := makeFakeNode("ecs-test-id-3", "ecs-test-3", "172.26.124.76")
	n4 := makeFakeNode("ecs-test-id-4", "ecs-test-4", "172.26.124.77")
	n5 := makeFakeNode("ecs-test-id-5", "ecs-test-5", "172.26.124.78")

	if clusterID == CLUSTER1_ID {
		resp.Data = append(resp.Data, n1, n2)
	} else if clusterID == CLUSTER2_ID {
		resp.Data = append(resp.Data, n3, n4, n5)
	}

	return resp, nil
}

type userServiceImpl struct {
	util.UserService
	factory util.Factory
}

func newUserService(factory util.Factory) *userServiceImpl {
	return &userServiceImpl{
		factory: factory,
	}
}
func (s *userServiceImpl) GetUserList(pageNumber, pageSize int) (*retailcloud.ListUsersResponse, error) {
	resp := &retailcloud.ListUsersResponse{}
	resp.Data = []retailcloud.ListUserResponse{
		{
			UserId:   "1",
			RealName: "user1",
			UserType: "TAOBAO",
		},
		{
			UserId:   "2",
			RealName: "user2",
			UserType: "TAOBAO",
		},
	}
	return resp, nil
}

func makeFakeApp(appId int64, title, bizTitle string) retailcloud.AppDetail {
	return retailcloud.AppDetail{
		AppId:           appId,
		Title:           title,
		AppStateType:    "stateless",
		BizName:         "JST",
		BizTitle:        bizTitle,
		DeployType:      "ContainerDeployment",
		Language:        "Java",
		OperatingSystem: "Linux",
	}
}

func makeFakeNode(id, name, ip string) retailcloud.ClusterNodeInfo {
	return retailcloud.ClusterNodeInfo{
		InstanceId:     id,
		InstanceName:   name,
		EcsPrivateIp:   ip,
		EcsCpu:         "16",
		EcsMemory:      "32768",
		OSName:         "Linux",
		RegionId:       "cn-zhangjiakou",
		EcsExpiredTime: "2099-12-01",
	}
}
