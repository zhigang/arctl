package util

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
)

// Factory provides abstractions that allow the arctl command to be extended across multiple types
// of resources and different API sets.
type Factory interface {
	Config() (*Config, error)
	RetailCloudClient() (*retailcloud.Client, error)
	AppService() AppService
	ClusterService() ClusterService
	EnvironmentService() EnvironmentService
	NodeService() NodeService
	UserService() UserService
}
