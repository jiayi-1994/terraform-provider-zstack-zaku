package view

import "time"

type ClusterCreateType string

type BaseClusterView struct {
	ID            int64             `json:"id"`            //ID
	Name          string            `json:"name"`          //名称
	CreateTime    time.Time         `json:"createTime"`    //创建时间
	PrometheusURL string            `json:"prometheusURL"` //prometheus地址
	CreateType    ClusterCreateType `json:"createType"`    //创建类型
}

type ClusterStatus struct {
	Status                   string `json:"status"`                   //当前集群状态
	Version                  string `json:"version"`                  //集群版本
	PlatformComponentVersion string `json:"platformComponentVersion"` //平台组件版本
	NodeCount                int    `json:"nodeCount"`                //集群节点数
}

type ClusterUsage struct {
	Cpu     string `json:"cpu"`     //cpu情况
	Memory  string `json:"memory"`  //内存情况
	Storage string `json:"storage"` //存储情况
}

type ClusterView struct {
	BaseClusterView

	ClusterStatus

	ClusterUsage
}

type ClusterSpinnerView struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	Unhealthy  bool              `json:"unhealthy"`
	CreateType ClusterCreateType `json:"createType"` // Inner,Outer
	Version    string            `json:"version"`    //集群版本
}
