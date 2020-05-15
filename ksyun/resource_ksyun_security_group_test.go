package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunSecurityGroup_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("ksyun_security_group.foo", &val),
					testAccCheckSecurityGroupAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunSecurityGroup_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("ksyun_security_group.foo", &val),
					testAccCheckSecurityGroupAttributes(&val),
				),
			},
			{
				Config: testAccSecurityGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("ksyun_security_group.foo", &val),
					testAccCheckSecurityGroupAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SecurityGroup id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		securityGroup := make(map[string]interface{})
		securityGroup["SecurityGroupId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeSecurityGroups(&securityGroup)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["SecurityGroupSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckSecurityGroupAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["SecurityGroupSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("SecurityGroup id is empty")
			}
		}
		return nil
	}
}
func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_security_group" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		securityGroup := make(map[string]interface{})
		securityGroup["SecurityGroupId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeSecurityGroups(&securityGroup)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["SecurityGroupSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("SecurityGroup still exist")
			}
		}
	}

	return nil
}

const testAccSecurityGroupConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "foo" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}`

const testAccSecurityGroupUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "foo" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group-update"
}
`
