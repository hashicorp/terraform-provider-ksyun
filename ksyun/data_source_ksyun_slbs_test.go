package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSlbsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSlbsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_slbs.foo"),
				),
			},
		},
	})
}

const testAccDataSlbsConfig = `
data "ksyun_slbs" "foo" {
  output_file="output_result"
  ids=[]
  state=""
  vpc_id=[]
}
`
