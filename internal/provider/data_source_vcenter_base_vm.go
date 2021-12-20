package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterBaseVM() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for base vm information.",

		ReadContext: dataSourcevCenterBaseVMRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "VM name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"datacenter_id": {
				Description: "Datacenter ID",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"filter_incompatible_vms": {
				Description: "Whether to filter out incompatible VMs",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"incompatible_reasons": {
				Description: "Reasons that may preclude this BaseVM from having its snapshots used in linked or instant clone desktop or farm creation.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network_type": {
				Description: "Type of network base VM belongs to. STANDARD_NETWORK: Standard network. OPAQUE_NETWORK: Opaque network. DISTRUBUTED_VIRTUAL_PORT_GROUP: DVS port group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"operating_system": {
				Description: "Operating system. UNKNOWN: Unknown WINDOWS_XP: Windows XP WINDOWS_VISTA: Windows Vista WINDOWS_7: Windows 7 WINDOWS_8: Windows 8 WINDOWS_10: Windows 10 WINDOWS_SERVER_2003: Windows Server 2003 WINDOWS_SERVER_2008: Windows Server 2008 WINDOWS_SERVER_2008_R2: Windows Server 2008 R2 WINDOWS_SERVER_2012: Windows Server 2012 WINDOWS_SERVER_2012_R2: Windows Server 2012 R2 WINDOWS_SERVER_2016_OR_ABOVE: Windows Server 2016 or above LINUX_OTHER: Linux (other) LINUX_SERVER_OTHER: Linux server (other) LINUX_UBUNTU: Linux (Ubuntu) LINUX_RHEL: Linux (Red Hat Enterprise) LINUX_SUSE: Linux (Suse) LINUX_CENTOS: Linux (CentOS)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"operating_system_display_name": {
				Description: "Operating system display name from Virtual Center.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterBaseVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	vCenterID := d.Get("vcenter_id").(string)
	listBaseVMs := client.ExternalApi.ListBaseVMs(ctx)

	if dcID, ok := d.GetOk("datacenter_id"); ok {
		listBaseVMs.DatacenterId(dcID.(string))
	}

	if filterIncompat, ok := d.GetOk("filter_incompatible_vms"); ok {
		listBaseVMs.FilterIncompatibleVms(filterIncompat.(bool))
	}

	baseVMs, _, err := listBaseVMs.VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, baseVM := range baseVMs {
		if *baseVM.Name == name {
			d.SetId(*baseVM.Id)
			d.Set("datacenter_id", baseVM.DatacenterId)
			d.Set("incompatible_reasons", baseVM.IncompatibleReasons)
			d.Set("network_type", baseVM.NetworkType)
			d.Set("operating_system", baseVM.OperatingSystem)
			d.Set("operating_system_display_name", baseVM.OperatingSystemDisplayName)
			d.Set("path", baseVM.Path)

			return nil
		}
	}

	return diag.Errorf("could not find any Base VM with name \"%s\"", name)
}
