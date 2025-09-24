package devs

import (
	"context"

	"terraform-provider-devops/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &devDataSource{}
	_ datasource.DataSourceWithConfigure = &devDataSource{}
)

// NewDevDataSource is a helper function to simplify the provider implementation.
func NewDevDataSource() datasource.DataSource {
	return &devDataSource{}
}

// devDataSource is the data source implementation.
type devDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *devDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

// Schema defines the schema for the data source.
func (d *devDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dev": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"engineers": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":    schema.StringAttribute{Computed: true},
									"name":  schema.StringAttribute{Computed: true},
									"email": schema.StringAttribute{Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure accepts provider configured data to set the API client.
func (d *devDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *client.Client for provider data, but received a different type.",
		)
		return
	}

	d.client = c
}

// Read refreshes the Terraform state with the latest data.
func (d *devDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DevDataSourceModel

	devs, err := d.client.GetDev()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Dev groups",
			err.Error(),
		)
		return
	}

	for _, dv := range devs {
		dvm := devDSModel{
			ID:   types.StringValue(dv.ID),
			Name: types.StringValue(dv.Name),
		}

		for _, eng := range dv.Engineers {
			dvm.Engineers = append(dvm.Engineers, devDSEngineer{
				ID:    types.StringValue(eng.ID),
				Name:  types.StringValue(eng.Name),
				Email: types.StringValue(eng.Email),
			})
		}

		state.Dev = append(state.Dev, dvm)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// DataSourceModel maps the data source schema data.
type DevDataSourceModel struct {
	Dev []devDSModel `tfsdk:"dev"`
}

// devModel maps Dev schema data.
type devDSModel struct {
	ID        types.String       `tfsdk:"id"`
	Name      types.String       `tfsdk:"name"`
	Engineers []devDSEngineer    `tfsdk:"engineers"`
}

type devDSEngineer struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

// DevInfoModel maps Dev info data
type DevInfoModel struct {
	ID types.String `tfsdk:"id"`
}
