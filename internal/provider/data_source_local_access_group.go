package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalAccessGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for reading information about a local access group from Horizon.",

		ReadContext: dataSourceLocalAccessGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Access group name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"deletable": {
				Description: "Indicates whether this access group can be deleted.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"description": {
				Description: "Access group description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLocalAccessGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)

	groups, _, err := client.ConfigApi.ListLocalAccessGroups(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, group := range groups {
		if *group.Name == name {
			d.Set("deletable", group.Deletable)
			d.Set("description", group.Description)
			d.SetId(*group.Id)
			return nil
		}
	}

	return diag.Errorf("Could not find Local Access Group with name %s", name)
}
