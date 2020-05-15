package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunLbListenerServersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLbListenerServersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lb_listener_servers.foo"),
				),
			},
		},
	})
}

const testAccDataLbListenerServersConfig = `
data "ksyun_lb_listener_servers" "foo" {
  output_file="output_result"
  ids=[]
  listener_id=[]
  real_server_ip=[]
}
`
