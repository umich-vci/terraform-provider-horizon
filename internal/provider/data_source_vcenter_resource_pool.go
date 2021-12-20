package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterResourcePool() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to find the ID of a vCenter Datacenter.",

		ReadContext: dataSourcevCenterResourcePoolRead,

		Schema: map[string]*schema.Schema{
			"host_or_cluster_id": {
				Description: "Host or Cluster ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Resource pool name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "Resource pool path.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Resource pool type. HOST: Host used as a resource pool suitable for use in desktop pool/farm. CLUSTER: Cluster used as a resource pool suitable for use in desktop pool/farm. RESOURCE_POOL: Regular resource pool suitable for use in desktop pool/farm. OTHER: Other resource type which cannot be used in desktop pool/farm.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterResourcePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	hcID := d.Get("host_or_cluster_id").(string)
	name := d.Get("name").(string)
	vCenterID := d.Get("vcenter_id").(string)

	resourcePools, _, err := client.ExternalApi.ListResourcePools(ctx).HostOrClusterId(hcID).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, resourcePool := range resourcePools {
		if *resourcePool.Name == name {
			d.SetId(*resourcePool.Id)
			d.Set("path", resourcePool.Path)
			d.Set("type", resourcePool.Type)

			return nil
		}
	}

	return diag.Errorf("could not find any resource pool with name \"%s\"", name)
}
