package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

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
	request := retailcloud.CreateListAppEnvironmentRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)
	if envType >= 0 {
		request.EnvType = requests.NewInteger(envType)
	}
	if envName != "" {
		request.EnvName = envName
	}
	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}

	response, err := client.ListAppEnvironment(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *environmentServiceImpl) GetEnvDetail(appID, envID int) (*retailcloud.DescribeAppEnvironmentDetailResponse, error) {
	request := retailcloud.CreateDescribeAppEnvironmentDetailRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.EnvId = requests.NewInteger(envID)
	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeAppEnvironmentDetail(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *environmentServiceImpl) CreateEnv(appID, schemaID, replicas, envType int, envName, region, clusterID string) (*retailcloud.CreateEnvironmentResponse, error) {
	request := retailcloud.CreateCreateEnvironmentRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.AppSchemaId = requests.NewInteger(schemaID)
	request.EnvName = envName
	request.EnvType = requests.NewInteger(envType)
	request.Replicas = requests.NewInteger(replicas)
	request.Region = region
	request.ClusterId = clusterID

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.CreateEnvironment(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *environmentServiceImpl) DeleteEnv(appID, envID int, force bool) (*retailcloud.DeleteAppEnvironmentResponse, error) {
	request := retailcloud.CreateDeleteAppEnvironmentRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.EnvId = requests.NewInteger(envID)
	request.Force = requests.NewBoolean(force)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DeleteAppEnvironment(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
