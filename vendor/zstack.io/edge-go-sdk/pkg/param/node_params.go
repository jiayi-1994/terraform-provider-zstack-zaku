package param

type ContainerRuntime string

const (
	ContainerRunTimeDocker     ContainerRuntime = "docker"
	ContainerRunTimeContainerd ContainerRuntime = "containerd"
)

type ClusterNodeRole string

const (
	NodeRoleMaster ClusterNodeRole = "Master"
	NodeRoleWorker ClusterNodeRole = "Worker"
	NodeRoleGPU    ClusterNodeRole = "GPU"
)

type GPUProduct string

const (
	GPUProductAscend   GPUProduct = "Ascend"
	GPUProductNvidia   GPUProduct = "Nvidia"
	GPUProductIluvatar GPUProduct = "Iluvatar"
	GPUProductHygon    GPUProduct = "Hygon"
	GPUProductEnflame  GPUProduct = "Enflame"
)

// NodeAddParamOpenApi 添加节点参数（OpenAPI）
type NodeAddParamOpenApi struct {
	NodeAddParam
	Password string `json:"password"` // node password
}

type NodeAddParam struct {
	ClusterID        int64               `json:"clusterID" example:"1"` //集群ID
	Nodes            []NodeAddObjParam   `json:"nodes"`
	ContainerRuntime ContainerRuntime    `json:"containerRuntime" example:"containerd"` //容器运行时间 containerd / docker
	DNSServer        string              `json:"dnsServer" example:"223.5.5.5"`         // dns server
	ImageDataDisk    map[string][]string `json:"imageDataDisk" yaml:"imageDataDisk"`    //镜像数据盘
	IluvatarLicense  string              `json:"iluvatarLicense" example:"xxxx"`        //iluvatar license
}

type NodeAddObjParam struct {
	Name       string            `json:"name" example:"k8s-node1"`     //node name for k8s
	IP         string            `json:"ip" example:"192.168.30.100"`  //manager accessible ip
	BusinessIp string            `json:"businessIp" example:""`        //business ip
	IP6        string            `json:"ip6" example:"fd00::2000:303"` //accessible ipv6
	Port       int               `json:"port" example:"22"`            // ssh port
	Roles      []ClusterNodeRole `json:"roles" example:"GPU"`          //节点角色  GPU,Master,Worker
	GPUProduct GPUProduct        `json:"gpuProduct" example:"Ascend"`  //GPU产品类型 Ascend,Nvidia
}
