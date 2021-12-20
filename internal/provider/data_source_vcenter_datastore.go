package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterDatastore() *schema.Resource {
	return &schema.Resource{
		Description: "Data source to find the ID of a vCenter Datacenter.",

		ReadContext: dataSourcevCenterDatastoreRead,

		Schema: map[string]*schema.Schema{
			"host_or_cluster_id": {
				Description: "Host or Cluster ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Datastore name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"capacity_mb": {
				Description: "Maximum capacity of this datastore, in MB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"datacenter_id": {
				Description: "Datacenter id for this datastore.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"disk_type": {
				Description: "Disk type of the datastore. SSD: Solid State Drive disk type. NON_SSD: NON-Solid State Drive disk type. UNKNOWN: Unknown disk type. NON_VMFS: NON-VMFS disk type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"file_system_type": {
				Description: "File system type of the datastore. VMFS: Virtual Machine File System. NFS: Network File System. VSAN: vSAN File System. VVOL: Virtual Volumes. UNKNOWN: Unknown File System type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"free_space_mb": {
				Description: "Available capacity of this datastore, in MB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"incompatible_reasons": {
				Description: "Reasons that may preclude this Datastore from being used in desktop pool/farm.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"local_datastore": {
				Description: "Indicates if this datastore is local to a single host.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"number_of_vms": {
				Description: "Indicates the number of virtual machines the datastore has for desktop pool/farm when applicable",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"path": {
				Description: "Datastore path.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vmfs_major_version": {
				Description: "The VMFS major version number.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterDatastoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	hcID := d.Get("host_or_cluster_id").(string)
	name := d.Get("name").(string)
	vCenterID := d.Get("vcenter_id").(string)

	datastores, _, err := client.ExternalApi.Listdatastores(ctx).HostOrClusterId(hcID).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, datastore := range datastores {
		if *datastore.Name == name {
			d.SetId(*datastore.Id)
			d.Set("capacity_mb", datastore.CapacityMb)
			d.Set("datacenter_id", datastore.DatacenterId)
			d.Set("disk_type", datastore.DiskType)
			d.Set("file_system_type", datastore.FileSystemType)
			d.Set("free_space_mb", datastore.FreeSpaceMb)
			d.Set("incompatible_reasons", datastore.IncompatibleReasons)
			d.Set("local_datastore", datastore.LocalDatastore)
			d.Set("number_of_vms", datastore.NumberOfVms)
			d.Set("path", datastore.Path)
			d.Set("vmfs_major_version", datastore.VmfsMajorVersion)

			return nil
		}
	}

	return diag.Errorf("could not find any datastore with name \"%s\"", name)
}
