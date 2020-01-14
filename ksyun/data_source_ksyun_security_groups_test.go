package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunSecurityGroupsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_security_groups.foo"),
				),
			},
		},
	})
}

const testAccDataSecurityGroupsConfig = `
data "ksyun_security_groups" "foo" {
  output_file="output_result"
  ids=[]
  vpc_id=[]
}
`
