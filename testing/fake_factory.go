package testing

import (
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

var fct *factoryImpl

type factoryImpl struct {
	util.Factory
	cfgFile string
	config  *util.Config
	loader  sync.Once
}

// NewFactory returns a fake impl factory
func NewFactory() util.Factory {
	fct = &factoryImpl{}
	return fct
}

// SetConfigFilePath sets a config file path
func SetConfigFilePath(cfgFile string) {
	fct.cfgFile = cfgFile
}

func (f *factoryImpl) Config() (*util.Config, error) {
	var err error = nil
	f.loader.Do(func() {
		f.config, err = util.LoadConfig(f.cfgFile)
	})
	return f.config, err
}

// RetailCloudClient is return a retailcloud client.
func (f *factoryImpl) RetailCloudClient() (*retailcloud.Client, error) {
	client, err := retailcloud.NewClientWithAccessKey(f.config.Aksk.RegionID, f.config.Aksk.AccessKeyID, f.config.Aksk.AccessKeySecret)
	return client, err
}

func (f *factoryImpl) AppService() util.AppService {
	return newAppService(f)
}

func (f *factoryImpl) ClusterService() util.ClusterService {
	return newClusterService(f)
}

func (f *factoryImpl) EnvironmentService() util.EnvironmentService {
	return newEnvironmentService(f)
}

func (f *factoryImpl) NodeService() util.NodeService {
	return newNodeService(f)
}

func (f *factoryImpl) UserService() util.UserService {
	return newUserService(f)
}
