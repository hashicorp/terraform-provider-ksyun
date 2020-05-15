package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunImagesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataImagesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_images.foo"),
				),
			},
		},
	})
}

const testAccDataImagesConfig = `
data "ksyun_images" "foo" {
  output_file="output_result"
  ids=[]
  name_regex="centos-7.0-20180927115228"
  is_public=true
  image_source="system"
}`
