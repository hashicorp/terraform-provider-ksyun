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
			volume, ok2 := l[0].(map[string]interface{})
			if !ok2 {
				continue
			}
			status, ok3 := volume["VolumeStatus"]
			if !ok3 {
				continue
			}
			if status == "recycling" || status == "deleting" {
				continue
			}
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
resource "ksyun_volume" "foo" {
  volume_name="ksyun_volume_tf_test"
  volume_type="SSD3.0"
  size=10
  charge_type="Daily"
  availability_zone="cn-beijing-6a"
  volume_desc="ksyun_volume_tf_test"
}
`
