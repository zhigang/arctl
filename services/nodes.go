package services

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

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
	request := retailcloud.CreateListClusterNodeRequest()
	request.Scheme = util.REQUEST_SCHEME

	if clusterID != "" {
		request.ClusterInstanceId = clusterID
	}
	request.PageNum = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)
	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListClusterNode(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *nodeServiceImpl) GetNodeLabels(clusterID, instanceID, labelKey, labelValue string, pageNumber, pageSize int) (*retailcloud.ListNodeLabelBindingsResponse, error) {
	request := retailcloud.CreateListNodeLabelBindingsRequest()
	request.Scheme = util.REQUEST_SCHEME

	if instanceID != "" {
		request.InstanceId = instanceID
	}

	if clusterID != "" {
		request.ClusterId = clusterID
	}
	if labelKey != "" && !strings.HasPrefix(labelKey, util.LABEL_PREFIX_JST) {
		request.LabelKey = util.LABEL_PREFIX_JST + labelKey
	}
	if labelValue != "" {
		request.LabelValue = labelValue
	}

	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListNodeLabelBindings(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *nodeServiceImpl) BindNodeLabel(clusterID, instanceID, labelKey, labelValue string) (*retailcloud.BindNodeLabelResponse, error) {
	request := retailcloud.CreateBindNodeLabelRequest()
	request.Scheme = util.REQUEST_SCHEME
	request.ClusterId = clusterID
	request.InstanceId = instanceID
	request.LabelKey = labelKey
	request.LabelValue = labelValue

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.BindNodeLabel(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *nodeServiceImpl) UnbindNodeLabel(clusterID, instanceID, labelKey, labelValue string) (*retailcloud.UnbindNodeLabelResponse, error) {
	request := retailcloud.CreateUnbindNodeLabelRequest()
	request.Scheme = util.REQUEST_SCHEME
	request.ClusterId = clusterID
	request.InstanceId = instanceID
	request.LabelKey = labelKey
	request.LabelValue = labelValue

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.UnbindNodeLabel(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *nodeServiceImpl) DeleteNodeLabel(clusterID, labelKey, labelValue string, force bool) (*retailcloud.DeleteNodeLabelResponse, error) {
	request := retailcloud.CreateDeleteNodeLabelRequest()
	request.Scheme = util.REQUEST_SCHEME
	request.ClusterId = clusterID
	request.LabelKey = labelKey
	request.LabelValue = labelValue
	request.Force = requests.NewBoolean(force)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DeleteNodeLabel(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
