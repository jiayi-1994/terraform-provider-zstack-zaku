package client

import (
	"fmt"

	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/util/utils"
	"zstack.io/edge-go-sdk/pkg/view"
)

func (cli *ZeClient) PageCluster(params param.QueryParam) ([]view.ClusterView, int, error) {
	var resp []view.ClusterView
	total, err := cli.Page("/open-api/v1/cluster", &params, &resp)
	return resp, total, err
}

func (cli *ZeClient) GetCluster(resourceId string) (*view.ClusterView, error) {
	var resp view.ClusterView
	return &resp, cli.Get("/open-api/v1/cluster", resourceId, nil, &resp)
}

// CreateCluster 创建集群
func (cli *ZeClient) CreateCluster(params param.ClusterCreateParam, async bool) (string, error) {
	params.Password = utils.EncryptByAccessKey(cli.ZeConfig.GetAccessKeySecret(), params.Password)
	cli.retryInterval = 10
	cli.retryTimes = 500
	return cli.PostWithAsync("/open-api/v1/cluster", "", "", "", params, nil, async)
}

// DeleteCluster 删除集群
func (cli *ZeClient) DeleteCluster(clusterId int, async bool) (string, error) {
	cli.retryInterval = 10
	cli.retryTimes = 500
	return cli.DeleteWithAsync("/open-api/v1/cluster", fmt.Sprintf("%d", clusterId), "", "", nil, async)
}

// GetClusterDetails 集群详情
func (cli *ZeClient) GetClusterDetails(clusterId int) (*view.ClusterDetailsView, error) {
	var resp view.ClusterDetailsView
	return &resp, cli.Get("/open-api/v1/cluster", fmt.Sprintf("%d", clusterId), nil, &resp)
}

// HasIluvatarLicense 判断集群是否有天数许可证
func (cli *ZeClient) HasIluvatarLicense(clusterId int) (bool, error) {
	var resp bool
	path := fmt.Sprintf("/open-api/v1/cluster/%d/has-iluvatar-license", clusterId)
	return resp, cli.Get(path, "", nil, &resp)
}

// GetClusterKubeconfig 获取集群kubeconfig接口
func (cli *ZeClient) GetClusterKubeconfig(clusterId int) (*view.ClusterConfigView, error) {
	var resp view.ClusterConfigView
	path := fmt.Sprintf("/open-api/v1/cluster/%d/kubeconfig", clusterId)
	return &resp, cli.Get(path, "", nil, &resp)
}

// GetClusterOperationLog 集群操作日志详情
func (cli *ZeClient) GetClusterOperationLog(clusterId, logId int) (string, error) {
	var resp string
	path := fmt.Sprintf("/open-api/v1/cluster/%d/log/%d", clusterId, logId)
	return resp, cli.Get(path, "", nil, &resp)
}

// PageClusterOperation 集群操作日志列表
func (cli *ZeClient) PageClusterOperation(clusterId int, params param.QueryParam) ([]view.ClusterOperationView, int, error) {
	var resp []view.ClusterOperationView
	path := fmt.Sprintf("/open-api/v1/cluster/%d/operation/list", clusterId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// RecreateCluster 重新安装集群
func (cli *ZeClient) RecreateCluster(clusterId int, async bool) (string, error) {
	cli.retryInterval = 10
	cli.retryTimes = 500
	path := fmt.Sprintf("/open-api/v1/cluster/%d/recreate", clusterId)
	return cli.PostWithAsync(path, "", "", "", nil, nil, async)
}
