// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"zstack.io/edge-go-sdk/pkg/client"
	"zstack.io/edge-go-sdk/pkg/param"
)

var _ datasource.DataSource = &NodesDataSource{}

func NewNodesDataSource() datasource.DataSource {
	return &NodesDataSource{}
}

type NodesDataSource struct {
	client *client.ZeClient
}

type NodesDataSourceModel struct {
	ClusterID types.Int64           `tfsdk:"cluster_id"`
	Name      types.String          `tfsdk:"name"`
	Nodes     []NodeDataSourceModel `tfsdk:"nodes"`
}

type NodeDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	ClusterID  types.Int64  `tfsdk:"cluster_id"`
	IP         types.String `tfsdk:"ip"`
	Role       types.String `tfsdk:"role"`
	Status     types.String `tfsdk:"status"`
	CPU        types.String `tfsdk:"cpu"`
	Memory     types.String `tfsdk:"memory"`
	Storage    types.String `tfsdk:"storage"`
	CreateTime types.String `tfsdk:"create_time"`
	UpdateTime types.String `tfsdk:"update_time"`
}

func (d *NodesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nodes"
}

func (d *NodesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "查询 ZStack Edge 集群节点列表数据源。",

		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.Int64Attribute{
				MarkdownDescription: "集群 ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "节点名称（可选，用于过滤）",
				Optional:            true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "节点列表",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "节点 ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "节点名称",
							Computed:            true,
						},
						"cluster_id": schema.Int64Attribute{
							MarkdownDescription: "集群 ID",
							Computed:            true,
						},
						"ip": schema.StringAttribute{
							MarkdownDescription: "节点 IP 地址",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "节点角色",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "节点状态",
							Computed:            true,
						},
						"cpu": schema.StringAttribute{
							MarkdownDescription: "CPU 信息",
							Computed:            true,
						},
						"memory": schema.StringAttribute{
							MarkdownDescription: "内存信息",
							Computed:            true,
						},
						"storage": schema.StringAttribute{
							MarkdownDescription: "存储信息",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "创建时间",
							Computed:            true,
						},
						"update_time": schema.StringAttribute{
							MarkdownDescription: "更新时间",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NodesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZeClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZeClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NodesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 构建查询参数
	queryParam := param.NewQueryParam()
	if !data.Name.IsNull() {
		queryParam.AddQ("name=" + data.Name.ValueString())
	}

	// 查询节点列表
	nodes, _, err := d.client.PageNode(int(data.ClusterID.ValueInt64()), queryParam)
	if err != nil {
		resp.Diagnostics.AddError("Failed to query nodes", err.Error())
		return
	}

	// 转换数据
	data.Nodes = make([]NodeDataSourceModel, 0, len(nodes))
	for _, node := range nodes {
		nodeModel := NodeDataSourceModel{
			ID:         types.Int64Value(node.ID),
			Name:       types.StringValue(node.Name),
			ClusterID:  types.Int64Value(node.ClusterID),
			IP:         types.StringValue(node.IP),
			Role:       types.StringValue(node.Role),
			Status:     types.StringValue(node.Status),
			CPU:        types.StringValue(node.CPU),
			Memory:     types.StringValue(node.Memory),
			Storage:    types.StringValue(node.Storage),
			CreateTime: types.StringValue(node.CreateTime.Format("2006-01-02 15:04:05")),
			UpdateTime: types.StringValue(node.UpdateTime.Format("2006-01-02 15:04:05")),
		}
		data.Nodes = append(data.Nodes, nodeModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
