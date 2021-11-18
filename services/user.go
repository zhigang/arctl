package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
)

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
	request := retailcloud.CreateListUsersRequest()
	request.Scheme = util.REQUEST_SCHEME
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListUsers(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
