package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunLbAclsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLbAclsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_lb_acls.foo"),
				),
			},
		},
	})
}

const testAccDataLbAclsConfig = `
data "ksyun_lb_acls" "foo" {
  output_file="output_result"
  ids=[]
}`
