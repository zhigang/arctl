package util

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
)

// AppService is a application service of retail client
type AppService interface {
	GetApp(appID int) (*retailcloud.DescribeAppDetailResponse, error)
	DeleteApp(appID int, force bool) (*retailcloud.DeleteAppDetailResponse, error)
	GetAppInstance(appID, envID, pageNumber, pageSize int) (*retailcloud.ListAppInstanceResponse, error)
	GetAppList(pageNumber, pageSize int) (*retailcloud.ListAppResponse, error)
	// TotalPartitions: 发布划分的批次数。
	// 建议至少划分为2个批次，划分多个批次的 好处主要有两点：
	// 1、做金丝雀灰度，防止本次变更有问题而 导致服务整体宕机；建议在第一批次发布 完成后暂停观察一段时间，若通过系统监 控或查看日志发现变更有问题，则可以停 止发布并回滚；
	// 2、分批重启应用实例，防止发布期间所有 实例同时重启而导致服务不可用；
	// 具体的错误码。枚举值：102000：没有发布准入凭据； 102001：发布准入凭据审核中
	DeployApp(envID, totalPartitions int, name, image string) (*retailcloud.DeployAppResponse, error)
	ScaleApp(envID, replicas int) (*retailcloud.ScaleAppResponse, error)
	// CreateApp is create a application
	// title, bizCode, ownerID, language, os, serviceType is required
	// bizTitle, desc, stateType, ns is optional
	CreateApp(title, bizCode, bizTitle, desc, ownerID, language, os, serviceType, ns string, stateType int) (*retailcloud.CreateAppResponse, error)
	CreateDeployConfig(appID int, envType, name, codePath, deployment, statefulSet, cronJob string, configMapList []string) (*retailcloud.CreateDeployConfigResponse, error)
	DeleteDeployConfig(schemaID int) (*retailcloud.DeleteDeployConfigResponse, error)
	GetDeployConfig(appID, id int, name, envType string) (*retailcloud.ListDeployConfigResponse, error)
}

// ClusterService is a cluster service of retail client
type ClusterService interface {
	GetClusterList(envType string, pageNumber, pageSize int) (*retailcloud.ListClusterResponse, error)
	GetCluster(id string) (*retailcloud.QueryClusterDetailResponse, error)
}

// EnvironmentService is a environment service of retail client
type EnvironmentService interface {
	GetEnvList(appID, pageNumber, pageSize, envType int, envName string) (*retailcloud.ListAppEnvironmentResponse, error)
	GetEnvDetail(appID, envID int) (*retailcloud.DescribeAppEnvironmentDetailResponse, error)
	CreateEnv(appID, schemaID, replicas, envType int, envName, region, clusterID string) (*retailcloud.CreateEnvironmentResponse, error)
	DeleteEnv(appID, envID int, force bool) (*retailcloud.DeleteAppEnvironmentResponse, error)
}

// NodeService is a node service of retail client
type NodeService interface {
	GetClusterNodes(clusterID string, pageNumber, pageSize int) (*retailcloud.ListClusterNodeResponse, error)
	GetNodeLabels(clusterID, instanceID, labelKey, labelValue string, pageNumber, pageSize int) (*retailcloud.ListNodeLabelBindingsResponse, error)
	BindNodeLabel(clusterID, instanceID, labelKey, labelValue string) (*retailcloud.BindNodeLabelResponse, error)
	UnbindNodeLabel(clusterID, instanceID, labelKey, labelValue string) (*retailcloud.UnbindNodeLabelResponse, error)
	DeleteNodeLabel(clusterID, labelKey, labelValue string, force bool) (*retailcloud.DeleteNodeLabelResponse, error)
}

// UserService is a user service of retail client
type UserService interface {
	GetUserList(pageNumber, pageSize int) (*retailcloud.ListUsersResponse, error)
}
