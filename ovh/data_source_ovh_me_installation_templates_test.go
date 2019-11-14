package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMeInstallationTemplatesDataSource_basic(t *testing.T) {
	templateName := acctest.RandomWithPrefix(test_prefix)
	presetup := fmt.Sprintf(
		testAccMeInstallationTemplatesDatasourceConfig_presetup,
		templateName,
	)
	config := fmt.Sprintf(
		testAccMeInstallationTemplatesDatasourceConfig_Basic,
		templateName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: presetup,
				Check: resource.TestCheckResourceAttr(
					"ovh_me_installation_template.template",
					"template_name",
					templateName,
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_me_installation_templates.templates",
						"result.#",
					),
					// resource.TestCheckResourceAttr(
					// 	fmt.Sprintf(
					// 		"data.ovh_me_installation_template.template[\\\"%s\\\"]",
					// 		templateName,
					// 	),
					// 	"template_name",
					// 	templateName,
					// ),
				),
			},
		},
	})
}

const testAccMeInstallationTemplatesDatasourceConfig_presetup = `
terraform {
  required_version = ">= 0.12"
}

resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}
`

const testAccMeInstallationTemplatesDatasourceConfig_Basic = `
terraform {
  required_version = ">= 0.12"
}

resource "ovh_me_installation_template" "template" {
  base_template_name = "centos7_64"
  template_name      = "%s"
  default_language   = "en"
}

data "ovh_me_installation_templates" "templates" {}

data "ovh_me_installation_template" "template" {
  for_each           = toset(data.ovh_me_installation_templates.templates.result)
  template_name      = each.value
}
`
