package param

const (
	// Harbor类型
	TypeHarborProjectPublic  = "Type_Harbor_Project_Public"
	TypeHarborProjectPrivate = "Type_Harbor_Project_Private"
)

// ContainerPackageParam 容器打包参数
type ContainerPackageParam struct {
	Repository string `json:"repository" example:"dev-test"` //仓库
	Name       string `json:"name" example:"aaa"`            //镜像名称
	Tag        string `json:"tag"`                           //标签
}

// RepositoryCreateParam 创建仓库参数
type RepositoryCreateParam struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type" example:"Type_Harbor_Project_Public" validate:"regexp=^(Type_Harbor_Project_Public)|(Type_Harbor_Project_Private)$"` // Type_Harbor_Project_Public Type_Harbor_Project_Private
}

// RepositoryUpdateParam 更新仓库参数
type RepositoryUpdateParam struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ImageUrlImportParam 在线导入镜像参数
type ImageUrlImportParam struct {
	SourceRepo string `json:"sourceRepo" example:"docker.io/nginx:v1.0"`                 //镜像地址
	Username   string `json:"username" example:"YWRtaW4=(admin)(base64 encoded)"`        //镜像仓库用户（需要base64）
	Password   string `json:"password" example:"cGFzc3dvcmQ=(password)(base64 encoded)"` //镜像仓库用户密码（需要base64）
	Platforms  string `json:"platforms" example:"linux/amd64,linux/arm64/v8"`            //指定镜像拉取的平台，多平台用,分隔，比如linux/amd64,linux/arm64/v8
}

// ImageImportApplyParam 离线导入镜像申请参数
type ImageImportApplyParam struct {
	ImageName string `json:"imageName"`
	Tag       string `json:"tag"`
	Size      int64  `json:"size"`
}

// ImageImportConfirmParam 离线导入镜像确认参数
type ImageImportConfirmParam struct {
	ImageName string `json:"imageName"`
	Tag       string `json:"tag"`
	MD5       string `json:"md5,omitempty"`
}
