// Copyright (c) LeoColomb.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/LeoColomb/terraform-provider-anytype/internal/client"
)

// Ensure AnytypeProvider satisfies the provider interfaces.
var _ provider.Provider = &AnytypeProvider{}

// AnytypeProvider is the Terraform provider for Anytype.
type AnytypeProvider struct {
	// version is the provider version ("dev" locally, set by goreleaser on release).
	version string
}

// AnytypeProviderModel is the user-facing provider configuration.
type AnytypeProviderModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	APIKey     types.String `tfsdk:"api_key"`
	APIVersion types.String `tfsdk:"api_version"`
}

func (p *AnytypeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "anytype"
	resp.Version = p.version
}

func (p *AnytypeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Anytype provider manages [Anytype](https://anytype.io) resources " +
			"through the official Anytype HTTP API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Anytype API endpoint. Defaults to the local desktop API " +
					"(`http://127.0.0.1:31009`). Can also be set with the `ANYTYPE_ENDPOINT` environment variable.",
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Anytype API key obtained via the `/v1/auth/api_keys` flow. " +
					"Can also be set with the `ANYTYPE_API_KEY` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"api_version": schema.StringAttribute{
				MarkdownDescription: "Value sent in the `Anytype-Version` header. Defaults to the " +
					"version the provider was generated against (`" + client.APIVersion + "`). " +
					"Can also be set with the `ANYTYPE_API_VERSION` environment variable.",
				Optional: true,
			},
		},
	}
}

func (p *AnytypeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AnytypeProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := stringValueOrEnv(data.Endpoint, "ANYTYPE_ENDPOINT", client.DefaultEndpoint)
	apiKey := stringValueOrEnv(data.APIKey, "ANYTYPE_API_KEY", "")
	apiVersion := stringValueOrEnv(data.APIVersion, "ANYTYPE_API_VERSION", client.APIVersion)

	if apiKey == "" {
		resp.Diagnostics.AddAttributeWarning(
			frameworkPath("api_key"),
			"Missing Anytype API key",
			"No api_key was configured and ANYTYPE_API_KEY is not set. "+
				"Anytype API calls will likely return 401 Unauthorized.",
		)
	}

	c, err := client.New(client.Config{
		Endpoint:   endpoint,
		APIKey:     apiKey,
		APIVersion: apiVersion,
		UserAgent:  "terraform-provider-anytype/" + p.version,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to configure Anytype client", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *AnytypeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSpaceResource,
		NewTypeResource,
		NewPropertyResource,
		NewTagResource,
		NewObjectResource,
	}
}

func (p *AnytypeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSpaceDataSource,
		NewSpacesDataSource,
		NewTypeDataSource,
		NewTypesDataSource,
		NewPropertyDataSource,
		NewPropertiesDataSource,
		NewTagDataSource,
		NewTagsDataSource,
		NewObjectDataSource,
		NewObjectsDataSource,
		NewMemberDataSource,
		NewMembersDataSource,
		NewTemplateDataSource,
		NewTemplatesDataSource,
	}
}

// New returns a constructor usable with providerserver.Serve.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AnytypeProvider{version: version}
	}
}

func stringValueOrEnv(v types.String, env, fallback string) string {
	if !v.IsNull() && !v.IsUnknown() && v.ValueString() != "" {
		return v.ValueString()
	}
	if fromEnv := os.Getenv(env); fromEnv != "" {
		return fromEnv
	}
	return fallback
}
