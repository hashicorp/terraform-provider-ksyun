package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunSqlServer_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_sqlserver.ks-ss-233",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSqlServerDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccSqlServerConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckSqlServerExists("ksyun_sqlserver.ks-ss-233", &val),
				),
			},
		},
	})
}

func testCheckSqlServerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found : %s", n)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("instance is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		resp, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			return err
		}
		if resp != nil {
			bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
			if !dataOk {
				return fmt.Errorf("error on reading Instance(krds)  %+v", (*resp)["Error"])
			}
			instances := bodyData["Instances"].([]interface{})
			if len(instances) == 0 {
				return fmt.Errorf("no instance find, instance number is 0")
			}
		}
		*val = *resp
		return nil
	}
}

func testAccCheckSqlServerDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_krds" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		_, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			if err.(awserr.Error).Code() == "NOT_FOUND" {
				return nil
			}
			return err
		}

	}

	return nil
}

const testAccSqlServerConfig = `

variable "available_zone" {
  default = "cn-shanghai-2a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_sqlserver" "ks-ss-233"{
 output_file = "output_file"
 db_instance_class= "db.ram.2|db.disk.100"
 db_instance_name = "ksyun_sqlserver_1"
 db_instance_type = "HRDS_SS"
 engine = "SQLServer"
 engine_version = "2008r2"
 master_user_name = "admin"
 master_user_password = "123qweASD"
 vpc_id = "${ksyun_vpc.default.id}"
 subnet_id = "${ksyun_subnet.foo.id}"
 bill_type = "DAY"

}
`
