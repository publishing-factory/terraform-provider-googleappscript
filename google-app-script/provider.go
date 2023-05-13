package googleappscript

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/script/v1"
	"net/http"
	"strings"
)

var (
	_ provider.Provider = &googleAppScriptProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &googleAppScriptProvider{}
}

type googleAppScriptProvider struct{}

// Metadata returns the provider type name.
func (p *googleAppScriptProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googleappscript"
}

func (p *googleAppScriptProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "generated token.json using quickstart folder",
				Required:    true,
				Sensitive:   true,
			},
			"credentials": schema.StringAttribute{
				Description: "Google oauth2 credentials json file",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *googleAppScriptProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config googleAppScriptProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Google API token",
			"The provider cannot create the Google Script API client as there is an unknown configuration value for token. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := getClient(config.Token.ValueString(), config.Credentials.ValueString())
	if client == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"failed to load google api config",
			"failed to load google api config",
		)
	}

	service, err := script.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = service
	resp.ResourceData = service

}

// DataSources defines the data sources implemented in the provider.
func (p *googleAppScriptProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *googleAppScriptProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
	}
}

// googleAppScriptProviderModel maps provider schema data to a Go type.
type googleAppScriptProviderModel struct {
	Token       types.String `tfsdk:"token"`
	Credentials types.String `tfsdk:"credentials"`
}

// func getClient(config *oauth2.Config) *http.Client {
// func getClient(token *oauth2.Token) *http.Client {
func getClient(token string, credentials string) *http.Client {
	tok := &oauth2.Token{}
	errorToken := json.NewDecoder(strings.NewReader(token)).Decode(tok)

	config, errorCreds := google.ConfigFromJSON([]byte(credentials), "https://www.googleapis.com/auth/script.projects")

	if errorToken != nil {
		return nil
	}

	if errorCreds != nil {
		return nil
	}

	return config.Client(context.Background(), tok)
}
