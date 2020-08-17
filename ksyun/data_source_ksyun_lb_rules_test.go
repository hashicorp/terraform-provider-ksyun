package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunLbRulesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLbRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lb_rules.foo"),
				),
			},
		},
	})
}

const testAccDataLbRulesConfig = `
data "ksyun_lb_rules" "foo" {
  output_file="output_result"
  ids=[]
  host_header_id=[]
}
`
