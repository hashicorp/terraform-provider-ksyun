package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKcsSecGroupAllocate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKcsSecurityGroupAllocateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKcsSecGroupAllocateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKcsSecurityGroupAllocateExists("ksyun_redis_sec_group_allocate.result"),
				),
			},
		},
	})
}

func testAccCheckKcsSecurityGroupAllocateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("kcs security group allocate is create failure")
		}
		return nil
	}
}

func testAccCheckKcsSecurityGroupAllocateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_redis_sec_group_allocate" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["SecurityGroupId"] = rs.Primary.ID
			ptr, err := client.kcsv1conn.DescribeSecurityGroup(&instanceCheck)
			// Verify the error is what we want
			if err != nil {
				return err
			}
			if ptr != nil {
				if (*ptr)["Data"] != nil {
					return errors.New("deallocate redis sec group failure")
				}
			}
		}
	}

	return nil
}

const testAccKcsSecGroupAllocateConfig = `
variable "available_zone" {
  default = "cn-beijing-6a"
}

resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "foo" {
  subnet_name      	= "ksyun-subnet-tf"
  cidr_block 		= "10.7.0.0/21"
  subnet_type 		= "Reserve"
  dhcp_ip_from 		= "10.7.0.2"
  dhcp_ip_to 		= "10.7.0.253"
  vpc_id 			= "${ksyun_vpc.default.id}"
  gateway_ip 		= "10.7.0.1"
  dns1 				= "198.18.254.41"
  dns2 				= "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_redis_sec_group" "allocate" {
  available_zone 	= "${var.available_zone}"
  name 				= "testTerraform000"
  description 		= "testTerraform000"
}

resource "ksyun_redis_instance" "default" {
  available_zone 		= "${var.available_zone}"
  name 					= "MyRedisInstance1101"
  mode 					= 2
  capacity 				= 1
  net_type 				= 2
  security_group_id 	= "${ksyun_redis_sec_group.allocate.id}"
  vnet_id 				= "${ksyun_subnet.default.id}"
  vpc_id 				= "${ksyun_vpc.default.id}"
  bill_type 			= 5
  duration 				= ""
  duration_unit 		= ""
  pass_word 			= "Shiwo1101"
  iam_project_id 		= "0"
  slave_num 			= 0
  protocol 				= "4.0"
  reset_all_parameters 	= false
  parameters = {
    "appendonly" 				= "no",
    "appendfsync" 				= "everysec",
    "maxmemory-policy" 			= "volatile-lru",
    "hash-max-ziplist-entries" 	= "512",
    "zset-max-ziplist-entries" 	= "128",
    "list-max-ziplist-size" 	= "-2",
    "hash-max-ziplist-value" 	= "64",
    "notify-keyspace-events" 	= "",
    "zset-max-ziplist-value" 	= "64",
    "maxmemory-samples" 		= "5",
    "set-max-intset-entries" 	= "512",
    "timeout" 					= "600",
  }
}

resource "ksyun_redis_sec_group" "add_allocate" {
  available_zone 	= "${var.available_zone}"
  name 				= "testTerraform111"
  description 		= "testTerraform111"
}

resource "ksyun_redis_sec_group_allocate" "result" {
  available_zone 		= "${var.available_zone}"
  security_group_id 	= "${ksyun_redis_sec_group.add_allocate.id}"
  cache_ids 			= ["${ksyun_redis_instance.default.id}"]
}
`
