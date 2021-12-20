package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceActiveDirectoryDomainUserOrGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Summary information related to AD Users or Groups. List API returning this summary information can use search filter query to filter on specific fields supported by filters. Supported Filters : 'And', 'Or', 'Equals', 'StartsWith', 'Contains'. See the field description to know the filter types it supports.",

		ReadContext: dataSourceActiveDirectoryDomainUserRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Description: "A JSON string containing the filter to use to find the Active Directory User or Group. The filter must find exactly 1 result or an error will be returned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_only": {
				Description: "If passed as \"true\", then only groups are returned. If passed as \"false\", then only users are returned. If not passed passed at all, then both types are returned",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"container": {
				Description: "AD container for this user or group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description number of this user or group. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_name": {
				Description: "Login name with domain of this user or group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"distinguished_name": {
				Description: "Active Directory distinguished name for this user or group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description: "DNS name of the domain in which this user or group belongs. Supported Filters : 'Equals'. Also, if 'Or' filter is used anywhere in filter string for this model class, then that 'Or' filter should nest only 'Equals' filter on 'domain' or 'id' field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Email address of this user or group. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"first_name": {
				Description: "First name of this user or group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group": {
				Description: "Indicates if this object represents a group.  This field is NOT supported in filter string. To use any filter on 'group', use 'group_only' query param directly.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"guid": {
				Description: "GUID of the user or group in RFC 4122 format. Supported Filters : 'Equals'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"kiosk_user": {
				Description: "Indicates if this user or group is a \"kiosk user\" that supports client authentication. Client authentication is the process of supporting client devices directly logging into resources.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"last_name": {
				Description: "Last name of this user or group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"login_name": {
				Description: "Login name of this user or group. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"long_display_name": {
				Description: "Login name, domain and name for this user or group, else display name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of this user or group. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"phone": {
				Description: "Phone number of this user. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_display_name": {
				Description: "User or group's display name. This corresponds with displayName attribute in AD.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_principal_name": {
				Description: "User Principal name(UPN) of this user. Supported Filters : 'Equals', 'StartsWith', 'Contains'.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceActiveDirectoryDomainUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).Client

	filter := d.Get("filter").(string)

	listUserOrGroup := client.ExternalApi.ListADUserOrGroupSummary(ctx)

	if g, ok := d.GetOk("group_only"); ok {
		listUserOrGroup.GroupOnly(strconv.FormatBool(g.(bool)))
	}

	entity, _, err := listUserOrGroup.Filter(filter).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	switch len(entity) {
	case 0:
		return diag.Errorf("could not find any user or group with filter\"%s\"", filter)
	case 1:
		d.SetId(*entity[0].Id)
		d.Set("container", entity[0].Container)
		d.Set("description", entity[0].Description)
		d.Set("display_name", entity[0].DisplayName)
		d.Set("distinguished_name", entity[0].DistinguishedName)
		d.Set("domain", entity[0].Domain)
		d.Set("email", entity[0].Email)
		d.Set("first_name", entity[0].FirstName)
		d.Set("group", entity[0].Group)
		d.Set("guid", entity[0].Guid)
		d.Set("kiosk_user", entity[0].KioskUser)
		d.Set("last_name", entity[0].LastName)
		d.Set("login_name", entity[0].LoginName)
		d.Set("long_display_name", entity[0].LongDisplayName)
		d.Set("name", entity[0].Name)
		d.Set("phone", entity[0].Phone)
		d.Set("user_display_name", entity[0].UserDisplayName)
		d.Set("user_principal_name", entity[0].UserPrincipalName)
		return nil
	default:
		return diag.Errorf("Multiple users/groups found with filter \"%s\"", filter)
	}
}
