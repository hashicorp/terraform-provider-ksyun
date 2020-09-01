package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKcsSecGroupRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKcsSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKcsSecGroupRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsSecurityGroupRuleExists("ksyun_redis_sec_group_rule.rule"),
				),
			},
		},
	})
}

func testAccCheckKcsSecurityGroupRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs security group rule is create failure")
		}
		return nil
	}
}

func testAccCheckKcsSecurityGroupRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_sec_group_rule" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["SecurityGroupId"] = rs.Primary.ID
			ptr, err := client.kcsv1conn.DescribeSecurityGroup(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				return err
			}
			if ptr != nil {
				if (*ptr)["Data"] != nil {
					return errors.New("delete redis sec group failure")
				}
			}
		}
	}

	return nil
}

const testAccKcsSecGroupRuleConfig = `
resource "ksyun_redis_sec_group" "rule" {
  name 			= "testTerraform888"
  description 	= "testTerraform888"
}

resource "ksyun_redis_sec_group_rule" "rule" {
  security_group_id = "${ksyun_redis_sec_group.rule.id}"
  rules 			= ["172.16.0.0/32","192.168.0.0/32"]
}
`
