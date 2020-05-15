package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunListenersDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataListenersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_listeners.foo"),
				),
			},
		},
	})
}

const testAccDataListenersConfig = `
data "ksyun_listeners" "foo" {
  output_file="output_result"
  ids=[""]
  load_balancer_id=[]

}
`
