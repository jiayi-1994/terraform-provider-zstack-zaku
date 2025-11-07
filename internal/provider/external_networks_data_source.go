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

var _ datasource.DataSource = &ExternalNetworksDataSource{}

func NewExternalNetworksDataSource() datasource.DataSource {
	return &ExternalNetworksDataSource{}
}

type ExternalNetworksDataSource struct {
	client *client.ZeClient
}

type ExternalNetworksDataSourceModel struct {
	ClusterID types.Int64                      `tfsdk:"cluster_id"`
	Name      types.String                     `tfsdk:"name"`
	Networks  []ExternalNetworkDataSourceModel `tfsdk:"networks"`
}

type ExternalNetworkDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	ClusterID   types.Int64  `tfsdk:"cluster_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Gateway     types.String `tfsdk:"gateway"`
	Netmask     types.String `tfsdk:"netmask"`
	Interface   types.String `tfsdk:"interface"`
	Status      types.String `tfsdk:"status"`
	CreateTime  types.String `tfsdk:"create_time"`
	UpdateTime  types.String `tfsdk:"update_time"`
}

func (d *ExternalNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_networks"
}

func (d *ExternalNetworksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "查询 ZStack Edge 外部网络列表数据源。",

		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.Int64Attribute{
				MarkdownDescription: "集群 ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "外部网络名称（可选，用于过滤）",
				Optional:            true,
			},
			"networks": schema.ListNestedAttribute{
				MarkdownDescription: "外部网络列表",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "外部网络 ID",
							Computed:            true,
						},
						"cluster_id": schema.Int64Attribute{
							MarkdownDescription: "集群 ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "外部网络名称",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "外部网络描述",
							Computed:            true,
						},
						"gateway": schema.StringAttribute{
							MarkdownDescription: "网关地址",
							Computed:            true,
						},
						"netmask": schema.StringAttribute{
							MarkdownDescription: "子网掩码",
							Computed:            true,
						},
						"interface": schema.StringAttribute{
							MarkdownDescription: "网卡接口名称",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "外部网络状态",
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

func (d *ExternalNetworksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExternalNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExternalNetworksDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 构建查询参数
	queryParam := param.NewQueryParam()
	if !data.Name.IsNull() {
		queryParam.AddQ("name=" + data.Name.ValueString())
	}

	// 查询外部网络列表
	networks, _, err := d.client.PageExternalNetwork(int(data.ClusterID.ValueInt64()), queryParam)
	if err != nil {
		resp.Diagnostics.AddError("Failed to query external networks", err.Error())
		return
	}

	// 转换数据
	data.Networks = make([]ExternalNetworkDataSourceModel, 0, len(networks))
	for _, network := range networks {
		// 计算状态
		status := "Inactive"
		if network.ExistNetwork {
			status = "Active"
		}

		networkModel := ExternalNetworkDataSourceModel{
			ID:          types.Int64Value(network.ID),
			ClusterID:   types.Int64Value(network.ClusterID),
			Name:        types.StringValue(network.Name),
			Description: types.StringValue(network.Description),
			Gateway:     types.StringValue(network.Gateway),
			Netmask:     types.StringValue(network.Netmask),
			Interface:   types.StringValue(network.Iface),
			Status:      types.StringValue(status),
			CreateTime:  types.StringValue(network.CreateTime.Format("2006-01-02 15:04:05")),
			// SDK 新版本没有 UpdateTime，使用 CreateTime
			UpdateTime: types.StringValue(network.CreateTime.Format("2006-01-02 15:04:05")),
		}
		data.Networks = append(data.Networks, networkModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
