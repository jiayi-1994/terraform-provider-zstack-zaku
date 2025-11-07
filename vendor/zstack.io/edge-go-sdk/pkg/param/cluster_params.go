package param

// ClusterCreateParam 创建集群参数
type ClusterCreateParam struct {
	EnableHA         bool                     `json:"enableHA" example:"true"` // 是否开启高可用
	Name             string                   `json:"name" example:"cluster-1" validate:"required"`
	Nodes            []ClusterCreateNodeParam `json:"nodes" validate:"dive,required"`                                                                                                                                                                                      // nodes for creating cluster
	NetCombined      bool                     `json:"netCombined" example:"true" `                                                                                                                                                                                         // 管理业务网是否复用
	Port             int                      `json:"port" example:"22" validate:"required"`                                                                                                                                                                               // ssh port
	Password         string                   `json:"password" validate:"required" example:"z3ODfYnKkhwXb9zgzy6F5hCYDI3j8/RWmiJpdQeO8h5GfOZnlefB/57I72AtsdPbJhssLwgp92yRKSPnkHY8GVbi4Zh0Y3Wccx8HMig+i2pGKDDMM5u0uE/IUKoKZ9kLjyKyD039s6sE2kPoRXcEJ+vIsa1OxGJ9vrahVM1MwoU="` // password should be encrypted
	ManagementVipV4  string                   `json:"managementVipV4" validate:"required,ipv4" example:"172.31.13.100"`                                                                                                                                                    // 管理网络cidr v4
	BusinessVipV4    string                   `json:"businessVipV4" validate:"required,ipv4" example:"172.32.4.100"`                                                                                                                                                       // 业务网络cidr v4
	MaxPodPerNode    int                      `json:"maxPodPerNode" example:"110"`                                                                                                                                                                                         // 节点最大pod数量
	DataDisk         map[string][]string      `json:"dataDisk" validate:"required"`
	ImageDataDisk    map[string][]string      `json:"imageDataDisk" `                                                 // 镜像数据盘
	PodCidrV4        string                   `json:"podCidrV4" example:"10.233.64.0/18" validate:"required,cidr"`    // k8s pod cidr v4
	ServiceCidrV4    string                   `json:"serviceCidrV4" example:"10.233.0.0/18" validate:"required,cidr"` // k8s service cidr v4
	DNSServer        string                   `json:"dnsServer" example:"223.5.5.5" validate:"required,ipv4"`         //single ip
	IstioEnabled     bool                     `json:"enableIstio" example:"true"`                                     // 是否开启istio
	K8sVersion       string                   `json:"k8sVersion" example:"1.24"`                                      // k8s版本
	IluvatarGpuModel string                   `json:"iluvatarGpuModel" example:"BI-V100"`                             // iluvatar 显卡型号
	IluvatarLicense  string                   `json:"iluvatarLicense" example:"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"`     // iluvatar 显卡授权码
}

type ClusterCreateNodeParam struct {
	Name               string            `json:"name" example:"k8s-node1" validate:"required"`                        // node name for k8s
	Roles              []ClusterNodeRole `json:"roles" example:"Master,Worker" validate:"required"`                   // node role
	GPUProduct         GPUProduct        `json:"gpuProduct" example:"Ascend"`                                         //GPU产品类型 Ascend,Nvidia
	ManagementIPv4Addr string            `json:"managementIPv4Addr" example:"172.31.13.100" validate:"required,ipv4"` // 管理网络v4 ip
	BusinessIPv4Addr   string            `json:"businessIPv4Addr" example:"172.32.4.100" validate:"required,ipv4"`    // 业务网络v4 ip
}

// NodeInfo 节点信息
type NodeInfo struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // master, worker
}

// ClusterConfig 集群配置
type ClusterConfig struct {
	PodCIDR     string                 `json:"podCIDR,omitempty"`
	ServiceCIDR string                 `json:"serviceCIDR,omitempty"`
	DNSDomain   string                 `json:"dnsDomain,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}
