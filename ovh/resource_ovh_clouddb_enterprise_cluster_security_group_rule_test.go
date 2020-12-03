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
	testCloudDBEnterpriseClusterSecurityGroupRuleConfig = `
data "ovh_clouddb_enterprise_cluster" "db" {
	cluster_id = "%s"
}
	
resource "ovh_clouddb_enterprise_cluster_security_group" "sg" {
  cluster_id = data.ovh_clouddb_enterprise_cluster.db.id
  name = "%s"
}
	
resource "ovh_clouddb_enterprise_cluster_security_group_rule" "rule" {
  cluster_id = data.ovh_clouddb_enterprise_cluster.db.id
  security_group_id = ovh_clouddb_enterprise_cluster_security_group.sg.id
  source = "%s"
}
`
	testCloudDBEnterpriseClusterSecurityGroupRuleSource1 = "51.51.51.0/24"
	testCloudDBEnterpriseClusterSecurityGroupRuleSource2 = "52.51.51.0/24"
)

func init() {
	resource.AddTestSweepers("ovh_clouddb_enterprise_cluster_security_group_rule", &resource.Sweeper{
		Name: "ovh_clouddb_enterprise_cluster_security_group_rule",
		F:    testSweepCloudDBEnterpriseClusterSecurityGroupRule,
	})
}

func testSweepCloudDBEnterpriseClusterSecurityGroupRule(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	clusterId := os.Getenv("OVH_CLOUDDB_ENTERPRISE")
	if clusterId == "" {
		log.Print("[DEBUG] No OVH_CLOUDDB_ENTERPRISE envvar specified. nothing to sweep")
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

		var ruleIds []string

		endpoint = fmt.Sprintf(
			CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl,
			url.PathEscape(clusterId),
			url.PathEscape(id),
		)

		if err := client.Get(endpoint, &ruleIds); err != nil {
			return fmt.Errorf("Error calling GET %s:\n\t%q", endpoint, err)
		}

		if len(ruleIds) == 0 {
			log.Printf("[DEBUG] No Security Group Rule to sweep on enterprise cloud db %s sg %s", clusterId, id)
			return nil
		}

		for _, rid := range ruleIds {
			log.Printf(
				"[INFO] Deleting Security Group Rule %v/%v on cluster %v",
				id,
				rid,
				clusterId,
			)
			endpoint = fmt.Sprintf(
				CloudDBEnterpriseClusterSecurityGroupRuleBaseUrl+"/%s",
				url.PathEscape(clusterId),
				url.PathEscape(id),
				url.PathEscape(rid),
			)
			if err := client.Delete(endpoint, nil); err != nil {
				return fmt.Errorf("Error calling DELETE %s:\n\t%q", endpoint, err)
			}
		}
	}
	return nil
}

func TestAccCloudDBEnterpriseClusterSecurityGroupRule(t *testing.T) {
	cluster := os.Getenv("OVH_CLOUDDB_ENTERPRISE_TEST")
	groupName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloudDBEnterpriseCluster(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testCloudDBEnterpriseClusterSecurityGroupRuleConfig,
					cluster,
					groupName,
					testCloudDBEnterpriseClusterSecurityGroupRuleSource1,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group.sg", "name", groupName),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group_rule.rule", "source", testCloudDBEnterpriseClusterSecurityGroupRuleSource1),
					resource.TestCheckResourceAttr(
						"data.ovh_clouddb_enterprise_cluster.db", "cluster_id", cluster),
				),
			},
			{
				Config: fmt.Sprintf(
					testCloudDBEnterpriseClusterSecurityGroupRuleConfig,
					cluster,
					groupName,
					testCloudDBEnterpriseClusterSecurityGroupRuleSource2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group.sg", "name", groupName),
					resource.TestCheckResourceAttr(
						"ovh_clouddb_enterprise_cluster_security_group_rule.rule", "source", testCloudDBEnterpriseClusterSecurityGroupRuleSource2),
				),
			},
		},
	})
}
