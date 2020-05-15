package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunSubnetsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSubnetsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_subnets.foo"),
				),
			},
		},
	})
}

const testAccDataSubnetsConfig = `
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
data "ksyun_subnets" "foo" {
    subnet_types = ["Normal"]
    vpc_ids = ["${ksyun_subnet.default.id}"]
	output_file = "output_result"
}
`
