package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"log"
	"testing"
)

func TestAccKsyunSecurityGroupEntry_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupEntryExists("ksyun_security_group_entry.foo", &val),
					testAccCheckSecurityGroupEntryAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunSecurityGroupEntry_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_security_group_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupEntryExists("ksyun_security_group_entry.foo", &val),
					testAccCheckSecurityGroupEntryAttributes(&val),
				),
			},
			{
				Config: testAccSecurityGroupEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupEntryExists("ksyun_security_group_entry.foo", &val),
					testAccCheckSecurityGroupEntryAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupEntryExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SecurityGroupEntry id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		securityGroupEntry := make(map[string]interface{})
		securityGroupEntry["SecurityGroupId.1"] = rs.Primary.Attributes["security_group_id"]
		securityGroupEntryId := rs.Primary.ID
		log.Printf("SecurityGroupId:%v", rs.Primary.Attributes["security_group_id"])
		ptr1, err := client.vpcconn.DescribeSecurityGroups(&securityGroupEntry)
		if err != nil {
			return err
		}
		ptr, ok := (*ptr1)["SecurityGroupSet"].([]interface{})
		if !ok {
			return errors.New("no SecurityGroupSet get")
		}
		if len(ptr) == 0 {
			return errors.New("no SecurityGroupSet get")
		}
		securityGroup, ok := ptr[0].(map[string]interface{})
		l := securityGroup["SecurityGroupEntrySet"].([]interface{})
		if len(l) == 0 {
			return errors.New("no SecurityGroupEntrySet get")
		}
		for _, securityGroupEntry := range l {
			securityGroupEntryItem, ok := securityGroupEntry.(map[string]interface{})
			if !ok {
				return errors.New("no securityGroupEntry get")
			}
			if securityGroupEntryItem["SecurityGroupEntryId"] == securityGroupEntryId {
				*val = securityGroupEntryItem
				return nil
			}
		}
		return errors.New("securityGroupEntry not exist")
	}
}
func testAccCheckSecurityGroupEntryAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//	if val != nil {
		//		l := (*val)["SecurityGroupEntrySet"].([]interface{})
		if len(*val) == 0 {
			return fmt.Errorf("SecurityGroupEntry id is empty")
		}
		//	}
		return nil
	}
}
func testAccCheckSecurityGroupEntryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_security_group_entry" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		securityGroupEntry := make(map[string]interface{})
		securityGroupEntry["SecurityGroupId.1"] = rs.Primary.Attributes["security_group_id"]
		securityGroupEntryId := rs.Primary.ID
		ptr1, err := client.vpcconn.DescribeSecurityGroups(&securityGroupEntry)
		// Verify the error is what we want
		if err != nil {
			return err
		}
		ptr, ok := (*ptr1)["SecurityGroupSet"].([]interface{})
		if !ok {
			return nil
		}
		if len(ptr) == 0 {
			return nil
		}
		securityGroup, ok := ptr[0].(map[string]interface{})
		l := securityGroup["SecurityGroupEntrySet"].([]interface{})
		if len(l) == 0 {
			return nil
		}
		for _, securityGroupEntry := range l {
			securityGroupEntryItem, ok := securityGroupEntry.(map[string]interface{})
			if !ok {
				return errors.New("no securityGroupEntry get")
			}
			if securityGroupEntryItem["SecurityGroupEntryId"] == securityGroupEntryId {
				return errors.New("securityGroupEntry still exist")
			}
		}
	}

	return nil
}

const testAccSecurityGroupEntryConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_security_group_entry" "foo" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.1.1/32"
  direction="in"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}`

const testAccSecurityGroupEntryUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_security_group_entry" "foo" {
  description = "test1"
  security_group_id="${ksyun_security_group.default.id}"
  cidr_block="10.0.3.1/32"
  direction="out"
  protocol="ip"
  icmp_type=0
  icmp_code=0
  port_range_from=0
  port_range_to=0
}
`
