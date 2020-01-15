package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunLbsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLbsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lbs.foo"),
				),
			},
		},
	})
}

const testAccDataLbsConfig = `
data "ksyun_lbs" "foo" {
  output_file="output_result"
  ids=[]
  state=""
  vpc_id=[]
}
`
