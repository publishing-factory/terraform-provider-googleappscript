package googleappscript

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"google.golang.org/api/script/v1"
	"log"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// orderResource is the resource implementation.
type projectResource struct {
	service *script.Service
}

// Metadata returns the resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"parent_id": schema.StringAttribute{
				Description: "Google Drive folder id",
				Optional:    true,
			},
			"title": schema.StringAttribute{
				Description: "The title of your Google Apps Script",
				Required:    true,
			},
			"script_id": schema.StringAttribute{
				Description: "The id of your Google Apps Script",
				Computed:    true,
			},
			"create_time": schema.StringAttribute{
				Description: "The creation timestamp of your Google Apps Script",
				Computed:    true,
			},
			"update_time": schema.StringAttribute{
				Description: "The last updated timestamp of your Google Apps Script",
				Computed:    true,
			},
			"deployment_id": schema.StringAttribute{
				Description: "The deployment id associated to your Google Apps Script",
				Computed:    true,
			},
		},
	}
}

// Create a new resource
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config googleAppScriptProviderModel
	req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new project
	request := script.CreateProjectRequest{
		Title:    plan.Title.ValueString(),
		ParentId: plan.ParentId.ValueString(),
	}

	project, err := r.service.Projects.Create(&request).Do()
	if err != nil {
		// The API encountered a problem.
		log.Fatalf("The API returned an error: %v", err)
	}

	content := &script.Content{
		ScriptId: project.ScriptId,
		Files: []*script.File{{
			Name:   "hello",
			Type:   "SERVER_JS",
			Source: "function helloWorld() {\n  console.log('Hello, world!');}",
		}, {
			Name: "appsscript",
			Type: "JSON",
			Source: "{\"timeZone\":\"America/New_York\",\"exceptionLogging\":" +
				"\"CLOUD\"}",
		}},
	}

	updatedContent, err := r.service.Projects.UpdateContent(project.ScriptId,
		content).Do()
	if err != nil {
		// The API encountered a problem.
		log.Fatalf("The API returned an error: %v", err)
	}

	log.Print(updatedContent)

	scriptConfig := script.Version{
		VersionNumber: 1,
		Description:   project.Title,
	}

	version, err := r.service.Projects.Versions.Create(project.ScriptId, &scriptConfig).Do()
	log.Print(version)

	deploymentConfig := script.DeploymentConfig{
		ScriptId:         project.ScriptId,
		Description:      project.Title,
		ManifestFileName: "appsscript",
		VersionNumber:    1,
	}

	deployment, err := r.service.Projects.Deployments.Create(project.ScriptId, &deploymentConfig).Do()
	if err != nil {
		// The API encountered a problem.
		log.Fatalf("The API returned an error: %v", err)
	}

	// Map response body to schema and populate Computed attribute values
	plan.ScriptId = types.StringValue(project.ScriptId)
	plan.CreateTime = types.StringValue(project.CreateTime)
	plan.UpdateTime = types.StringValue(project.UpdateTime)
	plan.DeploymentId = types.StringValue(deployment.DeploymentId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	project, err := r.service.Projects.Get(state.ScriptId.ValueString()).Do()
	if err != nil {
		return
	}
	// Overwrite items with refreshed state
	state.CreateTime = types.StringValue(project.CreateTime)
	state.ScriptId = types.StringValue(project.ScriptId)
	state.UpdateTime = types.StringValue(project.UpdateTime)
	state.Title = types.StringValue(project.Title)
	state.ParentId = types.StringValue(project.ParentId)

	deployments, err := r.service.Projects.Deployments.List(project.ScriptId).Do()
	if err != nil {
		// The API encountered a problem.
		log.Fatalf("The API returned an error: %v", err)
	}

	state.DeploymentId = types.StringValue(deployments.Deployments[0].DeploymentId)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No API available to delete project
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No API available to update project
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("script_id"), req, resp)
}

type projectResourceModel struct {
	Title        basetypes.StringValue `tfsdk:"title"`
	ParentId     basetypes.StringValue `tfsdk:"parent_id"`
	ScriptId     basetypes.StringValue `tfsdk:"script_id"`
	CreateTime   basetypes.StringValue `tfsdk:"create_time"`
	UpdateTime   basetypes.StringValue `tfsdk:"update_time"`
	DeploymentId basetypes.StringValue `tfsdk:"deployment_id"`
}

// Configure adds the provider configured client to the resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.service = req.ProviderData.(*script.Service)
}
