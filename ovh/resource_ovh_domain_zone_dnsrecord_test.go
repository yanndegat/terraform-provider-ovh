package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccDomainZoneName = "testaccdomainzone"

var testAccCheckOVHRecordConfig_basic = fmt.Sprintf(`
resource "ovh_domain_zone_dnsrecord" "foobar" {
	zone_name = "%s"
	sub_domain = "terraform"
	target = "192.168.0.10"
	field_type = "A"
	ttl = "3600"
}`, testAccDomainZoneName)

var testAccCheckOVHRecordConfig_new_value_1 = fmt.Sprintf(`
resource "ovh_domain_zone_dnsrecord" "foobar" {
	zone_name = "%s"
	sub_domain = "terraform"
	target = "192.168.0.11"
	field_type = "A"
	ttl = "3600"
}`, testAccDomainZoneName)

var testAccCheckOVHRecordConfig_new_value_2 = fmt.Sprintf(`
resource "ovh_domain_zone_dnsrecord" "foobar" {
	zone_name = "%s"
	sub_domain = "terraform2"
	target = "192.168.0.11"
	field_type = "A"
	ttl = "3600"
}`, testAccDomainZoneName)

var testAccCheckOVHRecordConfig_new_value_3 = fmt.Sprintf(`
resource "ovh_domain_zone_dnsrecord" "foobar" {
	zone = "%s"
	sub_domain = "terraform3"
	target = "192.168.0.13"
	field_type = "A"
	ttl = "3604"
}`, testAccDomainZoneName)

func TestAccDomainZoneDnsRecord_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOVHRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_dnsrecord.foobar"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "sub_domain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "zone", testAccDomainZoneName),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "target", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "ttl", "3600"),
				),
			},
		},
	})
}

func TestAccOVHRecord_Updated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOVHRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_dnsrecord.foobar"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "sub_domain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "zone", testAccDomainZoneName),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "target", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_dnsrecord.foobar"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "sub_domain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "zone", testAccDomainZoneName),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_dnsrecord.foobar"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "sub_domain", "terraform2"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "zone", testAccDomainZoneName),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_dnsrecord.foobar"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "sub_domain", "terraform3"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "zone", testAccDomainZoneName),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "target", "192.168.0.13"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_dnsrecord.foobar", "ttl", "3604"),
				),
			},
		},
	})
}

func testAccCheckOVHRecordDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_zone_dnsrecord" {
			continue
		}

		exists, err := domainZoneDnsRecordExists(testAccDomainZoneName, rs.Primary.ID, config)

		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckOVHRecordExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		exists, err := domainZoneDnsRecordExists(testAccDomainZoneName, rs.Primary.ID, config)

		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("Record doesn't exist")
		}

		return nil
	}
}
