package view

import "time"

type UserProjectSimpleView struct {
	ID                int64     `json:"ID"`                // 项目唯一标识
	Name              string    `json:"name"`              // 项目名称
	Readonly          bool      `json:"readonly"`          // 用户项目权限是否只读
	CreateTime        time.Time `json:"createTime"`        // 项目创建时间
	OperationFailFlag string    `json:"operationFailFlag"` // 操作失败标记 DELETE_FAILED: 删除失败
}
