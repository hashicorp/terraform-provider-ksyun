package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKsyunVolumesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataVolumesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_volumes.foo"),
				),
			},
		},
	})
}

const testAccDataVolumesConfig = `
data "ksyun_volumes" "foo" {
  output_file="output_result"
  ids = []
  volume_category=""
  volume_status=""
  volume_type=""
  availability_zone=""
}
`
