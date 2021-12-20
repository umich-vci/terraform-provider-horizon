package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umich-vci/gohorizon"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("HORIZON_USERNAME", nil),
					Description: "This is the username to use to access the VMware Horizon server. This must be provided in the config or in the environment variable `HORIZON_USERNAME`.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("HORIZON_PASSWORD", nil),
					Description: "This is the password to use to access the VMware Horizon server. This must be provided in the config or in the environment variable `HORIZON_PASSWORD`.",
				},
				"domain": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("HORIZON_DOMAIN", nil),
					Description: "This is the AD Domain of the `username` used to access the VMware Horizon server. This must be provided in the config or in the environment variable `HORIZON_USERNAME`.",
				},
				"horizon_host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("HORIZON_HOST", nil),
					Description: "This is the hostname or IP address of the VMware Horizon server. This must be provided in the config or in the environment variable `HORIZON_HOST`.",
				},
				"ssl_verify": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Verify the SSL certificate of the VMware Horizon server?",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"horizon_active_directory_domain":               dataSourceActiveDirectoryDomain(),
				"horizon_active_directory_domain_user_or_group": dataSourceActiveDirectoryDomainUserOrGroup(),
				"horizon_instant_clone_domain_account":          dataSourceInstantCloneDomainAccount(),
				"horizon_local_access_group":                    dataSourceLocalAccessGroup(),
				"horizon_vcenter_datacenter":                    dataSourcevCenterDatacenter(),
				"horizon_vcenter_server":                        dataSourcevCenter(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"horizon_desktop_pool_automated":    resourceDesktopPoolAutomated(),
				"horizon_desktop_pool_entitlements": resourceDesktopPoolEntitlements(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	Client gohorizon.APIClient
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-horizon", version)

		username := d.Get("username").(string)
		password := d.Get("password").(string)
		domain := d.Get("domain").(string)
		host := d.Get("horizon_host").(string)
		sslVerify := d.Get("ssl_verify").(bool)

		config := gohorizon.NewConfiguration()
		config.UserAgent = userAgent
		config.Host = host
		config.Scheme = "https"

		if !sslVerify {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			config.HTTPClient = &http.Client{Transport: tr}
		}

		client := gohorizon.NewAPIClient(config)

		body := gohorizon.NewAuthLogin(domain, password, username)
		tokens, _, err := client.AuthApi.LoginUser(ctx).Body(*body).Execute()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		client.GetConfig().AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", *tokens.AccessToken))
		return &apiClient{Client: *client}, nil
	}
}
