package client

import (
	"fmt"

	"zstack.io/edge-go-sdk/pkg/param"
	"zstack.io/edge-go-sdk/pkg/view"
)

// GetProductInfo 获取产品信息
func (cli *ZeClient) GetProductInfo() (*view.ProductInfoView, error) {
	var resp view.ProductInfoView
	return &resp, cli.Get("/open-api/v1/product-info", "", nil, &resp)
}

// PageProjectCluster 查询项目下的集群接口
func (cli *ZeClient) ListProjectCluster(projectId int) ([]view.ClusterSpinnerView, error) {
	var resp []view.ClusterSpinnerView
	path := fmt.Sprintf("/open-api/v1/project/%d/cluster", projectId)
	err := cli.Get(path, "", nil, &resp)
	return resp, err
}

// ListProjectNamespace 查询集群中的命名空间接口
func (cli *ZeClient) ListProjectNamespace(projectId, clusterId int) ([]view.NamespaceView, error) {
	var resp []view.NamespaceView
	path := fmt.Sprintf("/open-api/v1/project/%d/cluster/%d/namespace", projectId, clusterId)
	return resp, cli.Get(path, "", nil, &resp)
}

// SaveContainerAsImage 容器保存为镜像接口
func (cli *ZeClient) SaveContainerAsImage(projectId, clusterId int, namespace, podName, containerName string, params param.ContainerPackageParam) error {
	path := fmt.Sprintf("/open-api/v1/project/%d/cluster/%d/namespace/%s/pod/%s/container/%s/package",
		projectId, clusterId, namespace, podName, containerName)
	return cli.Post(path, params, nil)
}

// PageProjectRepository 查询项目下的本地仓库接口
func (cli *ZeClient) PageProjectRepository(projectId int, params param.QueryParam) ([]view.RepositoryView, int, error) {
	var resp []view.RepositoryView
	path := fmt.Sprintf("/open-api/v1/project/%d/repository", projectId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// CreateProjectRepository 创建本地仓库接口
func (cli *ZeClient) CreateProjectRepository(projectId int, params param.RepositoryCreateParam) (*view.RepositoryView, error) {
	var resp view.RepositoryView
	path := fmt.Sprintf("/open-api/v1/project/%d/repository", projectId)
	return &resp, cli.Post(path, params, &resp)
}

//// UpdateProjectRepository 更新本地仓库接口
//func (cli *ZeClient) UpdateProjectRepository(projectId int, params param.RepositoryUpdateParam) (*view.RepositoryView, error) {
//	var resp view.RepositoryView
//	path := fmt.Sprintf("/open-api/v1/project/%d/repository", projectId)
//	return &resp, cli.Put(path, "", params, &resp)
//}

// DeleteProjectRepository 删除本地仓库接口
//func (cli *ZeClient) DeleteProjectRepository(projectId, repositoryId int) error {
//	path := fmt.Sprintf("/open-api/v1/project/%d/repository", projectId)
//	return cli.Delete(path, fmt.Sprintf("%d", repositoryId), "")
//}

// PageRepositoryImage 查询本地仓库里的镜像接口
func (cli *ZeClient) PageRepositoryImage(projectId, repositoryId int, params param.QueryParam) ([]view.ImageView, int, error) {
	var resp []view.ImageView
	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%d/image", projectId, repositoryId)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// DeleteRepositoryImage 删除本地仓库里的镜像接口
//func (cli *ZeClient) DeleteRepositoryImage(projectId, repositoryId int, imageName string) error {
//	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%d/image", projectId, repositoryId)
//	return cli.Delete(path, imageName, "")
//}

// PageRepositoryImageTag 查询本地仓库里的镜像版本接口
func (cli *ZeClient) PageRepositoryImageTag(projectId, repositoryId int, imageName string, params param.QueryParam) ([]view.ImageTagView, int, error) {
	var resp []view.ImageTagView
	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%d/image/%s/tag", projectId, repositoryId, imageName)
	total, err := cli.Page(path, &params, &resp)
	return resp, total, err
}

// DeleteRepositoryImageTag 删除本地仓库里的镜像版本接口
//func (cli *ZeClient) DeleteRepositoryImageTag(projectId, repositoryId int, imageName, tag string) error {
//	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%d/image/%s/tag", projectId, repositoryId, imageName)
//	return cli.Delete(path, tag, "")
//}

// ImportImageByUrl 在线上传镜像接口
func (cli *ZeClient) ImportImageByUrl(projectId int, repositoryName string, params param.ImageUrlImportParam) error {
	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%s/image/url-import", projectId, repositoryName)
	return cli.Post(path, params, nil)
}

// ApplyBeforeImportImage 离线上传镜像前申请接口
func (cli *ZeClient) ApplyBeforeImportImage(projectId int, repositoryName, imageFileName string, params param.ImageImportApplyParam) (*view.ImageImportApplyView, error) {
	var resp view.ImageImportApplyView
	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%s/image/%s/apply-before-import",
		projectId, repositoryName, imageFileName)
	return &resp, cli.Post(path, params, &resp)
}

// ConfirmAfterImportImage 离线上传镜像后确认接口
func (cli *ZeClient) ConfirmAfterImportImage(projectId int, repositoryName, imageFileName string, params param.ImageImportConfirmParam) error {
	path := fmt.Sprintf("/open-api/v1/project/%d/repository/%s/image/%s/confirm-after-import",
		projectId, repositoryName, imageFileName)
	return cli.Post(path, params, nil)
}
