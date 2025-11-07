package view

import "time"

// ProjectView 项目视图
type ProjectView struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

// UserView 用户视图
type UserView struct {
	Name       string    `json:"name"`
	Password   string    `json:"password,omitempty"`
	IsAdmin    bool      `json:"isAdmin"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// CloudResourceQuotaView 云资源配额视图
type CloudResourceQuotaView struct {
	ResourceType string  `json:"resourceType"`
	Quota        float64 `json:"quota"`
	Used         float64 `json:"used"`
	Available    float64 `json:"available"`
	Unit         string  `json:"unit"`
}

// ClusterSimpleView 简单集群视图

// ClusterSimpleItem 简单集群项
type ClusterSimpleItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	NodeCount  int    `json:"nodeCount"`
	CreateTime string `json:"createTime"`
}

// ClusterDetailsView 集群详情视图
type ClusterDetailsView struct {
	BaseClusterView
	ClusterStatus
	ClusterUsage
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

// ClusterInventoryView 集群清单视图
type ClusterInventoryView struct {
	ID         int64                  `json:"id"`
	Name       string                 `json:"name"`
	Status     string                 `json:"status"`
	CreateTime time.Time              `json:"createTime"`
	Inventory  map[string]interface{} `json:"inventory"`
}

// ClusterConfigView 集群配置视图
type ClusterConfigView struct {
	Kubeconfig string `json:"kubeconfig"`
	Config     string `json:"config"`
}

// ClusterOperationView 集群操作视图
type ClusterOperationView struct {
	ID          int64     `json:"id"`
	ClusterID   int64     `json:"clusterId"`
	Operation   string    `json:"operation"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	CreateTime  time.Time `json:"createTime"`
	FinishTime  time.Time `json:"finishTime"`
	OperateUser string    `json:"operateUser"`
}

// ClusterPage 集群分页
type ClusterPage struct {
	TotalCount int `json:"totalCount"`
}

// ClusterOperationPage 集群操作分页
type ClusterOperationPage struct {
	TotalCount int `json:"totalCount"`
}
