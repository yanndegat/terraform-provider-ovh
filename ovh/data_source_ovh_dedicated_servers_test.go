package ovh

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDedicatedServersDataSource_basic(t *testing.T) {
	dedicated_server := os.Getenv("OVH_DEDICATED_SERVER_SERVICE_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDedicatedServer(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServersDatasourceConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"result",
						dedicated_server,
					),
				),
			},
		},
	})
}

const testAccDedicatedServersDatasourceConfig_Basic = `
data "ovh_dedicated_servers" "servers" {}

output result {
   value = data.ovh_dedicated_servers.servers.result[0]
}
`
