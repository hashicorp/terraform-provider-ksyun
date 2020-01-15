package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunListenerServersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataListenerServersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_listener_servers.foo"),
				),
			},
		},
	})
}

const testAccDataListenerServersConfig = `
data "ksyun_listener_servers" "foo" {
  output_file="output_result"
  ids=[]
  listener_id=[]
  real_server_ip=["10.72.20.126","172.31.16.20"]
}
`
