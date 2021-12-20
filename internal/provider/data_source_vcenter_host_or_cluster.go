package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterHostOrCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to find the ID of a vCenter Datacenter.",

		ReadContext: dataSourcevCenterHostOrClusterRead,

		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Description: "Datacenter ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Host or cluster display name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster": {
				Description: "Whether or not this is a cluster or a host.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"incompatible_reasons": {
				Description: "Reasons that may preclude this Host Or Cluster from being used in desktop pool creation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"path": {
				Description: "Host or cluster path.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vgpu_types": {
				Description: "Types of NVIDIA GRID vGPUs supported by this host or at least one host on this cluster. If unset, this host or cluster does not support NVIDIA GRID vGPUs and cannot be used for desktop creation with NVIDIA GRID vGPU support enabled.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterHostOrClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	datacenterID := d.Get("datacenter_id").(string)
	name := d.Get("name").(string)
	vCenterID := d.Get("vcenter_id").(string)
	cluster := d.Get("cluster").(bool)

	hostsOrClusters, _, err := client.ExternalApi.ListHostsOrClusters(ctx).DatacenterId(datacenterID).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, hostOrCluster := range hostsOrClusters {
		if *hostOrCluster.Details.Name == name && *hostOrCluster.Details.Cluster == cluster {
			d.SetId(*hostOrCluster.Id)
			d.Set("incompatible_reasons", hostOrCluster.Details.IncompatibleReasons)
			d.Set("path", hostOrCluster.Details.Path)
			d.Set("vgpu_types", hostOrCluster.Details.VgpuTypes)

			return nil
		}
	}

	return diag.Errorf("could not find any host or cluster with name \"%s\"", name)
}
