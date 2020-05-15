package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

func TestAccKsyunLbAclEntry_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_acl_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbAclEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbAclEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclEntryExists("ksyun_lb_acl_entry.foo", &val),
					testAccCheckLbAclEntryAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunLbAclEntry_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_acl_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbAclEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbAclEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclEntryExists("ksyun_lb_acl_entry.foo", &val),
					testAccCheckLbAclEntryAttributes(&val),
				),
			},
			{
				Config: testAccLbAclEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclEntryExists("ksyun_lb_acl_entry.foo", &val),
					testAccCheckLbAclEntryAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckLbAclEntryExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("LbAclEntry id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		lbAclEntry := make(map[string]interface{})
		ids := strings.Split(rs.Primary.ID, ":")
		if len(ids) == 0 {
			return fmt.Errorf("LbAclEntry id is error")
		}
		lbAclEntry["LoadBalancerAclId.1"] = ids[0]
		ptr, err := client.slbconn.DescribeLoadBalancerAcls(&lbAclEntry)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["LoadBalancerAclSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckLbAclEntryAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["LoadBalancerAclSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("LbAclEntry id is empty")
			}
		}
		return nil
	}
}
func testAccCheckLbAclEntryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_acl_entry" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		lbAclEntry := make(map[string]interface{})
		ids := strings.Split(rs.Primary.ID, ":")
		if len(ids) == 0 {
			return fmt.Errorf("LbAclEntry id is error")
		}
		lbAclEntry["LoadBalancerAclId.1"] = ids[0]
		ptr, err := client.slbconn.DescribeLoadBalancerAcls(&lbAclEntry)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["LoadBalancerAclSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("LbAclEntry still exist")
			}
		}
	}

	return nil
}

const testAccLbAclEntryConfig = `
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "ksyun-lb-acl-entry-tf"
}
resource "ksyun_lb_acl_entry" "foo" {
  load_balancer_acl_id = "${ksyun_lb_acl.default.id}"
  cidr_block = "192.168.11.1/32"
  rule_number = 10
  rule_action = "allow"
  protocol = "ip"
}
`

const testAccLbAclEntryUpdateConfig = `
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "ksyun-lb-acl-entry-tf"
}
resource "ksyun_lb_acl_entry" "foo" {
  load_balancer_acl_id = "${ksyun_lb_acl.default.id}"
  cidr_block = "196.168.11.1/32"
  rule_number = 9
  rule_action = "allow"
  protocol = "ip"
}
`
