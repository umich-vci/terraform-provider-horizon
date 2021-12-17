package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gohorizon"
)

func resourceDesktopPoolEntitlements() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for managing desktop pool entitlements in Horizon.",

		CreateContext: resourceDesktopPoolEntitlementsCreate,
		ReadContext:   resourceDesktopPoolEntitlementsRead,
		UpdateContext: resourceDesktopPoolEntitlementsUpdate,
		DeleteContext: resourceDesktopPoolEntitlementsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ad_user_or_group_ids": {
				Description: "List of ad-user-or-group SIDs for the entitlement operations on the given desktop pool.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pool_id": {
				Description: "Unique ID representing the desktop pool.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}
func resourceDesktopPoolEntitlementsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	poolID := d.Get("pool_id").(string)

	adIDs := d.Get("ad_user_or_group_ids").(*schema.Set).List()

	adIDsRaw := []string{}
	for _, adID := range adIDs {
		adIDstr := adID.(string)
		adIDsRaw = append(adIDsRaw, adIDstr)
	}

	bodyElem := gohorizon.NewEntitlementSpec()
	bodyElem.Id = &poolID
	bodyElem.AdUserOrGroupIds = &adIDsRaw
	body := []gohorizon.EntitlementSpec{*bodyElem}

	_, _, err := client.EntitlementsApi.BulkCreateDesktopPoolEntitlements(ctx).Body(body).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	return resourceDesktopPoolEntitlementsRead(ctx, d, meta)
}

func resourceDesktopPoolEntitlementsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	poolID := d.Id()

	entitlement, _, err := client.EntitlementsApi.GetDesktopPoolEntitlements(ctx, poolID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("pool_id", poolID)
	d.Set("ad_user_or_group_ids", entitlement.AdUserOrGroupIds)

	return nil
}

func resourceDesktopPoolEntitlementsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	poolID := d.Id()
	adIDs := d.Get("ad_user_or_group_ids").(*schema.Set).List()

	adIDsRaw := []string{}
	for _, adID := range adIDs {
		adIDstr := adID.(string)
		adIDsRaw = append(adIDsRaw, adIDstr)
	}

	currentADIDs, _, err := client.EntitlementsApi.GetDesktopPoolEntitlements(ctx, poolID).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	removeList := []string{}
	for _, currentADID := range *currentADIDs.AdUserOrGroupIds {
		found := false
		for _, adID := range adIDsRaw {
			if currentADID == adID {
				found = true
				break
			}
		}
		if !found {
			removeList = append(removeList, currentADID)
		}
	}

	addList := []string{}
	for _, adID := range adIDsRaw {
		found := false
		for _, currentADID := range *currentADIDs.AdUserOrGroupIds {
			if adID == currentADID {
				found = true
				break
			}
		}
		if !found {
			addList = append(addList, adID)
		}
	}

	if len(removeList) > 0 {
		removeBodyElem := gohorizon.NewEntitlementSpec()
		removeBodyElem.Id = &poolID
		removeBodyElem.AdUserOrGroupIds = &removeList
		removeBody := []gohorizon.EntitlementSpec{*removeBodyElem}
		_, _, err := client.EntitlementsApi.BulkDeleteDesktopPoolEntitlements(ctx).Body(removeBody).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addList) > 0 {
		addBodyElem := gohorizon.NewEntitlementSpec()
		addBodyElem.Id = &poolID
		addBodyElem.AdUserOrGroupIds = &addList
		addBody := []gohorizon.EntitlementSpec{*addBodyElem}
		_, _, err := client.EntitlementsApi.BulkCreateDesktopPoolEntitlements(ctx).Body(addBody).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDesktopPoolEntitlementsRead(ctx, d, meta)
}

func resourceDesktopPoolEntitlementsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	poolID := d.Get("pool_id").(string)

	bodyElem := gohorizon.NewEntitlementSpec()
	bodyElem.Id = &poolID
	body := []gohorizon.EntitlementSpec{*bodyElem}
	_, _, err := client.EntitlementsApi.BulkDeleteDesktopPoolEntitlements(ctx).Body(body).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
