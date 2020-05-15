package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunEipsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEipsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_eips.foo"),
				),
			},
		},
	})
}

const testAccDataEipsConfig = `

data "ksyun_eips" "foo" {
  project_id=["0"]
  instance_type=["Ipfwd"]
  network_interface_id=[]
  internet_gateway_id=[]
  band_width_share_id=[]
  line_id=[]
  public_ip=[]
  output_file = "output_result"
}
`
