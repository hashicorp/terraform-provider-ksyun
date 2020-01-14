package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunInstancesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_instances.foo"),
				),
			},
		},
	})
}

const testAccDataInstancesConfig = `
data "ksyun_instances" "foo" {
  output_file = "output_result"
  ids = []
  project_id = []
  network_interface {
    network_interface_id = []
    subnet_id = []
    group_id = []
  }
  instance_state {
    name =  []
  }
  availability_zone {
    name =  []
  }
}`
