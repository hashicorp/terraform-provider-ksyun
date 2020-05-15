package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"testing"
)

func TestAccKsyunMongodbShardInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbShardInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbShardInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbShardInstanceExists("ksyun_mongodb_shard_instance.default"),
				),
			},
		},
	})
}

func testAccCheckMongodbShardInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource or data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance create failure")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		readReq := make(map[string]interface{})
		readReq["InstanceId"] = rs.Primary.ID

		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstance", readReq)
		_, err := client.mongodbconn.DescribeMongoDBInstance(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance %q, %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckMongodbShardInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*KsyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ksyun_mongodb_shard_instance" {
			instanceCheck := make(map[string]interface{})
			instanceCheck["InstanceId"] = rs.Primary.ID
			_, err := client.mongodbconn.DescribeMongoDBInstance(&instanceCheck)

			if err != nil {
				if strings.Contains(err.Error(), "InstanceNotFound") {
					return nil
				} else {
					return fmt.Errorf("mongodb instance delete failure")
				}
			}
		}
	}

	return nil
}

const testAccMongodbShardInstanceConfig = `
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
resource "ksyun_vpc" "default" {
  vpc_name = "ksyun-vpc-mongodb-tf"
  cidr_block = "10.1.0.0/23"
}
resource "ksyun_subnet" "default" {
  subnet_name = "ksyun-subnet-mongodb-tf"
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
resource "ksyun_mongodb_shard_instance" "default" {
  name = "mongodb_shard_tf"
  instance_account = "root"
  instance_password = "admin"
  mongos_class = "1C2G"
  mongos_num = 2
  shard_class = "1C2G"
  shard_num = 2
  storage = 5
  vpc_id = "${ksyun_vpc.default.id}"
  vnet_id = "${ksyun_subnet.default.id}"
  db_version = "3.6"
  pay_type = "hourlyInstantSettlement"
  iam_project_id = "0"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
`
