package param

type ExternalNetworkIpPoolType string
type NetworkType string
type ShareType string

const (
	IpRange NetworkType = "IpRange"
	Cidr    NetworkType = "Cidr"

	ShareAssign                  ShareType                 = "assign"
	ShareGlobal                  ShareType                 = "global"
	ExternalNetworkIpPoolTypeSvc ExternalNetworkIpPoolType = "Svc" //服务外部网络
	ExternalNetworkIpPoolTypePod ExternalNetworkIpPoolType = "Pod" //容器组附加网络
)

type ExternalNetworkCreateParam struct {
	ClusterID   int    `json:"clusterID"`
	Description string `json:"description"`
	Gateway     string `json:"gateway"`
	Iface       string `json:"iface"`
	Name        string `json:"name"`
	Netmask     string `json:"netmask"`
}

// ExternalNetworkIpPoolCreateParam 创建外部网络IP池参数
type ExternalNetworkCreateIpPoolParam struct {
	ExternalNetworkIpPoolParam
	StartIp string `json:"startIp"  validate:"required"`
	EndIp   string `json:"endIp"  validate:"required"`
}

type ExternalNetworkIpPoolParam struct {
	Name       string                    `json:"name" validate:"min=2,max=50,regexp=^[a-z][-a-z0-9]*[a-z0-9]$"  validate:"required"`
	IpPoolType ExternalNetworkIpPoolType `json:"ipPoolType" `
}

// ExternalNetworkIpPoolUpdateParam 更新外部网络IP池参数
type ExternalNetworkIpPoolUpdateParam struct {
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	StartIP string `json:"startIp,omitempty"`
	EndIP   string `json:"endIp,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Netmask string `json:"netmask,omitempty"`
}
