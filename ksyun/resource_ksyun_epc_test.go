package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunEpc_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_epc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEpcDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccEpcConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
					resource.TestCheckResourceAttr("ksyun_epc.foo", "host_name", "tf-acc-epc"),
				),
			},
		},
	})
}

func TestAccKsyunEpc_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_epc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEpcDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccEpcConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
				),
			},
			{
				Config: testAccEpcConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckEpcExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("host id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		epc := make(map[string]interface{})
		epc["HostId.1"] = rs.Primary.ID
		ptr, err := client.epcconn.DescribeEpcs(&epc)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HostSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckEpcAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["HostSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("host id is empty")
			}
		}
		return nil
	}
}

func testAccCheckEpcDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_epc" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		epc := make(map[string]interface{})
		epc["HostId.1"] = rs.Primary.ID
		ptr, err := client.epcconn.DescribeEpcs(&epc)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HostSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Host still exist")
			}
		}
	}

	return nil
}

const testAccEpcConfig = `
resource "ksyun_epc" "foo" {
	host_name        = "tf-acc-epc"
}
`

const testAccEpcConfigUpdate = `
resource "ksyun_epc" "foo" {
	host_name        = "tf-acc-epc-1"
}
`
