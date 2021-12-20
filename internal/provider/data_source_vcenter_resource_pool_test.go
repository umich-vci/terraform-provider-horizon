package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcevCenterResourcePool(t *testing.T) {
	t.Skip("data source not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcevCenterResourcePool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.scaffolding_data_source.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccDataSourcevCenterResourcePool = `
data "scaffolding_data_source" "foo" {
  sample_attribute = "bar"
}
`
