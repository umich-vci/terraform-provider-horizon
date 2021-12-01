package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/umich-vci/gohorizon"
)

func resourceFarm() *schema.Resource {
	return &schema.Resource{
		Description: "Resource to manage Horizon Desktop Pools.",

		CreateContext: resourceDesktopPoolCreate,
		ReadContext:   resourceDesktopPoolRead,
		UpdateContext: resourceDesktopPoolUpdate,
		DeleteContext: resourceDesktopPoolDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the Desktop Pool. This property must contain only alphanumerics, underscores, and dashes.",
				Type:        schema.TypeString,
				Required:    true,
				// ValidateFunc: validation.StringMatch(),
			},
			"type": {
				Description:  "The type of the Desktop Pool.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AUTOMATED", "MANUAL", "RDS"}, false),
			},
		},
	}
}

func resourceDesktopPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	name := d.Get("name").(string)
	poolType := d.Get("type").(string)

	body := gohorizon.NewDesktopPoolCreateSpec(name, poolType)

	_, err := client.InventoryApi.CreateDesktopPool(ctx).Body(*body).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDesktopPoolRead(ctx, d, meta)
}

func resourceDesktopPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	id := d.Id()

	poolInfo, _, err := client.InventoryApi.GetDesktopPool(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", poolInfo.Name)
	d.Set("type", poolInfo.Type)

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
