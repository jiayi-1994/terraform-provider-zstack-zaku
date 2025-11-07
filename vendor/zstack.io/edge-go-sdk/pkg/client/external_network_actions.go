package client

import (
	"fmt"

	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/view"
)

func (cli *ZeClient) CreateExternalNetwork(params param.ExternalNetworkCreateParam) (string, error) {
	var resp string
	return resp, cli.Post(fmt.Sprintf("/open-api/v1/external-network/%d", params.ClusterID), params, &resp)
}

// PageExternalNetwork 查询外部网络列表
func (cli *ZeClient) PageExternalNetwork(clusterId int, params param.QueryParam) ([]view.ExternalNetworkView, int, error) {
	var resp []view.ExternalNetworkView
	path := fmt.Sprintf("/open-api/v1/external-network/%d", clusterId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// GetExternalNetworkCandidateInterface 查询外部网络列表可选的网卡
func (cli *ZeClient) GetExternalNetworkCandidateInterface(clusterId int, refresh bool) (*view.NodeIfaceAllResult, error) {
	var resp view.NodeIfaceAllResult
	path := fmt.Sprintf("/open-api/v1/external-network/%d/candidate-interface", clusterId)
	params := struct {
		Refresh bool `json:"refresh"`
	}{Refresh: refresh}
	return &resp, cli.Get(path, "", &params, &resp)
}

// PageExternalNetworkIpPool 查询外部网络Ip池列表
func (cli *ZeClient) PageExternalNetworkIpPool(clusterId, externalNetworkId int, params param.QueryParam) ([]view.ExternalNetworkIpPoolView, int, error) {
	var resp []view.ExternalNetworkIpPoolView
	path := fmt.Sprintf("/open-api/v1/external-network/%d/%d", clusterId, externalNetworkId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// CreateExternalNetworkIpPool 创建外部网络Ip池
func (cli *ZeClient) CreateExternalNetworkIpPool(clusterId, externalNetworkId int, params param.ExternalNetworkCreateIpPoolParam) error {
	path := fmt.Sprintf("/open-api/v1/external-network/%d/%d", clusterId, externalNetworkId)
	return cli.Post(path, params, nil)
}

// UpdateExternalNetworkIpPool 更新外部网络Ip池
//func (cli *ZeClient) UpdateExternalNetworkIpPool(clusterId, externalNetworkId int, params param.ExternalNetworkIpPoolUpdateParam) error {
//	path := fmt.Sprintf("/open-api/v1/external-network/%d/%d", clusterId, externalNetworkId)
//	return cli.Put(path, "", params, nil)
//}

// DeleteExternalNetworkIpPool 删除外部网络Ip池
//func (cli *ZeClient) DeleteExternalNetworkIpPool(clusterId, externalNetworkId, ipPoolId int) error {
//	path := fmt.Sprintf("/open-api/v1/external-network/%d/%d", clusterId, externalNetworkId)
//	return cli.Delete(path, fmt.Sprintf("%d", ipPoolId), "")
//}

// GetExternalNetworkAvailableIps 外部网络可用的ip
func (cli *ZeClient) GetExternalNetworkAvailableIps(clusterId int, networkName string) ([]string, error) {
	var resp []string
	path := fmt.Sprintf("/open-api/v1/external-network/%d/%s/availableIps", clusterId, networkName)
	return resp, cli.Get(path, "", nil, &resp)
}

// PageExternalNetworkForSvc 查询可选服务外部网络
func (cli *ZeClient) PageExternalNetworkForSvc(clusterId, projectId int, params param.QueryParam) ([]view.ExternalNetworkView, int, error) {
	var resp []view.ExternalNetworkView
	path := fmt.Sprintf("/open-api/v1/external-network/%d/%d/forSvc", clusterId, projectId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}
