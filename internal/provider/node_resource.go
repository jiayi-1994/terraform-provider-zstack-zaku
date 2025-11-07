// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

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
var _ resource.Resource = &NodeResource{}
var _ resource.ResourceWithImportState = &NodeResource{}

func NewNodeResource() resource.Resource {
	return &NodeResource{}
}

// NodeResource defines the resource implementation.
type NodeResource struct {
	client *client.ZeClient
}

// NodeResourceModel describes the resource data model.
type NodeResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ClusterID        types.Int64  `tfsdk:"cluster_id"`
	Password         types.String `tfsdk:"password"`
	ContainerRuntime types.String `tfsdk:"container_runtime"`
	DNSServer        types.String `tfsdk:"dns_server"`
	IluvatarLicense  types.String `tfsdk:"iluvatar_license"`
	Nodes            types.List   `tfsdk:"nodes"` // []NodeAddModel
	ImageDataDisk    types.Map    `tfsdk:"image_data_disk"`
}

// NodeAddModel describes the node add data model.
type NodeAddModel struct {
	Name       types.String `tfsdk:"name"`
	IP         types.String `tfsdk:"ip"`
	BusinessIP types.String `tfsdk:"business_ip"`
	IP6        types.String `tfsdk:"ip6"`
	Port       types.Int64  `tfsdk:"port"`
	Roles      types.List   `tfsdk:"roles"` // []string
	GPUProduct types.String `tfsdk:"gpu_product"`
}

func (r *NodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (r *NodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "管理 ZStack Edge 集群节点资源。提供节点的添加和删除功能。",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "节点资源唯一标识符（节点名称列表）",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.Int64Attribute{
				MarkdownDescription: "集群 ID",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "节点 SSH 密码",
				Required:            true,
				Sensitive:           true,
			},
			"container_runtime": schema.StringAttribute{
				MarkdownDescription: "容器运行时（containerd 或 docker）",
				Optional:            true,
			},
			"dns_server": schema.StringAttribute{
				MarkdownDescription: "DNS 服务器地址",
				Optional:            true,
			},
			"iluvatar_license": schema.StringAttribute{
				MarkdownDescription: "天数 GPU License",
				Optional:            true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "要添加的节点列表",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "节点名称",
							Required:            true,
						},
						"ip": schema.StringAttribute{
							MarkdownDescription: "管理网络 IP 地址",
							Required:            true,
						},
						"business_ip": schema.StringAttribute{
							MarkdownDescription: "业务网络 IP 地址",
							Optional:            true,
						},
						"ip6": schema.StringAttribute{
							MarkdownDescription: "IPv6 地址",
							Optional:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "SSH 端口",
							Required:            true,
						},
						"roles": schema.ListAttribute{
							MarkdownDescription: "节点角色列表（Master, Worker, GPU）",
							Required:            true,
							ElementType:         types.StringType,
						},
						"gpu_product": schema.StringAttribute{
							MarkdownDescription: "GPU 产品类型（Ascend, Nvidia, Iluvatar, Hygon, Enflame）",
							Optional:            true,
						},
					},
				},
			},
			"image_data_disk": schema.MapAttribute{
				MarkdownDescription: "镜像数据盘配置（节点名 -> 磁盘列表）",
				Optional:            true,
				ElementType:         types.ListType{ElemType: types.StringType},
			},
		},
	}
}

func (r *NodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 解析节点列表
	var nodes []NodeAddModel
	resp.Diagnostics.Append(data.Nodes.ElementsAs(ctx, &nodes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 构建添加节点参数
	addParam := param.NodeAddParamOpenApi{
		Password: data.Password.ValueString(),
		NodeAddParam: param.NodeAddParam{
			ClusterID:        data.ClusterID.ValueInt64(),
			Nodes:            make([]param.NodeAddObjParam, 0, len(nodes)),
			ContainerRuntime: param.ContainerRuntime(data.ContainerRuntime.ValueString()),
			DNSServer:        data.DNSServer.ValueString(),
			IluvatarLicense:  data.IluvatarLicense.ValueString(),
		},
	}

	// 如果没有指定容器运行时，使用默认值
	if addParam.ContainerRuntime == "" {
		addParam.ContainerRuntime = param.ContainerRunTimeContainerd
	}

	// 转换节点数据
	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		var roles []string
		resp.Diagnostics.Append(node.Roles.ElementsAs(ctx, &roles, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		nodeRoles := make([]param.ClusterNodeRole, 0, len(roles))
		for _, role := range roles {
			nodeRoles = append(nodeRoles, param.ClusterNodeRole(role))
		}

		nodeParam := param.NodeAddObjParam{
			Name:       node.Name.ValueString(),
			IP:         node.IP.ValueString(),
			BusinessIp: node.BusinessIP.ValueString(),
			IP6:        node.IP6.ValueString(),
			Port:       int(node.Port.ValueInt64()),
			Roles:      nodeRoles,
			GPUProduct: param.GPUProduct(node.GPUProduct.ValueString()),
		}
		addParam.Nodes = append(addParam.Nodes, nodeParam)
		nodeNames = append(nodeNames, node.Name.ValueString())
	}

	// 解析镜像数据盘
	if !data.ImageDataDisk.IsNull() {
		imageDataDisk := make(map[string][]string)
		for k, v := range data.ImageDataDisk.Elements() {
			var disks []string
			if listVal, ok := v.(types.List); ok {
				resp.Diagnostics.Append(listVal.ElementsAs(ctx, &disks, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				imageDataDisk[k] = disks
			}
		}
		addParam.ImageDataDisk = imageDataDisk
	}

	tflog.Debug(ctx, "Adding nodes to cluster", map[string]interface{}{
		"cluster_id": addParam.ClusterID,
		"node_count": len(addParam.Nodes),
		"node_names": nodeNames,
	})

	// 调用 SDK 添加节点（异步操作）
	_, err := r.client.AddNode(int(data.ClusterID.ValueInt64()), addParam, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add nodes", err.Error())
		return
	}

	// 设置资源 ID（使用节点名称列表作为标识）
	data.ID = types.StringValue(fmt.Sprintf("%v", nodeNames))

	tflog.Trace(ctx, "Added nodes to cluster", map[string]interface{}{"node_names": nodeNames})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 节点读取操作可以通过查询节点列表验证节点是否存在
	// 这里简化处理，直接返回当前状态
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NodeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 节点资源不支持更新操作
	resp.Diagnostics.AddError(
		"Update not supported",
		"Node resource does not support update operations. Please destroy and recreate the resource to modify nodes.",
	)
}

func (r *NodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NodeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 解析节点列表
	var nodes []NodeAddModel
	resp.Diagnostics.Append(data.Nodes.ElementsAs(ctx, &nodes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 收集节点名称
	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name.ValueString())
	}

	tflog.Debug(ctx, "Deleting nodes from cluster", map[string]interface{}{
		"cluster_id": data.ClusterID.ValueInt64(),
		"node_names": nodeNames,
	})

	// 调用 SDK 删除节点
	err := r.client.DeleteNode(int(data.ClusterID.ValueInt64()), nodeNames)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete nodes", err.Error())
		return
	}

	tflog.Trace(ctx, "Deleted nodes from cluster", map[string]interface{}{"node_names": nodeNames})
}

func (r *NodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
