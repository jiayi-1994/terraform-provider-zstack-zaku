// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"zstack.io/edge-go-sdk/pkg/client"
	"zstack.io/edge-go-sdk/pkg/param"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClustersDataSource{}

func NewClustersDataSource() datasource.DataSource {
	return &ClustersDataSource{}
}

// ClustersDataSource defines the data source implementation.
type ClustersDataSource struct {
	client *client.ZeClient
}

// ClustersDataSourceModel describes the data source data model.
type ClustersDataSourceModel struct {
	Clusters []ClusterListItemModel `tfsdk:"clusters"`
	Limit    types.Int64            `tfsdk:"limit"`
	Offset   types.Int64            `tfsdk:"offset"`
	Total    types.Int64            `tfsdk:"total"`
}

// ClusterListItemModel describes a cluster item in the list.
type ClusterListItemModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	Version       types.String `tfsdk:"version"`
	NodeCount     types.Int64  `tfsdk:"node_count"`
	CreateTime    types.String `tfsdk:"create_time"`
	PrometheusURL types.String `tfsdk:"prometheus_url"`
	CreateType    types.String `tfsdk:"create_type"`
	CPUUsage      types.String `tfsdk:"cpu_usage"`
	MemoryUsage   types.String `tfsdk:"memory_usage"`
	StorageUsage  types.String `tfsdk:"storage_usage"`
}

func (d *ClustersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

func (d *ClustersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "获取 ZStack Edge 集群列表的数据源。支持分页查询。",

		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				MarkdownDescription: "每页返回的记录数，默认为 20",
				Optional:            true,
			},
			"offset": schema.Int64Attribute{
				MarkdownDescription: "偏移量，用于分页，默认为 0",
				Optional:            true,
			},
			"total": schema.Int64Attribute{
				MarkdownDescription: "集群总数",
				Computed:            true,
			},
			"clusters": schema.ListNestedAttribute{
				MarkdownDescription: "集群列表",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "集群唯一标识符",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "集群名称",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "集群状态",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "集群版本",
							Computed:            true,
						},
						"node_count": schema.Int64Attribute{
							MarkdownDescription: "集群节点数量",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "集群创建时间",
							Computed:            true,
						},
						"prometheus_url": schema.StringAttribute{
							MarkdownDescription: "Prometheus 监控地址",
							Computed:            true,
						},
						"create_type": schema.StringAttribute{
							MarkdownDescription: "集群创建类型",
							Computed:            true,
						},
						"cpu_usage": schema.StringAttribute{
							MarkdownDescription: "CPU 使用情况",
							Computed:            true,
						},
						"memory_usage": schema.StringAttribute{
							MarkdownDescription: "内存使用情况",
							Computed:            true,
						},
						"storage_usage": schema.StringAttribute{
							MarkdownDescription: "存储使用情况",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ClustersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClustersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClustersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default values
	limit := 20
	offset := 0

	if !data.Limit.IsNull() {
		limit = int(data.Limit.ValueInt64())
	}

	if !data.Offset.IsNull() {
		offset = int(data.Offset.ValueInt64())
	}

	tflog.Info(ctx, "Reading clusters list", map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	})

	// Query clusters
	queryParam := param.NewQueryParam()
	queryParam.Limit(limit).Start(offset)

	clusters, total, err := d.client.PageCluster(queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading clusters",
			fmt.Sprintf("Unable to read clusters, got error: %s", err),
		)
		return
	}

	// Map response to data source model
	data.Total = types.Int64Value(int64(total))
	data.Clusters = make([]ClusterListItemModel, len(clusters))

	for i, cluster := range clusters {
		data.Clusters[i] = ClusterListItemModel{
			ID:            types.Int64Value(cluster.ID),
			Name:          types.StringValue(cluster.Name),
			Status:        types.StringValue(cluster.Status),
			Version:       types.StringValue(cluster.Version),
			NodeCount:     types.Int64Value(int64(cluster.NodeCount)),
			CreateTime:    types.StringValue(cluster.CreateTime.String()),
			PrometheusURL: types.StringValue(cluster.PrometheusURL),
			CreateType:    types.StringValue(string(cluster.CreateType)),
			CPUUsage:      types.StringValue(cluster.Cpu),
			MemoryUsage:   types.StringValue(cluster.Memory),
			StorageUsage:  types.StringValue(cluster.Storage),
		}
	}

	tflog.Debug(ctx, "Clusters list retrieved successfully", map[string]interface{}{
		"total": total,
		"count": len(clusters),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
