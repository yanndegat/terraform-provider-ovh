package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testCloudDBEnterpriseClusterSecurityGroupConfig = `
data "ovh_clouddb_enterprise_cluster" "db" {
	cluster_id = "%s"
}
	
resource "ovh_clouddb_enterprise_cluster_security_group" "sg" {
  cluster_id = data.ovh_clouddb_enterprise_cluster.db.id
  name = "%s"
}
`
)

func init() {
	resource.AddTestSweepers("ovh_clouddb_enterprise_cluster_security_group", &resource.Sweeper{
		Name: "ovh_clouddb_enterprise_cluster_security_group",
		Dependencies: []string{
			"ovh_clouddb_enterprise_cluster_security_group_rule",
		},
		F: testSweepCloudDBEnterpriseClusterSecurityGroup,
	})
}

func testSweepCloudDBEnterpriseClusterSecurityGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	clusterId := os.Getenv("OVH_CLOUDDB_ENTERPRISE_TEST")
	if clusterId == "" {
		log.Print("[DEBUG] No OVH_CLOUDDB_ENTERPRISE_TEST envvar specified. nothing to sweep")
		return nil
	}

	var Ids []string
	endpoint := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl, url.PathEscape(clusterId))

	if err := client.Get(endpoint, &Ids); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
	}

	if len(Ids) == 0 {
		log.Printf("[DEBUG] No SG to sweep on enterprise cloud db %s", clusterId)
		return nil
	}

	for _, id := range Ids {
		urlGet := fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", clusterId, id)
		var sgResp CloudDBEnterpriseClusterSecurityGroup
		err := client.Get(urlGet, &sgResp)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(sgResp.Name, test_prefix) {
			continue
		}

		log.Printf(
			"[INFO] Deleting Security Group %s %v on cluster %v",
			id,
			sgResp.Name,
			clusterId,
		)
		endpoint = fmt.Sprintf(CloudDBEnterpriseClusterSecurityGroupBaseUrl+"/%s", url.PathEscape(clusterId), id)
		if err := client.Delete(endpoint, nil); err != nil {
			return fmt.Errorf("Error calling DELETE %s:\n\t%q", endpoint, err)
		}
	}
	return nil
}

func TestAccCloudDBEnterpriseClusterSecurityGroup(t *testing.T) {
	cluster := os.Getenv("OVH_CLOUDDB_ENTERPRISE_TEST")
	groupName1 := acctest.RandomWithPrefix(test_prefix)
	groupName2 := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDBEnterpriseCluster(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testCloudDBEnterpriseClusterSecurityGroupConfig, cluster, groupName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group.sg", "name", groupName1),
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise_cluster.db", "cluster_id", cluster),
				),
			},
			{
				Config: fmt.Sprintf(testCloudDBEnterpriseClusterSecurityGroupConfig, cluster, groupName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group.sg", "name", groupName2),
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise_cluster.db", "cluster_id", cluster),
				),
			},
		},
	})
}
