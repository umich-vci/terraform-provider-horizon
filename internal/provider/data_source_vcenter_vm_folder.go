package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterVMFolder() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to find the ID of a vCenter Datacenter.",

		ReadContext: dataSourcevCenterVMFolderRead,

		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Description: "Datacenter ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "VM folder path.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"incompatible_reasons": {
				Description: "Reasons that may preclude this VM folder from being used in desktop pool or farm.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Description: "VM folder name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "VM folder type. DATACENTER: A datacenter that serves as a folder suitable for use in desktop pool/farm. FOLDER: A regular folder suitable for use in desktop pool/farm. OTHER: Other folder type that cannot be used in desktop pool/farm.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterVMFolderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	dcID := d.Get("datacenter_id").(string)
	path := d.Get("path").(string)
	vCenterID := d.Get("vcenter_id").(string)

	vmFolders, _, err := client.ExternalApi.ListVMFolders(ctx).DatacenterId(dcID).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, vmFolder := range vmFolders {
		if *vmFolder.Path == path {
			d.SetId(*vmFolder.Id)
			d.Set("incompatible_reasons", vmFolder.IncompatibleReasons)
			d.Set("name", vmFolder.Name)
			d.Set("type", vmFolder.Type)

			return nil
		}
	}

	return diag.Errorf("could not find any VM folder with path \"%s\"", path)
}
