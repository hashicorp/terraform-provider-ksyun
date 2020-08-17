package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunListenerHostHeadersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataListenerHostHeadersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_listener_host_headers.foo"),
				),
			},
		},
	})
}

const testAccDataListenerHostHeadersConfig = `
data "ksyun_listener_host_headers" "foo" {
  output_file="output_result"
  ids=[]
  listener_id=[]
}
`
