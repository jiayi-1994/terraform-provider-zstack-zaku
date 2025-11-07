// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"zstack.io/edge-go-sdk/pkg/client"
)

// Ensure ZakuProvider satisfies various provider interfaces.
var _ provider.Provider = &ZakuProvider{}
var _ provider.ProviderWithFunctions = &ZakuProvider{}
var _ provider.ProviderWithEphemeralResources = &ZakuProvider{}

// ZakuProvider defines the provider implementation.
type ZakuProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ZakuProviderModel describes the provider data model.
// 定义 Provider 的配置数据模型
type ZakuProviderModel struct {
	Host        types.String `tfsdk:"host"`         // ZStack Edge 主机地址
	AccessKey   types.String `tfsdk:"access_key"`   // 访问密钥
	SecretKey   types.String `tfsdk:"secret_key"`   // 密钥
}

func (p *ZakuProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	// 设置 Provider 的类型名称，这将作为所有资源和数据源名称的前缀
	// 例如：zstack_cluster, zstack_node 等
	resp.TypeName = "zstack"
	resp.Version = p.version
}

func (p *ZakuProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		// Provider 级别的配置项
		MarkdownDescription: "ZStack Edge Terraform Provider，用于管理 ZStack Edge 集群资源",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "ZStack Edge 主机地址",
				Required:            true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "ZStack Edge 访问密钥",
				Required:            true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "ZStack Edge 密钥",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ZakuProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ZakuProviderModel

	// 从 Terraform 配置中读取 Provider 的配置项
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// 验证必需参数
	if data.Host.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Host Configuration",
			"The provider requires a host configuration. "+
				"Set the host value in the provider configuration or use the ZSTACK_HOST environment variable.",
		)
	}

	if data.AccessKey.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Access Key Configuration",
			"The provider requires an access key configuration. "+
				"Set the access_key value in the provider configuration or use the ZSTACK_ACCESS_KEY environment variable.",
		)
	}

	if data.SecretKey.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Secret Key Configuration",
			"The provider requires a secret key configuration. "+
				"Set the secret_key value in the provider configuration or use the ZSTACK_SECRET_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// 创建 ZStack Edge 客户端
	zeConfig := client.DefaultZeConfig(data.Host.ValueString()).
		AccessKey(data.AccessKey.ValueString(), data.SecretKey.ValueString())
	zeClient := client.NewZeClient(zeConfig)

	// 将客户端传递给 Data Sources 和 Resources
	resp.DataSourceData = zeClient
	resp.ResourceData = zeClient
}

func (p *ZakuProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewExternalNetworkResource,  // 外部网络资源
		NewNodeResource,              // 节点资源
	}
}

func (p *ZakuProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *ZakuProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// 注册所有的 Data Sources
	// 每个 Data Source 都需要在这里注册，Terraform 才能识别和使用
	return []func() datasource.DataSource{
		NewClusterDataSource,          // 单个集群数据源
		NewClustersDataSource,         // 集群列表数据源
		NewExternalNetworksDataSource, // 外部网络列表数据源
		NewNodesDataSource,            // 节点列表数据源
	}
}

func (p *ZakuProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZakuProvider{
			version: version,
		}
	}
}
