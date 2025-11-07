package param

// CloudProjectCreateParam 创建项目参数
type CloudProjectCreateParam struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CloudProjectUpdateParam 更新项目参数
type CloudProjectUpdateParam struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CloudProjectUserAddParam 项目添加用户参数
type CloudProjectUserAddParam struct {
	ProjectUuids []string `json:"projectUuids"` //项目UUID
	Usernames    []string `json:"usernames"`    //登录用户名
}

// ProjectQuotaParam 项目配额参数
type ProjectQuotaParam struct {
	Quotas []ResourceQuota `json:"quotas"`
}

// ResourceQuota 资源配额
type ResourceQuota struct {
	ResourceType string  `json:"resourceType"`
	Quota        float64 `json:"quota"`
}

// CloudUserCreateParam 创建用户参数
type CloudUserCreateParam struct {
	Name        string `json:"name"`
	ThirdParty  bool   `json:"thirdParty"` //是否第三方
	Description string `json:"description,omitempty"`
}
