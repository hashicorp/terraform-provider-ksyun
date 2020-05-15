package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunInstance_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("ksyun_instance.foo", &val),
					testAccCheckInstanceAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunInstance_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("ksyun_instance.foo", &val),
					testAccCheckInstanceAttributes(&val),
				),
			},
			{
				Config: testAccInstanceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("ksyun_instance.foo", &val),
					testAccCheckInstanceAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		instance := make(map[string]interface{})
		instance["InstanceId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeInstances(&instance)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["InstancesSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckInstanceAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["InstancesSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Instance id is empty")
			}
		}
		return nil
	}
}
func testAccCheckInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_instance" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		instance := make(map[string]interface{})
		instance["InstanceId.1"] = rs.Primary.ID
		ptr, err := client.kecconn.DescribeInstances(&instance)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["InstancesSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Instance still exist")
			}
		}
	}

	return nil
}

const testAccInstanceConfig = `
data "ksyun_images" "centos-7_5" {
  output_file=""
  platform= "centos-7.5"
  is_public=true
}
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_instance" "foo" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="S4.1A"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  #max_count=1
  #min_count=1
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  private_ip_address=""
  instance_name="ksyun-kec-tf"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}

`

const testAccInstanceUpdateConfig = `
data "ksyun_images" "centos-7_5" {
  output_file=""
  platform= "centos-7.5"
  is_public=true
}
data "ksyun_availability_zones" "default" {
  output_file=""
  ids=[]
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}
resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  security_group_name="ksyun-security-group"
}
resource "ksyun_instance" "foo" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="S4.1A"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  #max_count=1
  #min_count=1
  subnet_id="${ksyun_subnet.default.id}"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["${ksyun_security_group.default.id}"]
  private_ip_address=""
  instance_name="ksyun-kec-tf-update"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}

`
