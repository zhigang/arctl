package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

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
	request := retailcloud.CreateListClusterRequest()
	request.Scheme = util.RequestScheme
	if envType != "" {
		request.EnvType = envType
	}
	request.PageNum = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListCluster(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *clusterServiceImpl) GetCluster(id string) (*retailcloud.QueryClusterDetailResponse, error) {
	request := retailcloud.CreateQueryClusterDetailRequest()
	request.Scheme = util.RequestScheme
	request.ClusterInstanceId = id

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}

	response, err := client.QueryClusterDetail(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
