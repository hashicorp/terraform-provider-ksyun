package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func TestAccKsyunMongodbSecurityRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbSecurityRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbSecurityRuleExists("ksyun_mongodb_security_rule.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbSecurityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance is not exist")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		securityRuleCheck := make(map[string]interface{})
		securityRuleCheck["InstanceId"] = rs.Primary.ID
		resp, err := client.mongodbconn.ListSecurityGroupRules(&securityRuleCheck)
		if err != nil {
			return fmt.Errorf("error on reading mongodb instance security rule %q, %s", rs.Primary.ID, err)
		}
		rules := (*resp)["MongoDBSecurityGroupRule"].([]interface{})
		if len(rules) == 0 {
			return fmt.Errorf("mongodb instance security rule is not exist")
		}

		return nil
	}
}

func testAccCheckMongodbSecurityRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_mongodb_security_rule" {
			securityRuleCheck := make(map[string]interface{})
			securityRuleCheck["InstanceId"] = rs.Primary.ID
			resp, err := client.mongodbconn.ListSecurityGroupRules(&securityRuleCheck)

			if err != nil {
				if strings.Contains(err.Error(), "InstanceNotFound") {
					return nil
				}
				return fmt.Errorf("error on reading mongodb instance security rule %q, %s", rs.Primary.ID, err)
			}

			rules := (*resp)["MongoDBSecurityGroupRule"].([]interface{})
			if len(rules) > 0 {
				return fmt.Errorf("delete mongodb security rule failure")
			}

			return nil
		}
	}

	return nil
}

const testAccMongodbSecurityRuleConfig = `
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
resource "ksyun_vpc" "default" {
  vpc_name = "ksyun-vpc-tf"
  cidr_block = "10.1.0.0/23"
}
resource "ksyun_subnet" "default" {
  subnet_name = "ksyun-subnet-tf"
  cidr_block = "10.1.0.0/23"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
resource "ksyun_mongodb_instance" "default" {
  name = "mongodb_repset"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 5
  node_num = 3
  vpc_id = "${ksyun_vpc.default.id}"
  vnet_id = "${ksyun_subnet.default.id}"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "0"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
resource "ksyun_mongodb_security_rule" "default" {
  instance_id = "${ksyun_mongodb_instance.default.id}"
  cidrs = "192.168.10.1/32"
}
`
