package client

import (
	"fmt"

	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/util/utils"
	"zstack.io/edge-go-sdk/pkg/view"
)

// PageNode 节点列表
func (cli *ZeClient) PageNode(clusterId int, params param.QueryParam) ([]view.NodeView, int, error) {
	var resp []view.NodeView
	path := fmt.Sprintf("/open-api/v1/cluster/%d/node", clusterId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// AddNode 添加节点
func (cli *ZeClient) AddNode(clusterId int, params param.NodeAddParamOpenApi, async bool) (string, error) {
	params.Password = utils.EncryptByAccessKey(cli.ZeConfig.GetAccessKeySecret(), params.Password)
	path := fmt.Sprintf("/open-api/v1/cluster/%d/node", clusterId)
	cli.retryTimes = 300
	return cli.PostWithAsync(path, "", "", "", params, nil, async)
}

// DeleteNode 删除节点
func (cli *ZeClient) DeleteNode(clusterId int, nodenames []string) error {
	path := fmt.Sprintf("/open-api/v1/cluster/%d/node", clusterId)
	paramsStr := fmt.Sprintf("nodenames=%s", joinStrings(nodenames))
	cli.retryTimes = 300
	return cli.DeleteWithSpec(path, "", "", paramsStr, nil)
}

// GetNodeDisk 获取节点磁盘(非系统盘)
func (cli *ZeClient) GetNodeDisk(sshIP string, sshPort int, sshPassword string) ([]view.DiskInfoView, error) {
	pwd := utils.EncryptByAccessKey(cli.ZeConfig.GetAccessKeySecret(), sshPassword)
	var resp []view.DiskInfoView
	params := struct {
		SshIP       string `json:"sshIP"`
		SshPort     int    `json:"sshPort"`
		SshPassword string `json:"sshPassword"`
	}{
		SshIP:       sshIP,
		SshPort:     sshPort,
		SshPassword: pwd,
	}
	return resp, cli.Get("/open-api/v1/cluster/node-disk", "", &params, &resp)
}
