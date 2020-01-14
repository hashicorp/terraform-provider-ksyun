package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunNetworkInterfacesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworkInterfacesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_network_interfaces.foo"),
				),
			},
		},
	})
}

const testAccDataNetworkInterfacesConfig = `
data "ksyun_network_interfaces" "foo" {
  output_file="output_result"
  ids=[]
  vpc_id=[]
  subnet_id=[]
  securitygroup_id=[]
  instance_type=[]
  instance_id=[]
  private_ip_address=[]
}
`
