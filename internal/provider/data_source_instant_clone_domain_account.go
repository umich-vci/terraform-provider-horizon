package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstantCloneDomainAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for reading information about an instant clone domain account from Horizon.",

		ReadContext: dataSourceInstantCloneDomainAccountRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Description: "User name of the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ad_domain_id": {
				Description: "SID of the AD Domain that this account user belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceInstantCloneDomainAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
