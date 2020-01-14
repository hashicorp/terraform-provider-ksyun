package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunVPC_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVPCDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVPCConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists("ksyun_vpc.foo", &val),
					testAccCheckVPCAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_vpc.foo", "vpc_name", "tf-acc-vpc"),
					resource.TestCheckResourceAttr("ksyun_vpc.foo", "cidr_block", "192.168.0.0/16"),
				),
			},
		},
	})
}

func TestAccKsyunVPC_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVPCDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccVPCConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists("ksyun_vpc.foo", &val),
					testAccCheckVPCAttributes(&val),
					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "vpc_name", "tf-acc-vpc"),
					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "cidr_block", "192.168.0.0/16"),
				),
			},
			{
				Config: testAccVPCConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCExists("ksyun_vpc.foo", &val),
					testAccCheckVPCAttributes(&val),
					//resource.TestCheckResourceAttr("ksyun_vpc.foo", "vpc_name", "tf-acc-vpc-1"),
				),
			},
		},
	})
}

func testAccCheckVPCExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		vpc := make(map[string]interface{})
		vpc["VpcId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeVpcs(&vpc)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["VpcSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckVPCAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["VpcSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("vpc id is empty")
			}
		}
		return nil
	}
}

func testAccCheckVPCDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_vpc" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		vpc := make(map[string]interface{})
		vpc["VpcId.1"] = rs.Primary.ID
		ptr, err := client.vpcconn.DescribeVpcs(&vpc)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["VpcSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("VPC still exist")
			}
		}
	}

	return nil
}

const testAccVPCConfig = `
resource "ksyun_vpc" "foo" {
	vpc_name        = "tf-acc-vpc"
	cidr_block = "192.168.0.0/16"
}
`

const testAccVPCConfigUpdate = `
resource "ksyun_vpc" "foo" {
	vpc_name        = "tf-acc-vpc-1"
    cidr_block      = "192.168.0.0/16"
}
`
