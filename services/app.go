package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	"github.com/zhigang/arctl/util"
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

func (s *appServiceImpl) GetAppList(pageNumber, pageSize int) (*retailcloud.ListAppResponse, error) {
	request := retailcloud.CreateListAppRequest()
	request.Scheme = util.RequestScheme

	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListApp(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// TotalPartitions: 发布划分的批次数。
// 建议至少划分为2个批次，划分多个批次的 好处主要有两点：
// 1、做金丝雀灰度，防止本次变更有问题而 导致服务整体宕机；建议在第一批次发布 完成后暂停观察一段时间，若通过系统监 控或查看日志发现变更有问题，则可以停 止发布并回滚；
// 2、分批重启应用实例，防止发布期间所有 实例同时重启而导致服务不可用；
// 具体的错误码。枚举值：102000：没有发布准入凭据； 102001：发布准入凭据审核中
func (s *appServiceImpl) DeployApp(envID, totalPartitions int, name, image string) (*retailcloud.DeployAppResponse, error) {
	request := retailcloud.CreateDeployAppRequest()
	request.Scheme = util.RequestScheme
	request.Name = name
	request.EnvId = requests.NewInteger(envID)
	request.TotalPartitions = requests.NewInteger(totalPartitions)

	request.ContainerImageList = &[]string{image}
	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DeployApp(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) GetAppInstance(appID, envID, pageNumber, pageSize int) (*retailcloud.ListAppInstanceResponse, error) {
	request := retailcloud.CreateListAppInstanceRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.EnvId = requests.NewInteger(envID)
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListAppInstance(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) GetApp(appID int) (*retailcloud.DescribeAppDetailResponse, error) {
	request := retailcloud.CreateDescribeAppDetailRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeAppDetail(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) DeleteApp(appID int, force bool) (*retailcloud.DeleteAppDetailResponse, error) {
	request := retailcloud.CreateDeleteAppDetailRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.Force = requests.NewBoolean(force)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DeleteAppDetail(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) ScaleApp(envID, replicas int) (*retailcloud.ScaleAppResponse, error) {
	request := retailcloud.CreateScaleAppRequest()
	request.Scheme = util.RequestScheme
	request.EnvId = requests.NewInteger(envID)
	request.Replicas = requests.NewInteger(replicas)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ScaleApp(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// CreateApp is create a application
// title, bizCode, ownerID, language, os, serviceType is required
// bizTitle, desc, stateType, ns is optional
func (s *appServiceImpl) CreateApp(title, bizCode, bizTitle, desc, ownerID, language, os, serviceType, ns string, stateType int) (*retailcloud.CreateAppResponse, error) {
	request := retailcloud.CreateCreateAppRequest()
	request.Scheme = util.RequestScheme

	request.BizCode = bizCode
	request.OperatingSystem = os
	request.Language = language
	request.ServiceType = serviceType
	request.Title = title

	var owner retailcloud.CreateAppUserRoles
	// RoleName is one of `Owner`,`PE`,`Dev`,`Test`
	owner.RoleName = "Owner"
	// UserType is one of  of `DING_TALK`,`TAOBAO`
	owner.UserType = "TAOBAO"
	owner.UserId = ownerID
	request.UserRoles = &[]retailcloud.CreateAppUserRoles{owner}

	if bizTitle != "" {
		request.BizTitle = bizTitle
	}
	if desc != "" {
		request.Description = desc
	}
	if stateType > 0 {
		request.StateType = requests.NewInteger(stateType)
	}
	if ns != "" {
		request.Namespace = ns
	}

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.CreateApp(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) DeleteDeployConfig(schemaID int) (*retailcloud.DeleteDeployConfigResponse, error) {
	request := retailcloud.CreateDeleteDeployConfigRequest()
	request.Scheme = util.RequestScheme

	request.SchemaId = requests.NewInteger(schemaID)

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DeleteDeployConfig(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) CreateDeployConfig(appID int, envType, name, codePath, deployment, statefulSet, cronJob string, configMapList []string) (*retailcloud.CreateDeployConfigResponse, error) {
	request := retailcloud.CreateCreateDeployConfigRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)
	request.EnvType = envType
	request.Name = name

	if codePath != "" {
		request.CodePath = codePath
	}

	if deployment != "" {
		request.Deployment = deployment
	}

	if statefulSet != "" {
		request.StatefulSet = statefulSet
	}

	if cronJob != "" {
		request.CronJob = cronJob
	}

	if len(configMapList) > 0 {
		request.ConfigMapList = &configMapList
	}

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.CreateDeployConfig(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *appServiceImpl) GetDeployConfig(appID, id int, name, envType string) (*retailcloud.ListDeployConfigResponse, error) {
	request := retailcloud.CreateListDeployConfigRequest()
	request.Scheme = util.RequestScheme

	request.AppId = requests.NewInteger(appID)

	if id > 0 {
		request.Id = requests.NewInteger(id)
	}

	if name != "" {
		request.Name = name
	}

	if envType != "" {
		request.EnvType = envType
	}

	client, err := s.factory.RetailCloudClient()
	if err != nil {
		return nil, err
	}
	response, err := client.ListDeployConfig(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
