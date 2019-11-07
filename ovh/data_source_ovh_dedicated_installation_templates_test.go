package ovh

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDedicatedInstallationTemplatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInstallationTemplatesDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"result",
						"true",
					),
				),
			},
		},
	})
}

const testAccDedicatedInstallationTemplatesDatasourceConfig_Basic = `
data "ovh_dedicated_installation_templates" "templates" {}

output result {
   value = tostring(length(data.ovh_dedicated_installation_templates.templates.result) > 0)
}
`
