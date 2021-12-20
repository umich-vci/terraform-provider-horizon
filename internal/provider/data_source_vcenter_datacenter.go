package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterDatacenter() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to find the ID of a vCenter Datacenter.",

		ReadContext: dataSourcevCenterDatacenterRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the datacenter.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "Datacenter path.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	vCenterID := d.Get("vcenter_id").(string)

	datacenters, _, err := client.ExternalApi.ListDatacenters(ctx).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, datacenter := range datacenters {
		if *datacenter.Name == name {
			d.SetId(*datacenter.Id)
			d.Set("path", datacenter.Path)
			return nil
		}
	}

	return diag.Errorf("could not find any datacenter with name \"%s\"", name)

}
