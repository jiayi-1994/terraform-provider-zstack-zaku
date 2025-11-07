package view

import "time"

// ProductInfoView 产品信息视图
type ProductInfoView struct {
	Version string `json:"version"` //产品版本
}

// NamespaceView 命名空间视图
type NamespaceView struct {
	Name       string    `json:"name"`
	ClusterID  int64     `json:"clusterId"`
	ProjectID  int64     `json:"projectId"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"createTime"`
}

// RepositoryView 仓库视图
type RepositoryView struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// ImageView 镜像视图
type ImageView struct {
	Name         string    `json:"name"`
	RepositoryID int64     `json:"repositoryId"`
	Size         int64     `json:"size"`
	TagCount     int       `json:"tagCount"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

// ImageTagView 镜像标签视图
type ImageTagView struct {
	Tag        string    `json:"tag"`
	ImageName  string    `json:"imageName"`
	Size       int64     `json:"size"`
	Digest     string    `json:"digest"`
	CreateTime time.Time `json:"createTime"`
	PushTime   time.Time `json:"pushTime"`
}

// ImageImportApplyView 镜像导入申请视图
type ImageImportApplyView struct {
	UploadURL   string `json:"uploadUrl"`
	UploadToken string `json:"uploadToken"`
	ExpireTime  int64  `json:"expireTime"`
}

// ApiResult API结果
type ApiResult struct {
	Success bool        `json:"success"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Content interface{} `json:"content,omitempty"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// ApiLocation API位置（用于异步操作）
type ApiLocation struct {
	Location string `json:"location"`
	ActionID string `json:"actionId"`
}
