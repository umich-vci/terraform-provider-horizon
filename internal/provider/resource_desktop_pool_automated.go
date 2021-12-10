package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gohorizon"
)

func resourceDesktopPoolAutomated() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage Horizon Desktop Pools.",

		CreateContext: resourceDesktopPoolCreate,
		ReadContext:   resourceDesktopPoolRead,
		UpdateContext: resourceDesktopPoolUpdate,
		DeleteContext: resourceDesktopPoolDelete,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Description: "Access groups can organize the entities such as desktop pools in the organization. They can also be used for delegated administration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Name of the Desktop Pool. This property must contain only alphanumerics, underscores, and dashes.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				// ValidateFunc: validation.StringMatch(),
			},
			"naming_method": {
				Description:  "Naming method for the desktop pool.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SPECIFIED", "PATTERN"}, false),
			},
			"user_assignment": {
				Description:  "User assignment scheme. DEDICATED: With dedicated assignment, a user returns to the same machine at each session. FLOATING: With floating assignment, a user may return to one of the available machines for the next session.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEDICATED", "FLOATING"}, false),
			},
			"vcenter_id": {
				Description: "ID of the virtual center server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"allow_multiple_user_assignments": {
				Description: "Only applies to automated desktop pools with manual user assignment. Whether assignment of multiple users to a single machine is allowed. If this is true then automatic_user_assignment should be false. ",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"automatic_user_assignment": {
				Description: "Automatic assignment of a user the first time they access the machine. This property is applicable if user_assignment is set to DEDICATED with default value as true.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"category_folder_name": {
				Description: "Name of the category folder in the user's OS containing a shortcut to the desktop pool. Will be unset if the desktop does not belong to a category.This property defines valid folder names with a max length of 64 characters and up to 4 subdirectory levels.The subdirectories can be specified using a backslash, e.g. (dir1\\dir2\\dir3\\dir4). Folder names can't start orend with a backslash nor can there be 2 or more backslashes together. Combinations such as(\\dir1, dir1\\dir2, dir1\\\\dir2, dir1\\\\\\dir2) are invalid. The windows reserved keywords(CON, PRN, NUL, AUX, COM1 - COM9, LPT1 - LPT9 etc.) are not allowed in subdirectory names.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_assigned": {
				Description: "Indicates whether this desktop is assigned to a workspace in Horizon Cloud Services. This can be set to true from cloud session only and only when cloud_managed is set to true.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"cloud_managed": {
				Description: "Indicates whether this desktop is managed by Horizon Cloud Services. This can be set to false only when cloud_assigned is set to false. Default value is false. This property cannot be set to true, if any of the conditions are satisfied: user is provided. enabled is false. supported_session_type is not DESKTOP. global_entitlement is set. user_assignment is DEDICATED and automatic_user_assignment is false. Local entitlements are configured. Any of the machines in the pool have users assigned. cs_restriction_tags is not set. Desktop pool type is MANUAL.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"cs_restriction_tags": {
				Description: "List of Connection server restriction tags to which the access to the desktop pool is restricted. If this property is not set it indicates that desktop pool can be accessed from any connection server.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clone_prep_settings": {
				Description: "ClonePrep is a VMware system tool executed by Instant Clone Engine during a instant clone machine deployment. ClonePrep personalizes each machine created from the Master image.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_container_rdn": {
							Description: "Instant Clone Engine Active Directory container for ClonePrep.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"instant_clone_domain_account_id": {
							Description: "This is required for instant clone desktop pools. This is the administrator which will add the machines to its domain upon creation.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"priming_computer_account": {
							Description: "Instant Clone publishing needs an additional computer account in the same AD domain as the clones. This field accepts the pre-created computer accounts.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"post_synchronization_script_name": {
							Description: "Post synchronization script. ClonePrep can run a customization script on instant-clone machines after they are created or recovered or a new image is pushed. Provide the path to the script on the parent virtual machine.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"post_synchronization_script_parameters": {
							Description: "Post synchronization script parameters.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"power_off_script_name": {
							Description: "Power off script. ClonePrep can run a customization script on instant-clone machines before they are powered off. Provide the path to the script on the parent virtual machine.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"power_off_script_parameters": {
							Description: "Power off script parameters.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"reuse_pre_existing_accounts": {
							Description: "Indicates whether to allow the use of existing AD computer accounts when the VM names of newly created clones match the existing computer account names.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"sys_prep_settings": {
				Description: "Microsoft Sysprep is a tool to deploy the configured operating system installation from a base image. The machine can then be customized based on an answer script. Sysprep can modify a larger number of configurable parameters than QuickPrep.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysprep_customization_spec_id": {
							Description: "This is required when customization_type is set as SYS_PREP. Customization specification to use when Sysprep customization is requested.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"description": {
				Description: "Description of the desktop pool.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"display_assigned_machine_name": {
				Description: "Applicable To: Dedicated desktop pools with default value as false. Indicates whether users should see the hostname of the machine assigned to them instead of display_name when they connect using Horizon Client. If no machine is assigned to the user then \"display_name (No machine assigned)\" will be displayed in the client.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"display_machine_alias": {
				Description: "Applicable To: Dedicated desktop pools with default value as false. If no machine is assigned to the user then \"displayName No machine assigned)\" will be displayed in the Horizon client. If both display_assigned_machine_name and this property is set to true, machine alias of the assigned machine is displayed if the user has machine alias set. Otherwise hostname will be displayed.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"display_name": {
				Description: "Display name of the desktop pool. If the display name is left blank, it defaults to name.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"display_protocol_settings": {
				Description: "Display protocol settings.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_users_to_choose_protocol": {
							Description: "Indicates whether the users can choose the protocol.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"default_display_protocol": {
							Description:  "The default display protocol for the desktop pool.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "PCOIP",
							ValidateFunc: validation.StringInSlice([]string{"RDP", "PCOIP", "BLAST"}, false),
						},
						"grid_vgpus_enabled": {
							Description: "When 3D rendering is managed by the vSphere Client, this enables support for NVIDIA GRID vGPUs. This will be false if 3D rendering is not managed by the vSphere Client. If this is true, the host or cluster associated with the desktop pool must support NVIDIA GRID and vGPU types required by the desktop pool's VirtualMachines, VmTemplate or BaseImageSnapshot. If this is false, the desktop pool's VirtualMachines, VmTemplate or BaseImageSnapshot must not support NVIDIA GRID vGPUs. Since suspending VMs with passthrough devices such as vGPUs is not possible, power_policy cannot be set to SUSPEND if this is enabled. Default value is false.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"max_number_of_monitors": {
							Description: "When render3D is disabled, the max_number_of_monitors and max_resolution_of_any_one_monitor settings determine the amount of vRAM assigned to machines in this desktop. The greater these values are, the more memory will be consume on the associated ESX hosts. Existing virtual machines must be powered off and subsequently powered on for the change to take effect. A restart will not cause the changes to take effect. If 3D is enabled and managed by View, the maximum number of monitors must be 1 or 2. For Instant Clones, this value is inherited from snapshot of Master VM. This property is required if renderer3D is set to AUTOMATIC, SOFTWARE, HARDWARE or DISABLED.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     2,
						},
						"max_resolution_of_any_one_monitor": {
							Description:  "If 3D rendering is enabled and managed by View, this must be set to the default value. When 3D rendering is disabled, the max_number_of_monitors and max_resolution_of_any_one_monitor settings determine the amount of vRAM assigned to machines in this desktop. The greater these values are, the more memory will be consumed on the associated ESX hosts. This setting is only relevant on managed machines. Existing virtual machines must be powered off and subsequently powered on for the change to take effect. A restart will not cause the changes to take effect. For Instant Clones, this value is inherited from snapshot of Master VM.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "WUXGA",
							ValidateFunc: validation.StringInSlice([]string{"WSXGA_PLUS", "WUXGA", "WQXGA", "UHD", "UHD_5K", "UHD_8K"}, false),
						},
						"renderer_3d": {
							Description:  "3D rendering is supported on Windows 7 or later guests running on VMs with virtual hardware version 8 or later. The default_display_protocol must set to PCOIP and allow_users_to_choose_protocol must be set to false to enable 3D rendering. For instant clone source desktop 3D rendering always mapped to MANAGE_BY_VSPHERE_CLIENT.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "DISABLED",
							ValidateFunc: validation.StringInSlice([]string{"MANAGE_BY_VSPHERE_CLIENT", "AUTOMATIC", "SOFTWARE", "HARDWARE", "DISABLED"}, false),
						},
						"session_collaboration_enabled": {
							Description: "Enable session collaboration feature. Session collaboration allows a user to share their remote session with other users. BLAST must be configured as a supported protocol in supported_display_protocols.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"vram_size_mb": {
							Description: "vRAM size for View managed 3D rendering. More VRAM can improve 3D performance. Size is in MB. On ESXi 5.0 hosts, the renderer allows a maximum VRAM size of 128MB. On ESXi 5.1 and later hosts, the maximum VRAM size is 512MB. For Instant Clones, this value is inherited from snapshot of Master VM. This property is required if renderer_3d is set to AUTOMATIC, SOFTWARE or HARDWARE.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     96,
						},
					},
				},
			},
			"do_not_power_on_vms_after_creation": {
				Description: "Indicates whether to power on VMs after creation. This is the settings when customization will be done manually. This property is required if customization_type is set to NONE with default value as false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"enable_client_restrictions": {
				Description: "Client restrictions to be applied to the desktop pool.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"enable_provisioning": {
				Description: "Indicates whether provisioning is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"enabled": {
				Description: "Indicates whether the desktop pool is enabled for brokering.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			// "nics": {},
			"pattern_naming_settings": {
				Description: "Naming pattern settings for Automated desktop pool.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"naming_pattern": {
							Description: "Virtual machines will be named according to the specified naming pattern. By default, view manager appends a unique number to the specified pattern to provide a unique name for each virtual machine. To place this unique number elsewhere in the pattern, use '{n}'. (For example: vm-{n}-sales.) The unique number can also be made a fixed length. (For example: vm-{n:fixed=3}-sales will name VMs from vm-001-sales to vm-999-sales). Machine names are constrained to a maximum size of 15 characters including the unique number. Therefore, care must be taken when choosing a pattern. If the maximum desktop size is 9 machines, the pattern must be at most 14 characters. For 99 machines, 13 characters, for 999 machines, 12 characters. For 9999 machines, 11 characters. If using a fixed size token, use a maximum of 14 characters for \"n=1\", 13 characters for \"n=2\", 12 characters for \"n=3\", and 11 characters for \"n=4\". If {n} is specified with no size, a size of 2 is automatically used and if no {} is specified, {n=2} is automatically appended to the end of the pattern.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"provisioning_time": {
							Description:  "Determines when the machines are provisioned. ON_DEMAND: Provision machines on demand. UP_FRONT: Provision all machines up-front.",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "UP_FRONT",
							ValidateFunc: validation.StringInSlice([]string{"UP_FRONT", "ON_DEMAND"}, false),
						},
						"max_number_of_machines": {
							Description: "Maximum number of machines in the desktop pool.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
						},
						"min_number_of_machines": {
							Description: "This is applicable if provisioning_time is set to ON_DEMAND with default value of 0.",
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
						},
						"number_of_spare_machines": {
							Description:  "Number of spare powered on machines.",
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
			"provisioning_settings": {
				Description: "Virtual center provisioning settings for Automated desktop pool.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_or_cluster_id": {
							Description: "Host or cluster where the machines are deployed in. For Instant clone desktops it can only be set to a cluster id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"resource_pool_id": {
							Description: "Resource pool to deploy the machines.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"vm_folder_id": {
							Description: "VM folder where the machines are deployed to.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"add_virtual_tpm": {
							Description: "Indicates whether to add Virtual TPM device.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
						},
						"base_snapshot_id": {
							Description: "This property can be set only when source is set to INSTANT_CLONE, vm_template_id is unset and parent_vm_id is set.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"datacenter_id": {
							Description: "Datacenter within which the desktop pool is configured.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"im_stream_id": {
							Description: "This is required when vm_template_id, parent_vm_id and base_snapshot_id are not set.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"im_tag_id": {
							Description: "This is required when im_stream_id is set.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"parent_vm_id": {
							Description: "This property can be set only when source is set to INSTANT_CLONE.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"vm_template_id": {
							Description: "Applicable To: Full clone desktop pools. This is required if parent_vm_id and base_snapshot_id are not set.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			// "session_settings": {},
			"session_type": {
				Description:  "Supported session types for this desktop pool. If this property is set to APPLICATION then this desktop pool can be used for application pool creation. This will be useful when the machines in the pool support application remoting.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DESKTOP", "APPLICATION", "DESKTOP_AND_APPLICATION"}, false),
				Default:      "DESKTOP",
			},
			"shortcut_locations_v2": {
				Description: "Locations of the category folder in the user's OS containing a shortcut to the desktop pool. This is required if the category_folder_name is set.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				//ValidateFunc: validation.StringInSlice([]string{"START_MENU", "DESKTOP"}, false),
			},
			"source": {
				Description:  "Source of the Machines in this Desktop Pool.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"INSTANT_CLONE", "VIRTUAL_CENTER"}, false),
			},
			// "specific_naming_settings": {},
			"stop_provisioning_on_error": {
				Description: "Disable provisioning on the pool if there is a provisioning error.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"storage_settings": {
				Description: "Virtual center storage settings for Automated desktop pool.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datastores": {
							Description: "Datastores to store the machine.",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datastore_id": {
										Description: "Id of the datastore.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"sdrs_cluster": {
										Description: "Id of the datastore.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
								},
							},
						},
						"reclaim_vm_disk_space": {
							Description: "With vSphere 5.x, virtual machines can be configured to use a space efficient disk format that supports reclamation of unused diskspace (such as deleted files). This option reclaims unused diskspace on each virtual machine. The operation is initiated when an estimate of used disk space exceeds the specified threshold.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"reclamation_threshold_mb": {
							Description: "Initiate reclamation when unused space on virtual machine exceeds the threshold in MB.  This property is required if reclaim_vm_disk_space is set to true.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"replica_disk_datastore_id": {
							Description: "Datastore to store replica disks for instant clone machines. This property is required if use_separate_datastores_replica_and_os_disks is set to true.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"use_separate_datastores_replica_and_os_disks": {
							Description: "Indicates whether to use separate datastores for replica and OS disks.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"use_vsan": {
							Description: "Indicates whether to use vSphere vSAN.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"transparent_page_sharing_scope": {
				Description:  "Transparent page sharing scope for this Desktop Pool. VM: Inter-VM page sharing is not permitted. DESKTOP_POOL: Inter-VM page sharing among VMs belonging to the same Desktop pool is permitted. POD: Inter-VM page sharing among VMs belonging to the same Pod is permitted. GLOBAL: Inter-VM page sharing among all VMs on the same host is permitted.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "VM",
				ValidateFunc: validation.StringInSlice([]string{"VM", "DESKTOP_POOL", "POD", "GLOBAL"}, false),
			},
			"view_storage_accelerator_settings": {
				Description: "View Storage Accelerator settings for Managed desktop pool.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blackout_times": {
							Description: "Fields for specifying blackout time for View Storage Accelerator. Storage accelerator regeneration and VM disk space reclamation do not occur during blackout times. The same blackout policy applies to both operations.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Description: "List of days for a given range of time.",
										Type:        schema.TypeList,
										Required:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"end_time": {
										Description: "Ending time for the blackout in 24-hour format.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"start_time": {
										Description: "Starting time for the blackout in 24-hour format.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
						"regenerate_view_storage_accelerator_days": {
							Description: "How often to regenerate the View Storage Accelerator cache. Measured in Days. This property is required if useViewStorageAccelerator is set to true. ",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     7,
						},
						"use_view_storage_accelerator": {
							Description: "Indicates whether to use View Storage Accelerator. ",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"view_storage_accelerator_disk_types": {
							Description: "Disk types to enable for the View Storage Accelerator feature. Not applicable to full clone pools",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"delete_in_progress": {
				Description: "Indicates whether the desktop pool is in the process of being deleted.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"image_source": {
				Description: "Source of image used in the desktop pool. Possible values are VIRTUAL_CENTER: Image was created in virtual center. IMAGE_CATALOG: Image was created in image catalog.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_group_count": {
				Description: "Count of user or group entitlements present for the desktop pool.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceDesktopPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	poolType := "AUTOMATED"
	source := d.Get("source").(string)
	userAssignment := d.Get("user_assignment").(string)
	enableProvisioning := d.Get("enable_provisioning").(bool)
	stopOnErr := d.Get("stop_provisioning_on_error").(bool)
	namingMethod := d.Get("naming_method").(string)
	vCenterID := d.Get("vcenter_id").(string)
	agID := d.Get("access_group_id").(string)

	//access_group_id

	body := gohorizon.NewDesktopPoolCreateSpec(name, poolType)
	body.Source = &source
	body.UserAssignment = &userAssignment
	body.EnableProvisioning = &enableProvisioning
	body.StopProvisioningOnError = &stopOnErr
	body.NamingMethod = &namingMethod
	body.VcenterId = &vCenterID
	body.AccessGroupId = &agID

	autoAssign := false
	if aa, ok := d.GetOk("automatic_user_assignment"); ok {
		if userAssignment == "FLOATING" {
			return diag.Errorf("automatic_user_assignment should not be set when user_assignment is \"FLOATING\"")
		}
		autoAssign = aa.(bool)
		body.AutomaticUserAssignment = &autoAssign
	}

	multiAssign := false
	if ma, ok := d.GetOk("allow_multiple_user_assignments"); ok {
		if userAssignment == "FLOATING" {
			return diag.Errorf("allow_multiple_user_assignments should not be set when user_assignment is \"FLOATING\"")
		}
		multiAssign = ma.(bool)
		body.AllowMultipleUserAssignments = &multiAssign
	}

	if autoAssign && multiAssign {
		return diag.Errorf("automatic_user_assignment and allow_multiple_user_assignments cannot both be true.")
	}

	if namingMethod == "PATTERN" {
		if pns, ok := d.GetOk("pattern_naming_settings"); ok {
			patternNamingRaw := pns.([]interface{})[0].(map[string]interface{})
			namingPattern := patternNamingRaw["naming_pattern"].(string)
			patternNaming := gohorizon.NewDesktopPoolVirtualMachinePatternNamingSettingsCreateSpec(namingPattern)
			provTime := patternNamingRaw["provisioning_time"].(string)
			patternNaming.ProvisioningTime = &provTime
			maxMachine := int32(patternNamingRaw["max_number_of_machines"].(int))
			patternNaming.MaxNumberOfMachines = &maxMachine

			if patternNamingRaw["min_number_of_machines"].(int) > 0 {
				if namingPattern == "UP_FRONT" {
					return diag.Errorf("min_number_of_machines can not be set when naming_pattern is \"UP_FRONT\"")
				}
				minMachine := int32(patternNamingRaw["min_number_of_machines"].(int))
				patternNaming.MinNumberOfMachines = &minMachine
			}

			if spare, ok := patternNamingRaw["number_of_spare_machines"]; ok {
				if namingPattern == "UP_FRONT" {
					return diag.Errorf("number_of_spare_machines can not be set when naming_pattern is \"UP_FRONT\"")
				}
				numSpare := int32(spare.(int))
				patternNaming.NumberOfSpareMachines = &numSpare
			}

			body.PatternNamingSettings = patternNaming

		} else {
			return diag.Errorf("pattern_naming_settings must be set if naming_method is \"PATTERN\"")
		}
	}

	provSettingsRaw := d.Get("provisioning_settings").([]interface{})[0].(map[string]interface{})
	hcID := provSettingsRaw["host_or_cluster_id"].(string)
	rpID := provSettingsRaw["resource_pool_id"].(string)
	folderID := provSettingsRaw["vm_folder_id"].(string)
	provSettings := gohorizon.NewDesktopPoolProvisioningSettingsCreateSpec(hcID, rpID, folderID)
	vTPM := provSettingsRaw["add_virtual_tpm"].(bool)
	provSettings.AddVirtualTpm = &vTPM

	if dcid, ok := provSettingsRaw["datacenter_id"]; ok {
		dcID := dcid.(string)
		provSettings.DatacenterId = &dcID
	}

	switch source {
	case "INSTANT_CLONE":
		var icErr diag.Diagnostics

		if provSettingsRaw["parent_vm_id"].(string) != "" {
			pvmid := provSettingsRaw["parent_vm_id"].(string)
			provSettings.ParentVmId = &pvmid
		} else {
			icErr = append(icErr, diag.Errorf("parent_vm_id must be set when source is \"INSTANT_CLONE\"")...)
		}

		if provSettingsRaw["base_snapshot_id"].(string) != "" {
			basesnapID := provSettingsRaw["base_snapshot_id"].(string)
			provSettings.BaseSnapshotId = &basesnapID
		} else {
			icErr = append(icErr, diag.Errorf("base_snapshot_id must be set when source is \"INSTANT_CLONE\"")...)
		}

		if provSettingsRaw["vm_template_id"].(string) != "" {
			icErr = append(icErr, diag.Errorf("vm_template_id must not be set when source is \"INSTANT_CLONE\"")...)
		}

		if icErr != nil {
			return icErr
		}
	case "VIRTUAL_CENTER":
		var fcErr diag.Diagnostics

		if provSettingsRaw["vm_template_id"].(string) != "" {
			templateID := provSettingsRaw["vm_template_id"].(string)
			body.ProvisioningSettings.VmTemplateId = &templateID
		} else {
			fcErr = append(fcErr, diag.Errorf("vm_template_id must be set when source is \"VIRTUAL_CENTER\"")...)
		}

		if provSettingsRaw["parent_vm_id"].(string) != "" {
			fcErr = append(fcErr, diag.Errorf("parent_vm_id must not be set when source is \"VIRTUAL_CENTER\"")...)
		}

		if provSettingsRaw["base_snapshot_id"] != "" {
			fcErr = append(fcErr, diag.Errorf("base_snapshot_id must not be set when source is \"VIRTUAL_CENTER\"")...)
		}

		if fcErr != nil {
			return fcErr
		}
	default:
		return diag.Errorf("invalid source - should not be possible to get here")
	}

	body.ProvisioningSettings = provSettings

	resp, err := client.InventoryApi.CreateDesktopPool(ctx).Body(*body).Execute()
	if err != nil {
		return returnResponseErr(resp, err)
	}

	return resourceDesktopPoolRead(ctx, d, meta)
}

func resourceDesktopPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	id := d.Id()

	poolInfo, _, err := client.InventoryApi.GetDesktopPoolV5(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("access_group_id", poolInfo.AccessGroupId)
	d.Set("allow_multiple_user_assignments", poolInfo.AllowMultipleUserAssignments)
	d.Set("automatic_user_assignment", poolInfo.AutomaticUserAssignment)
	d.Set("category_folder_name", poolInfo.CategoryFolderName)
	d.Set("cloud_assigned", poolInfo.CloudAssigned)
	d.Set("cloud_managed", poolInfo.CloudManaged)
	d.Set("delete_in_progress", poolInfo.DeleteInProgress)
	d.Set("display_assigned_machine_name", poolInfo.DisplayAssignedMachineName)
	d.Set("display_machine_alias", poolInfo.DisplayMachineAlias)
	d.Set("display_name", poolInfo.DisplayName)
	d.Set("enable_client_restrictions", poolInfo.EnableClientRestrictions)
	d.Set("enable_provisioning", poolInfo.EnableProvisioning)
	d.Set("enabled", poolInfo.Enabled)
	d.Set("image_source", poolInfo.ImageSource)
	d.Set("name", poolInfo.Name)
	d.Set("naming_method", poolInfo.NamingMethod)
	d.Set("session_type", poolInfo.SessionType)
	d.Set("shortcut_locations_v2", poolInfo.ShortcutLocationsV2)
	d.Set("source", poolInfo.Source)
	d.Set("stop_provisioning_on_error", poolInfo.StopProvisioningOnError)
	d.Set("transparent_page_sharing_scope", poolInfo.TransparentPageSharingScope)
	d.Set("type", poolInfo.Type)
	d.Set("user_assignment", poolInfo.UserAssignment)
	d.Set("user_group_count", poolInfo.UserGroupCount)
	d.Set("vcenter_id", poolInfo.VcenterId)

	return nil
}

func resourceDesktopPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceDesktopPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	id := d.Id()

	_, err := client.InventoryApi.DeleteDesktopPool(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
