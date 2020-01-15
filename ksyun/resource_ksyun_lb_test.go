package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunLb_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbExists("ksyun_lb.foo", &val),
					testAccCheckLbAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunLb_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbExists("ksyun_lb.foo", &val),
					testAccCheckLbAttributes(&val),
				),
			},
			{
				Config: testAccLbUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbExists("ksyun_lb.foo", &val),
					testAccCheckLbAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckLbExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Lb id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		lb := make(map[string]interface{})
		lb["LoadBalancerId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeLoadBalancers(&lb)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["LoadBalancerDescriptions"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckLbAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["LoadBalancerDescriptions"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Lb id is empty")
			}
		}
		return nil
	}
}
func testAccCheckLbDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		lb := make(map[string]interface{})
		lb["LoadBalancerId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeLoadBalancers(&lb)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["LoadBalancerDescriptions"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Lb still exist")
			}
		}
	}

	return nil
}

const testAccLbConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb" "foo" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "ksyun-lb-tf"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
`

const testAccLbUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb" "foo" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "ksyun-lb-tf-update"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
`
