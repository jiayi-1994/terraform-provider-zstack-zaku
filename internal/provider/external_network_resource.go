// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"zstack.io/edge-go-sdk/pkg/client"
	"zstack.io/edge-go-sdk/pkg/param"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExternalNetworkResource{}
var _ resource.ResourceWithImportState = &ExternalNetworkResource{}

func NewExternalNetworkResource() resource.Resource {
	return &ExternalNetworkResource{}
}

// ExternalNetworkResource defines the resource implementation.
type ExternalNetworkResource struct {
	client *client.ZeClient
}

// ExternalNetworkResourceModel describes the resource data model.
type ExternalNetworkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ClusterID   types.Int64  `tfsdk:"cluster_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Gateway     types.String `tfsdk:"gateway"`
	Netmask     types.String `tfsdk:"netmask"`
	Interface   types.String `tfsdk:"interface"`

	// Computed fields
	Status     types.String `tfsdk:"status"`
	CreateTime types.String `tfsdk:"create_time"`
	UpdateTime types.String `tfsdk:"update_time"`
}

func (r *ExternalNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_network"
}

func (r *ExternalNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "管理 ZStack Edge 外部网络资源。提供外部网络的创建、读取和删除功能。",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "外部网络唯一标识符",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.Int64Attribute{
				MarkdownDescription: "集群 ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "外部网络名称",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "外部网络描述",
				Optional:            true,
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "网关地址",
				Required:            true,
			},
			"netmask": schema.StringAttribute{
				MarkdownDescription: "子网掩码",
				Required:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "网卡接口名称",
				Required:            true,
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
	}
}

func (r *ExternalNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZeClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZeClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ExternalNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExternalNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 创建外部网络参数
	createParam := param.ExternalNetworkCreateParam{
		ClusterID:   int(data.ClusterID.ValueInt64()),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Gateway:     data.Gateway.ValueString(),
		Netmask:     data.Netmask.ValueString(),
		Iface:       data.Interface.ValueString(),
	}

	tflog.Debug(ctx, "Creating external network", map[string]interface{}{
		"cluster_id": createParam.ClusterID,
		"name":       createParam.Name,
		"interface":  createParam.Iface,
	})

	// 调用 SDK 创建外部网络
	actionID, err := r.client.CreateExternalNetwork(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create external network",
			fmt.Sprintf("API returned error: %s\n\nPlease check:\n"+
				"1. Cluster ID %d exists and is accessible\n"+
				"2. Network interface '%s' exists on cluster nodes\n"+
				"3. Gateway %s and netmask %s are valid",
				err.Error(), createParam.ClusterID, createParam.Iface,
				createParam.Gateway, createParam.Netmask),
		)
		return
	}

	if actionID == "" {
		resp.Diagnostics.AddError(
			"Failed to create external network",
			"API returned empty network ID. The creation may have failed on the server side.",
		)
		return
	}

	tflog.Info(ctx, "External network created successfully", map[string]interface{}{
		"action_id": actionID,
		"name":      createParam.Name,
	})

	// 读取创建后的外部网络详情，获取 ID、create_time 等 Computed 字段
	if err := r.readExternalNetwork(ctx, &data); err != nil {
		resp.Diagnostics.AddError("Failed to read external network after creation", err.Error())
		return
	}

	tflog.Trace(ctx, "Created external network resource", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExternalNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ExternalNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Reading external network", map[string]interface{}{
		"cluster_id": data.ClusterID.ValueInt64(),
		"name":       data.Name.ValueString(),
	})

	if err := r.readExternalNetwork(ctx, &data); err != nil {
		// 如果资源不存在，从状态中移除（Terraform 会在下次 apply 时重新创建）
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "External network not found in backend, removing from state", map[string]interface{}{
				"id":   data.ID.ValueString(),
				"name": data.Name.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		
		resp.Diagnostics.AddError("Failed to read external network", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExternalNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ExternalNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 外部网络目前不支持更新操作
	resp.Diagnostics.AddError(
		"Update not supported",
		"External network does not support update operations. Please destroy and recreate the resource.",
	)
}

func (r *ExternalNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ExternalNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting external network", map[string]interface{}{"id": data.ID.ValueString()})

	// 注意：SDK 中没有提供删除外部网络的 API，这里只记录日志
	tflog.Warn(ctx, "External network deletion is not implemented in SDK")
	resp.Diagnostics.AddWarning(
		"Delete not implemented",
		"The SDK does not provide an API to delete external networks. The resource will be removed from state only.",
	)
}

func (r *ExternalNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ErrResourceNotFound 资源未找到错误
var ErrResourceNotFound = errors.New("resource not found")

// readExternalNetwork 读取外部网络详情
func (r *ExternalNetworkResource) readExternalNetwork(ctx context.Context, data *ExternalNetworkResourceModel) error {
	clusterId := int(data.ClusterID.ValueInt64())

	// 查询外部网络列表
	queryParam := param.NewQueryParam()
	queryParam.AddQ("name=" + data.Name.ValueString())

	networks, _, err := r.client.PageExternalNetwork(clusterId, queryParam)
	if err != nil {
		return fmt.Errorf("failed to query external network: %w", err)
	}

	if len(networks) == 0 {
		return ErrResourceNotFound
	}

	// 使用第一个匹配的网络
	network := networks[0]

	data.ID = types.StringValue(strconv.FormatInt(network.ID, 10))
	data.ClusterID = types.Int64Value(network.ClusterID)
	data.Name = types.StringValue(network.Name)
	data.Description = types.StringValue(network.Description)
	data.Gateway = types.StringValue(network.Gateway)
	data.Netmask = types.StringValue(network.Netmask)
	data.Interface = types.StringValue(network.Iface)

	// 计算状态：根据网络是否存在来判断
	if network.ExistNetwork {
		data.Status = types.StringValue("Active")
	} else {
		data.Status = types.StringValue("Inactive")
	}

	data.CreateTime = types.StringValue(network.CreateTime.Format("2006-01-02 15:04:05"))
	// SDK 新版本没有 UpdateTime，使用 CreateTime
	data.UpdateTime = types.StringValue(network.CreateTime.Format("2006-01-02 15:04:05"))

	return nil
}
