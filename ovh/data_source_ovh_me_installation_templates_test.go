package ovh

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccMeInstallationTemplatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMeInstallationTemplatesDatasourceConfig_Basic,
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

const testAccMeInstallationTemplatesDatasourceConfig_Basic = `
data "ovh_me_installation_templates" "templates" {}

output result {
   value = tostring(length(data.ovh_me_installation_templates.templates.result) == 0)
}
`
