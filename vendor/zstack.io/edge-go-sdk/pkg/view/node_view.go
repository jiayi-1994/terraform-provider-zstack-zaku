package view

import "time"

// NodeView 节点视图
type NodeView struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ClusterID  int64     `json:"clusterId"`
	IP         string    `json:"ip"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	CPU        string    `json:"cpu"`
	Memory     string    `json:"memory"`
	Storage    string    `json:"storage"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// NodePage 节点分页
type NodePage struct {
	TotalCount int `json:"totalCount"`
}

// DiskInfoView 磁盘信息视图
type DiskInfoView struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	SizeStr string `json:"sizeStr"`
	Used    bool   `json:"used"` // 是否已经使用
}
