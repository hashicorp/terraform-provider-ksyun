package ksyun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKsyunVolume_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("ksyun_volume.foo"),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("volume id is empty")
		}
		return nil
	}
}

func testAccCheckVolumeDestory(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_volume" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		volume := make(map[string]interface{})
		volume["VolumeId.1"] = rs.Primary.ID
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
				return fmt.Errorf("volume still exist")
			}
		}
	}
	return nil
}

const testAccVolumeConfig = `
resource "ksyun_volume" "default" {
  volume_name="test"
  volume_type="SSD3.0"
  size=15
  charge_type="Daily"
  availability_zone="cn-shanghai-3a"
  volume_desc="test"
}
`
