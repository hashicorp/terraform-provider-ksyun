package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunEpcsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEpcsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_epcs.foo"),
				),
			},
		},
	})
}

const testAccDataEpcsConfig = `

data "ksyun_epcs" "foo" {
    ids = ["4c3656de-9907-4c0d-868e-969cee396f37"]
	output_file = "output_result_foo"
}
`
