package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFarm(t *testing.T) {
	t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDesktopPool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"horizon_desktop_pool.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceDesktopPool = `
resource "scaffolding_resource" "foo" {
  sample_attribute = "bar"
}
`
