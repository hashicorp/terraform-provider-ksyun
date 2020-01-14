package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunVPCsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVPCsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_vpcs.foo"),
				),
			},
		},
	})
}

const testAccDataVPCsConfig = `
resource "ksyun_vpc" "default" {
	vpc_name        = "tf-acc-vpc-data"
    cidr_block      = "192.168.0.0/16"
}
data "ksyun_vpcs" "foo" {
    ids = ["${ksyun_vpc.default.id}"]
	output_file = "output_result"
}
`
