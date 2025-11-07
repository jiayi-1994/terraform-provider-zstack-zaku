package view

import (
	"time"

	"zstack.io/edge-go-sdk/pkg/param"
)

type NetWorkType string

const (
	NetTypeManager  NetWorkType = "manager"
	NetTypeBusiness NetWorkType = "business"
)

// ExternalNetworkView 外部网络视图
type ExternalNetworkView struct {
	IpTotalNum      int64       `json:"ipTotalNum"`
	IpUsedNum       int64       `json:"ipUsedNum"`
	ExistNetwork    bool        `json:"existNetwork"`
	ID              int64       `json:"id"` // 主键ID
	ClusterID       int64       `json:"cluster_id" `
	Name            string      `json:"name"` // 网络名称
	Description     string      `json:"description" `
	Iface           string      `json:"iface"` // 网卡名称
	Type            NetWorkType `json:"type"`
	SpiderPoolReady bool        `json:"spiderPoolReady"`
	MetallbReady    bool        `json:"metallbReady"`
	Netmask         string      `json:"netmask"` // 子网掩码
	Gateway         string      `json:"gateway"` // 网关
	Cidr            string      `json:"cidr"`    // ipv4 CIDR
	StartIp         uint        `json:"-"`       // 起始IP
	EndIp           uint        `json:"-"`       // 结束IP
	CreateTime      time.Time   `json:"createTime"`
}

// ExternalNetworkPage 外部网络分页
type ExternalNetworkPage struct {
	TotalCount int `json:"totalCount"`
}

// ExternalNetworkIpPoolView 外部网络IP池视图
type ExternalNetworkIpPoolView struct {
	Name        string                          `json:"name"`
	L2Name      string                          `json:"l2Name"`
	Type        param.ExternalNetworkIpPoolType `json:"type"`
	ShareType   param.ShareType                 `json:"shareType"`  // assign global
	ProjectIDs  []int64                         `json:"projectIDs"` //项目ID
	IpRanges    []string                        `json:"ipRanges"`
	Ip6Ranges   []string                        `json:"ip6Ranges"`
	IpTotalNum  int64                           `json:"ipTotalNum"`
	IpUsedNum   int64                           `json:"ipUsedNum"`
	CreateTime  time.Time                       `json:"createTime"`
	ExistIpUsed bool                            `json:"existIpUsed"`
	Disabled    bool                            `json:"disabled"`
}

type NodeIfaceWithRouteIface struct {
	Host          string       `json:"host"`
	Name          string       `json:"name"`
	IsMaster      bool         `json:"isMaster"`
	Ifaces        []string     `json:"ifaces"`
	IfaceWithVlan []string     `json:"ifaceWithVlan"`
	NodeIfaceMap  NodeIfaceMap `json:"nodeIfaceMap"`
}
type NodeIfaceMap struct {
	VlanMap  map[string]string
	BrMap    map[string]string
	RouteMap map[string]NodeInterfaceInfo
}
type NodeInterfaceInfo struct {
	Interface string `json:"interface"`
	IpRange   string `json:"ipRange"`
}

type NodeIfaceAllResult struct {
	SameIfaces       []string                  `json:"sameIfaces"`
	MasterSameIfaces []string                  `json:"masterSameIfaces"`
	NodeIfaces       []NodeIfaceWithRouteIface `json:"nodeIfaces"`
}

// NodeIfaceAllResult 节点网卡结果

// NodeIfaceInfo 节点网卡信息
type NodeIfaceInfo struct {
	NodeName   string      `json:"nodeName"`
	Interfaces []IfaceInfo `json:"interfaces"`
}

// IfaceInfo 网卡信息
type IfaceInfo struct {
	Name       string   `json:"name"`
	IP         string   `json:"ip"`
	MAC        string   `json:"mac"`
	Status     string   `json:"status"`
	Speed      string   `json:"speed"`
	Duplex     string   `json:"duplex"`
	MTU        int      `json:"mtu"`
	Slaves     []string `json:"slaves,omitempty"`
	BondMode   string   `json:"bondMode,omitempty"`
	BondStatus string   `json:"bondStatus,omitempty"`
}
