package ksyun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKsyunVolumeAttach_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeAttachDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAttachConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeAttachExists("ksyun_volume_attach.foo"),
				),
			},
		},
	})
}

func testAccCheckVolumeAttachExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("attach id is empty")
		}
		return nil
	}
}

func testAccCheckVolumeAttachDestory(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_volume_attach" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		volume := make(map[string]interface{})
		volume["VolumeId.1"] = strings.Split(rs.Primary.ID, ":")[0]
		ptr, err := client.ebsconn.DescribeVolumes(&volume)
		if err != nil {
			return err
		}
		if ptr != nil {
			resp, ok := (*ptr)["Volumes"]
			if !ok {
				continue
			}
			l := resp.([]interface{})
			if len(l) == 0 {
				continue
			} else {
				if volume, ok := l[0].(map[string]interface{}); ok {
					if volume["InstanceId"] == strings.Split(rs.Primary.ID, ":")[1] {
						return fmt.Errorf("VolumeAttach still exist")
					}
				}
			}
		}
	}
	return nil
}

const testAccVolumeAttachConfig = `
resource "ksyun_volume_attach" "default" {
volume_id=""
instance_id=""
}
`
