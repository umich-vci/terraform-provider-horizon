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

	username := d.Get("username").(string)
	domainSID := d.Get("ad_domain_id").(string)

	accounts, _, err := client.ConfigApi.ListICDomainAccounts(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, account := range accounts {
		if *account.Username == username && *account.AdDomainId == domainSID {
			d.SetId(*account.Id)
			return nil
		}
	}

	return diag.Errorf("Could not find Instant Clone Domain Account with username %s and ad_domain_id %s", username, domainSID)
}
