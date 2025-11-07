package client

import (
	"fmt"

	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/view"
)

// ListSimpleCluster 查询集群列表接口
func (cli *ZeClient) ListSimpleCluster() ([]view.ClusterSimpleItem, error) {
	resp := make([]view.ClusterSimpleItem, 0)
	return resp, cli.Get("/open-api/v1/cloud/clusters", "", nil, &resp)
}

// CreateProject 创建项目接口
func (cli *ZeClient) CreateProject(params param.CloudProjectCreateParam) (*view.ProjectView, error) {
	var resp view.ProjectView
	return &resp, cli.Post("/open-api/v1/cloud/projects", params, &resp)
}

// UpdateProject 编辑项目接口
func (cli *ZeClient) UpdateProject(params param.CloudProjectUpdateParam) (*view.ProjectView, error) {
	var resp view.ProjectView
	return &resp, cli.Put("/open-api/v1/cloud/projects", "", params, &resp)
}

// DeleteProject 删除项目接口
func (cli *ZeClient) DeleteProject(uuid string) error {
	return cli.Delete("/open-api/v1/cloud/projects", uuid, "")
}

// AddProjectUser 项目添加用户接口
func (cli *ZeClient) AddProjectUser(params param.CloudProjectUserAddParam) error {
	return cli.Post("/open-api/v1/cloud/projects/users", params, nil)
}

// RemoveProjectUser 项目移除用户接口
func (cli *ZeClient) RemoveProjectUser(projectUuids, usernames []string) error {
	paramsStr := fmt.Sprintf("projectUuids=%s&usernames=%s",
		joinStrings(projectUuids), joinStrings(usernames))
	return cli.DeleteWithSpec("/open-api/v1/cloud/projects/users", "", "", paramsStr, nil)
}

// GetProjectQuota 查询项目资源配额情况
func (cli *ZeClient) GetProjectQuota(projectUuid string, clusterId int) ([]view.CloudResourceQuotaView, error) {
	var resp []view.CloudResourceQuotaView
	path := fmt.Sprintf("/open-api/v1/cloud/projects/%s/clusters/%d/quota", projectUuid, clusterId)
	return resp, cli.Get(path, "", nil, &resp)
}

// UpdateProjectResourceQuota 更新项目资源配额情况
func (cli *ZeClient) UpdateProjectResourceQuota(projectUuid string, clusterId int, params param.ProjectQuotaParam) ([]view.CloudResourceQuotaView, error) {
	var resp []view.CloudResourceQuotaView
	path := fmt.Sprintf("/open-api/v1/cloud/projects/%s/clusters/%d/quota", projectUuid, clusterId)
	return resp, cli.Put(path, "", params, &resp)
}

// CreateUser 创建用户接口
func (cli *ZeClient) CreateUser(params param.CloudUserCreateParam) (*view.UserView, error) {
	var resp view.UserView
	return &resp, cli.Post("/open-api/v1/cloud/users", params, &resp)
}

// DeleteUser 删除用户接口
func (cli *ZeClient) DeleteUser(name string) error {
	return cli.Delete("/open-api/v1/cloud/users", name, "")
}

// SetAdministrator 用户设为管理员接口
func (cli *ZeClient) SetAdministrator(name string) error {
	return cli.PutWithSpec("/open-api/v1/cloud/users", name, "setadmin", "", nil, nil)
}

// UnsetAdministrator 用户取消管理员接口
func (cli *ZeClient) UnsetAdministrator(name string) error {
	return cli.PutWithSpec("/open-api/v1/cloud/users", name, "unsetadmin", "", nil, nil)
}

// joinStrings 辅助函数：将字符串数组用逗号连接
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}
