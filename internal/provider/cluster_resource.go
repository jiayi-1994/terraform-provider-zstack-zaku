// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

// ClusterResource defines the resource implementation.
type ClusterResource struct {
	client *client.ZeClient
}

// ClusterResourceModel describes the resource data model.
type ClusterResourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	EnableHA         types.Bool   `tfsdk:"enable_ha"`
	NetCombined      types.Bool   `tfsdk:"net_combined"`
	Port             types.Int64  `tfsdk:"port"`
	Password         types.String `tfsdk:"password"`
	ManagementVipV4  types.String `tfsdk:"management_vip_v4"`
	BusinessVipV4    types.String `tfsdk:"business_vip_v4"`
	MaxPodPerNode    types.Int64  `tfsdk:"max_pod_per_node"`
	PodCidrV4        types.String `tfsdk:"pod_cidr_v4"`
	ServiceCidrV4    types.String `tfsdk:"service_cidr_v4"`
	DNSServer        types.String `tfsdk:"dns_server"`
	IstioEnabled     types.Bool   `tfsdk:"istio_enabled"`
	K8sVersion       types.String `tfsdk:"k8s_version"`
	IluvatarGpuModel types.String `tfsdk:"iluvatar_gpu_model"`
	IluvatarLicense  types.String `tfsdk:"iluvatar_license"`
	Nodes            types.List   `tfsdk:"nodes"` // []ClusterNodeModel
	DataDisk         types.Map    `tfsdk:"data_disk"`
	ImageDataDisk    types.Map    `tfsdk:"image_data_disk"`

	// Computed fields
	Status        types.String `tfsdk:"status"`
	Version       types.String `tfsdk:"version"`
	NodeCount     types.Int64  `tfsdk:"node_count"`
	CreateTime    types.String `tfsdk:"create_time"`
	PrometheusURL types.String `tfsdk:"prometheus_url"`
}

// ClusterNodeModel describes the cluster node data model.
type ClusterNodeModel struct {
	Name               types.String `tfsdk:"name"`
	Roles              types.List   `tfsdk:"roles"` // []string
	GPUProduct         types.String `tfsdk:"gpu_product"`
	ManagementIPv4Addr types.String `tfsdk:"management_ipv4_addr"`
	BusinessIPv4Addr   types.String `tfsdk:"business_ipv4_addr"`
}

func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "管理 ZStack Edge 集群资源。提供集群的创建、读取、更新和删除功能。",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "集群唯一标识符",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "集群名称",
				Required:            true,
			},
			"enable_ha": schema.BoolAttribute{
				MarkdownDescription: "是否启用高可用",
				Optional:            true,
			},
			"net_combined": schema.BoolAttribute{
				MarkdownDescription: "管理网络和业务网络是否复用",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "SSH 端口",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "SSH 密码（加密后的）",
				Required:            true,
				Sensitive:           true,
			},
			"management_vip_v4": schema.StringAttribute{
				MarkdownDescription: "管理网络 VIP IPv4 地址",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"business_vip_v4": schema.StringAttribute{
				MarkdownDescription: "业务网络 VIP IPv4 地址",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_pod_per_node": schema.Int64Attribute{
				MarkdownDescription: "每个节点的最大 Pod 数量",
				Optional:            true,
			},
			"pod_cidr_v4": schema.StringAttribute{
				MarkdownDescription: "Kubernetes Pod CIDR IPv4",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_cidr_v4": schema.StringAttribute{
				MarkdownDescription: "Kubernetes Service CIDR IPv4",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dns_server": schema.StringAttribute{
				MarkdownDescription: "DNS 服务器 IP 地址",
				Required:            true,
			},
			"istio_enabled": schema.BoolAttribute{
				MarkdownDescription: "是否启用 Istio",
				Optional:            true,
			},
			"k8s_version": schema.StringAttribute{
				MarkdownDescription: "Kubernetes 版本",
				Optional:            true,
			},
			"iluvatar_gpu_model": schema.StringAttribute{
				MarkdownDescription: "天数 GPU 型号",
				Optional:            true,
			},
			"iluvatar_license": schema.StringAttribute{
				MarkdownDescription: "天数 GPU 许可证",
				Optional:            true,
				Sensitive:           true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "集群节点列表",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "节点名称",
							Required:            true,
						},
						"roles": schema.ListAttribute{
							MarkdownDescription: "节点角色列表 (例如: Master, Worker)",
							Required:            true,
							ElementType:         types.StringType,
						},
						"gpu_product": schema.StringAttribute{
							MarkdownDescription: "GPU 产品类型 (Ascend, Nvidia)",
							Optional:            true,
						},
						"management_ipv4_addr": schema.StringAttribute{
							MarkdownDescription: "管理网络 IPv4 地址",
							Required:            true,
						},
						"business_ipv4_addr": schema.StringAttribute{
							MarkdownDescription: "业务网络 IPv4 地址",
							Required:            true,
						},
					},
				},
			},
			"data_disk": schema.MapAttribute{
				MarkdownDescription: "数据磁盘配置",
				Required:            true,
				ElementType:         types.ListType{ElemType: types.StringType},
			},
			"image_data_disk": schema.MapAttribute{
				MarkdownDescription: "镜像数据磁盘配置",
				Optional:            true,
				ElementType:         types.ListType{ElemType: types.StringType},
			},

			// Computed attributes
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
				MarkdownDescription: "Prometheus 地址",
				Computed:            true,
			},
		},
	}
}

func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating cluster", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Build cluster create parameters
	createParam := param.ClusterCreateParam{
		Name:            data.Name.ValueString(),
		EnableHA:        !data.EnableHA.IsNull() && data.EnableHA.ValueBool(),
		NetCombined:     !data.NetCombined.IsNull() && data.NetCombined.ValueBool(),
		Port:            int(data.Port.ValueInt64()),
		Password:        data.Password.ValueString(),
		ManagementVipV4: data.ManagementVipV4.ValueString(),
		BusinessVipV4:   data.BusinessVipV4.ValueString(),
		PodCidrV4:       data.PodCidrV4.ValueString(),
		ServiceCidrV4:   data.ServiceCidrV4.ValueString(),
		DNSServer:       data.DNSServer.ValueString(),
		IstioEnabled:    !data.IstioEnabled.IsNull() && data.IstioEnabled.ValueBool(),
	}

	if !data.MaxPodPerNode.IsNull() {
		createParam.MaxPodPerNode = int(data.MaxPodPerNode.ValueInt64())
	}

	if !data.K8sVersion.IsNull() {
		createParam.K8sVersion = data.K8sVersion.ValueString()
	}

	if !data.IluvatarGpuModel.IsNull() {
		createParam.IluvatarGpuModel = data.IluvatarGpuModel.ValueString()
	}

	if !data.IluvatarLicense.IsNull() {
		createParam.IluvatarLicense = data.IluvatarLicense.ValueString()
	}

	// Parse nodes
	var nodes []ClusterNodeModel
	resp.Diagnostics.Append(data.Nodes.ElementsAs(ctx, &nodes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam.Nodes = make([]param.ClusterCreateNodeParam, len(nodes))
	for i, node := range nodes {
		var roles []string
		resp.Diagnostics.Append(node.Roles.ElementsAs(ctx, &roles, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		nodeParam := param.ClusterCreateNodeParam{
			Name:               node.Name.ValueString(),
			ManagementIPv4Addr: node.ManagementIPv4Addr.ValueString(),
			BusinessIPv4Addr:   node.BusinessIPv4Addr.ValueString(),
		}

		// Convert roles to ClusterNodeRole type
		nodeParam.Roles = make([]param.ClusterNodeRole, len(roles))
		for j, role := range roles {
			nodeParam.Roles[j] = param.ClusterNodeRole(role)
		}

		if !node.GPUProduct.IsNull() {
			nodeParam.GPUProduct = param.GPUProduct(node.GPUProduct.ValueString())
		}

		createParam.Nodes[i] = nodeParam
	}

	// Parse data_disk
	var dataDisk map[string][]string
	resp.Diagnostics.Append(data.DataDisk.ElementsAs(ctx, &dataDisk, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createParam.DataDisk = dataDisk

	// Parse image_data_disk
	if !data.ImageDataDisk.IsNull() {
		var imageDataDisk map[string][]string
		resp.Diagnostics.Append(data.ImageDataDisk.ElementsAs(ctx, &imageDataDisk, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createParam.ImageDataDisk = imageDataDisk
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ(fmt.Sprintf("name=%s", data.Name.ValueString()))

	clusters, _, err := r.client.PageCluster(queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying created cluster",
			fmt.Sprintf("Unable to query cluster by name '%s', got error: %s", data.Name.ValueString(), err),
		)
		return
	}
	taskID := ""
	if len(clusters) > 0 {
		if clusters[0].Status == "Status_Cluster_Create_Failed" {
			taskID, err = r.client.RecreateCluster(int(clusters[0].ID), false)
			if err != nil {
				resp.Diagnostics.AddError("Error recreate cluster",
					fmt.Sprintf("Unable to create cluster, got error: %s", err))
				return
			}
		}
	} else {
		// Create cluster
		taskID, err = r.client.CreateCluster(createParam, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating cluster",
				fmt.Sprintf("Unable to create cluster, got error: %s", err),
			)
			return
		}
	}

	tflog.Info(ctx, "Cluster creation initiated", map[string]interface{}{
		"task_id":      taskID,
		"cluster_name": data.Name.ValueString(),
	})

	// Query cluster by name to get the cluster ID
	// taskID is just a task UUID, not the cluster ID
	queryParam = param.NewQueryParam()
	queryParam.AddQ(fmt.Sprintf("name=%s", data.Name.ValueString()))

	clusters, _, err = r.client.PageCluster(queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying created cluster",
			fmt.Sprintf("Unable to query cluster by name '%s', got error: %s", data.Name.ValueString(), err),
		)
		return
	}

	if len(clusters) == 0 {
		resp.Diagnostics.AddError(
			"Cluster not found",
			fmt.Sprintf("Created cluster '%s' not found in query results", data.Name.ValueString()),
		)
		return
	}

	// Use the first matching cluster
	data.ID = types.Int64Value(int64(clusters[0].ID))

	// Read cluster details
	r.readCluster(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Cluster created successfully", map[string]interface{}{
		"id": data.ID.ValueInt64(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readCluster(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "Cluster update is not fully supported, most fields require replacement")

	// For now, just read the cluster to sync state
	r.readCluster(ctx, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := int(data.ID.ValueInt64())

	tflog.Info(ctx, "Deleting cluster", map[string]interface{}{
		"id": clusterID,
	})

	_, err := r.client.DeleteCluster(clusterID, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			fmt.Sprintf("Unable to delete cluster %d, got error: %s", clusterID, err),
		)
		return
	}

	tflog.Info(ctx, "Cluster deleted successfully", map[string]interface{}{
		"id": clusterID,
	})
}

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	clusterID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing cluster ID",
			fmt.Sprintf("Unable to parse cluster ID '%s': %s", req.ID, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), clusterID)...)
}

// Helper function to read cluster details
func (r *ClusterResource) readCluster(ctx context.Context, data *ClusterResourceModel, diags *diag.Diagnostics) {
	clusterID := int(data.ID.ValueInt64())

	clusterDetails, err := r.client.GetClusterDetails(clusterID)
	if err != nil {
		diags.AddError(
			"Error reading cluster",
			fmt.Sprintf("Unable to read cluster %d, got error: %s", clusterID, err),
		)
		return
	}

	// Update computed fields
	data.Status = types.StringValue(clusterDetails.Status)
	data.Version = types.StringValue(clusterDetails.Version)
	data.NodeCount = types.Int64Value(int64(clusterDetails.NodeCount))
	data.CreateTime = types.StringValue(clusterDetails.CreateTime.String())
	data.PrometheusURL = types.StringValue(clusterDetails.PrometheusURL)

	tflog.Debug(ctx, "Cluster details retrieved", map[string]interface{}{
		"id":     clusterID,
		"status": clusterDetails.Status,
	})
}
