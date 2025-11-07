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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClusterDataSource{}

func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterDataSource defines the data source implementation.
type ClusterDataSource struct {
	client *client.ZeClient
}

// ClusterDataSourceModel describes the data source data model.
type ClusterDataSourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	Version       types.String `tfsdk:"version"`
	NodeCount     types.Int64  `tfsdk:"node_count"`
	CreateTime    types.String `tfsdk:"create_time"`
	PrometheusURL types.String `tfsdk:"prometheus_url"`
	CreateType    types.String `tfsdk:"create_type"`
	Description   types.String `tfsdk:"description"`
	CPUUsage      types.String `tfsdk:"cpu_usage"`
	MemoryUsage   types.String `tfsdk:"memory_usage"`
	StorageUsage  types.String `tfsdk:"storage_usage"`
}

func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "获取 ZStack Edge 集群信息的数据源。可以通过集群 ID 查询集群的详细信息。",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "集群唯一标识符",
				Required:            true,
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
			"description": schema.StringAttribute{
				MarkdownDescription: "集群描述",
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
	}
}

func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := int(data.ID.ValueInt64())

	tflog.Info(ctx, "Reading cluster details", map[string]interface{}{
		"id": clusterID,
	})

	// Get cluster details
	clusterDetails, err := d.client.GetClusterDetails(clusterID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			fmt.Sprintf("Unable to read cluster %d, got error: %s", clusterID, err),
		)
		return
	}

	// Map response to data source model
	data.Name = types.StringValue(clusterDetails.Name)
	data.Status = types.StringValue(clusterDetails.Status)
	data.Version = types.StringValue(clusterDetails.Version)
	data.NodeCount = types.Int64Value(int64(clusterDetails.NodeCount))
	data.CreateTime = types.StringValue(clusterDetails.CreateTime.String())
	data.PrometheusURL = types.StringValue(clusterDetails.PrometheusURL)
	data.CreateType = types.StringValue(string(clusterDetails.CreateType))
	data.Description = types.StringValue(clusterDetails.Description)
	data.CPUUsage = types.StringValue(clusterDetails.Cpu)
	data.MemoryUsage = types.StringValue(clusterDetails.Memory)
	data.StorageUsage = types.StringValue(clusterDetails.Storage)

	tflog.Debug(ctx, "Cluster details retrieved successfully", map[string]interface{}{
		"id":     clusterID,
		"name":   clusterDetails.Name,
		"status": clusterDetails.Status,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
