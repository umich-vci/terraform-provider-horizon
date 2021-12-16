package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcevCenter() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for reading information about a vCenter from Horizon.",

		ReadContext: dataSourcevCenterRead,

		Schema: map[string]*schema.Schema{
			"server_name": {
				Description: "Virtual Center's server name or IP address.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// "certificate_override": {
			// 	Description: "Certificate details and type information, which can be used to override thumbprint details.",
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// },
			"deployment_type": {
				Description: "Indicates different environments that Horizon can be deployed into. GENERAL: Horizon is deployed on On-premises. AZURE: Horizon is deployed on Azure. AWS: Horizon is deployed on AWS. DELL_EMC: Horizon is deployed on Dell EMC. GOOGLE: Horizon is deployed on Google Cloud. ORACLE: Horizon is deployed on Oracle Cloud.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Human readable description of the Virtual Center instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_name": {
				Description: "Human readable name of the Virtual Center instance.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Indicates if the virtual center is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"has_virtual_tpm_pools": {
				Description: "Indicates if there is any instant clone Desktop pool associated with this Virtual Center which has addVirtualTPM set",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"instance_uuid": {
				Description: "Virtual center's instanceUuid.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// "limits": {
			// 	Description: "Information about the limits configured for Virtual Center",
			// 	Type:        schema.TypeSet,
			// 	Computed:    true,
			// },
			"maintenance_mode": {
				Description: "Indicates if maintenance or upgrade task is scheduled on Virtual center or hosts",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"port": {
				Description: "Port of the virtual center to connect to.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"se_sparse_reclamation_enabled": {
				Description: "Indicates if Storage Efficiency Sparse (seSparse) reclamation is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			// "storage_accelerator_data": {
			// 	Description: "Information about the Storage Accelerator Data",
			// 	Type:        schema.TypeSet,
			// 	Computed:    true,
			// },
			"use_ssl": {
				Description: "Indicates if SSL should be used when connecting to the server.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"user_name": {
				Description: "User name to use for the connection.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"version": {
				Description: "Version of the Virtual Center.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcevCenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	serverName := d.Get("server_name").(string)

	vCenters, _, err := client.ConfigApi.ListVCInfoV2(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, vCenter := range vCenters {
		if *vCenter.ServerName == serverName {
			d.Set("deployment_type", vCenter.DeploymentType)
			d.Set("description", vCenter.Description)
			d.Set("display_name", vCenter.DisplayName)
			d.Set("enabled", vCenter.Enabled)
			d.Set("has_virtual_tpm_pools", vCenter.HasVirtualTpmPools)
			d.Set("instance_uuid", vCenter.InstanceUuid)
			d.Set("maintenance_mode", vCenter.MaintenanceMode)
			d.Set("port", vCenter.Port)
			d.Set("se_sparse_reclamation_enabled", vCenter.SeSparseReclamationEnabled)
			d.Set("use_ssl", vCenter.UseSsl)
			d.Set("user_name", vCenter.UserName)
			d.Set("version", vCenter.Version)
			d.SetId(*vCenter.Id)
			return nil
		}
	}

	return diag.Errorf("Could not find vCenter server with server_name %s", serverName)
}
