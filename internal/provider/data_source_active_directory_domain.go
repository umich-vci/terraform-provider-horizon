package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceActiveDirectoryDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for reading information about an Active Directory domain from Horizon.",

		ReadContext: dataSourceActiveDirectoryDomainRead,

		Schema: map[string]*schema.Schema{
			"dns_name": {
				Description: "DNS name of the AD Domain.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ad_domain_auto_discovery": {
				Description: "Auto discovers domain controllers. Auto discovery, AD domain controllers and preferred site name are mutually exclusive. Only one of them can be defined at a time. Default value is true.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"ad_domain_context": {
				Description: "Active directory domain Context.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ad_domain_controllers": {
				Description: "One or more AD domain controllers. Auto discovery, AD domain controllers and preferred site name are mutually exclusive. Only one of them can be defined at a time.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ad_domain_preferred_site": {
				Description: "ADDomain preferred domain site. Auto discovery, AD domain controllers and preferred site name are mutually exclusive. Only one of them can be defined at a time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"auxiliary_accounts": {
				Description: "Auxiliary service accounts information of untrusted domain.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Unique SID representing auxiliary account.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username": {
							Description: "Auxiliary Service account username.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"domain_type": {
				Description: "AD Domain Type. CONNECTION_SERVER_DOMAIN: The domain having trust with connection server domain. NO_TRUST_DOMAIN: The domain not having any trust with connection server domain.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"netbios_name": {
				Description: "NetBIOS name of the AD Domain.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"port": {
				Description: "Port of the server to connect to.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"primary_account_password": {
				Description: "Information related to untrusted Domain service accounts. Service account user password.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"primary_account_username": {
				Description: "Information related to untrusted Domain service accounts. Service account username.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceActiveDirectoryDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	dnsName := d.Get("dns_name").(string)

	domains, _, err := client.ExternalApi.ListADDomainsV3(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, domain := range domains {
		if *domain.DnsName == dnsName {
			d.Set("domain_type", domain.DomainType)
			d.Set("netbios_name", domain.NetbiosName)

			auxAccounts := []map[string]interface{}{}
			if domain.AuxiliaryAccounts != nil {
				for _, acct := range *domain.AuxiliaryAccounts {
					auxAccount := make(map[string]interface{})
					auxAccount["id"] = acct.Id
					auxAccount["username"] = acct.Username
					auxAccounts = append(auxAccounts, auxAccount)
				}
			}
			d.Set("auxiliary_accounts", auxAccounts)

			if domain.PrimaryAccount != nil {
				d.Set("primary_account_password", domain.PrimaryAccount.Password)
				d.Set("primary_account_username", domain.PrimaryAccount.Username)
			}

			if domain.AdDomainAdvancedSettings != nil {
				d.Set("ad_domain_auto_discovery", domain.AdDomainAdvancedSettings.AdDomainAutoDiscovery)
				d.Set("ad_domain_context", domain.AdDomainAdvancedSettings.AdDomainContext)
				d.Set("ad_domain_controllers", domain.AdDomainAdvancedSettings.AdDomainControllers)
				d.Set("ad_domain_preferred_site", domain.AdDomainAdvancedSettings.AdDomainPreferredSite)
				d.Set("port", domain.AdDomainAdvancedSettings.Port)
			}

			d.SetId(*domain.Id)
			return nil
		}
	}

	return diag.Errorf("Could not find Active Directory Domain with dns_name %s", dnsName)
}
