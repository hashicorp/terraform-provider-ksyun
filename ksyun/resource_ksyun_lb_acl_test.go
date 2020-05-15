package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunLbAcl_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_acl.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbAclConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclExists("ksyun_lb_acl.foo", &val),
					testAccCheckLbAclAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunLbAcl_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_acl.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbAclConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclExists("ksyun_lb_acl.foo", &val),
					testAccCheckLbAclAttributes(&val),
				),
			},
			{
				Config: testAccLbAclUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbAclExists("ksyun_lb_acl.foo", &val),
					testAccCheckLbAclAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckLbAclExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("LbAcl id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		lbAcl := make(map[string]interface{})
		lbAcl["LoadBalancerAclId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeLoadBalancerAcls(&lbAcl)
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
func testAccCheckLbAclAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["LoadBalancerAclSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("LbAcl id is empty")
			}
		}
		return nil
	}
}
func testAccCheckLbAclDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_acl" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		lbAcl := make(map[string]interface{})
		lbAcl["LoadBalancerAclId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeLoadBalancerAcls(&lbAcl)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["LoadBalancerAclSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("LbAcl still exist")
			}
		}
	}

	return nil
}

const testAccLbAclConfig = `
resource "ksyun_lb_acl" "foo" {
  load_balancer_acl_name = "ksyun-lb-acl-tf"
}`

const testAccLbAclUpdateConfig = `
resource "ksyun_lb_acl" "foo" {
  load_balancer_acl_name = "ksyun-lb-acl-tf-update"
}
`
