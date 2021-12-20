package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenterBaseVMSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all the VM snapshots from the vCenter for a given VM.",

		ReadContext: dataSourcevCenterBaseVMSnapshotRead,

		Schema: map[string]*schema.Schema{
			"base_vm_id": {
				Description: "VM ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "VM snapshot path.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vcenter_id": {
				Description: "Virtual Center ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"created_timestamp": {
				Description: "Epoch time in milli seconds, when the VM snapshot was created.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"description": {
				Description: "Description of the VM snapshot.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"disk_size_mb": {
				Description: "Sum of capacities of all the virtual disks in the VM snapshot, in MB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"hardware_version": {
				Description: "VM snapshot hardware version",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"incompatible_reasons": {
				Description: "Reasons that may preclude this VM snapshot from being used in linked/instant clone desktop pool or farm creation.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"max_number_of_monitors": {
				Description: "Maximum number of monitors set in SVGA settings for the VM snapshot in vCenter.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_resolution_of_any_one_monitor": {
				Description: "Maximum resolution of any one monitor set in SVGA settings for the VM snapshot in vCenter.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"memory_mb": {
				Description: "The physical memory size of VM snapshot, in MB",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"memory_reservation_mb": {
				Description: "Amount of memory that is guaranteed available to the virtual machine, in MB.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"name": {
				Description: "VM snapshot name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"renderer3d": {
				Description: "Indicate how the virtual video device for the VM snapshot renders 3D graphics. Will be set only if VM snapshot supports 3D functions. MANAGE_BY_VSPHERE_CLIENT: 3D rendering managed by vSphere Client. AUTOMATIC: 3D rendering is automatic. SOFTWARE: 3D rendering is software dependent. The software renderer is supported (at minimum) on virtual hardware version 8 in a vSphere 5.0 environment. HARDWARE: 3D rendering is hardware dependent. The hardware-based renderer is supported (at minimum) on virtual hardware version 9 in a vSphere 5.1 environment. DISABLED: 3D rendering is disabled.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"total_video_memory_mb": {
				Description: "Total video memory in MB set in SVGA settings for the VM snapshot in vCenter.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"vgpu_type": {
				Description: "NVIDIA GRID vGPU type configured on this VM snapshot.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterBaseVMSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	bvmID := d.Get("base_vm_id").(string)
	path := d.Get("path").(string)
	vCenterID := d.Get("vcenter_id").(string)

	snapshots, _, err := client.ExternalApi.ListBaseSnapshots(ctx).BaseVmId(bvmID).VcenterId(vCenterID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, snapshot := range snapshots {
		if *snapshot.Path == path {
			d.SetId(*snapshot.Id)
			d.Set("created_timestamp", snapshot.CreatedTimestamp)
			d.Set("description", snapshot.Description)
			d.Set("disk_size_mb", snapshot.DiskSizeMb)
			d.Set("hardware_version", snapshot.HardwareVersion)
			d.Set("incompatible_reasons", snapshot.IncompatibleReasons)
			d.Set("max_number_of_monitors", snapshot.MaxNumberOfMonitors)
			d.Set("max_resolution_of_any_one_monitor", snapshot.MaxResolutionOfAnyOneMonitor)
			d.Set("memory_mb", snapshot.MemoryMb)
			d.Set("memory_reservation_mb", snapshot.MemoryReservationMb)
			d.Set("name", snapshot.Name)
			d.Set("renderer3d", snapshot.Renderer3d)
			d.Set("total_video_memory_mb", snapshot.TotalVideoMemoryMb)
			d.Set("vgpu_type", snapshot.VgpuType)

			return nil
		}
	}

	return diag.Errorf("could not find any Base VM Snapshot with path \"%s\"", path)
}
