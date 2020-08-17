package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunLbBackendServerGroup_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_backend_server_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbBackendServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbBackendServerGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbBackendServerGroupExists("ksyun_lb_backend_server_group.foo", &val),
					testAccCheckLbBackendServerGroupAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunLbBackendServerGroup_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_backend_server_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbBackendServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbBackendServerGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbBackendServerGroupExists("ksyun_lb_backend_server_group.foo", &val),
					testAccCheckLbBackendServerGroupAttributes(&val),
				),
			},
			{
				Config: testAccLbBackendServerGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbBackendServerGroupExists("ksyun_lb_backend_server_group.foo", &val),
					testAccCheckLbBackendServerGroupAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckLbBackendServerGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("LbBackendServerGroup id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		backendServerGroup := make(map[string]interface{})
		backendServerGroup["BackendServerGroupId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeBackendServerGroups(&backendServerGroup)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckLbBackendServerGroupAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("BackendServerGroup id is empty")
			}
		}
		return nil
	}
}
func testAccCheckLbBackendServerGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_backend_server_group" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		backendServerGroup := make(map[string]interface{})
		backendServerGroup["BackendServerGroupId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeBackendServerGroups(&backendServerGroup)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerGroupSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("BackendServerGroup still exist")
			}
		}
	}

	return nil
}

const testAccLbBackendServerGroupConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb_backend_server_group" "foo" {
  backend_server_group_name="xuan-tf"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
}
`

const testAccLbBackendServerGroupUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb_backend_server_group" "foo" {
  backend_server_group_name="xuan-tf-update"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
}
`
