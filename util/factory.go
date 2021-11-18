package util

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
)

type Factory interface {
	Config() (*Config, error)
	RetailCloudClient() (*retailcloud.Client, error)
	AppService() AppService
	ClusterService() ClusterService
	EnvironmentService() EnvironmentService
	NodeService() NodeService
	UserService() UserService
}
