package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunRegisterBackendServer_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_register_backend_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRegisterBackendServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisterBackendServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegisterBackendServerExists("ksyun_lb_register_backend_server.foo", &val),
					testAccCheckRegisterBackendServerAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunRegisterBackendServer_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_register_backend_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRegisterBackendServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisterBackendServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegisterBackendServerExists("ksyun_lb_register_backend_server.foo", &val),
					testAccCheckRegisterBackendServerAttributes(&val),
				),
			},
			{
				Config: testAccRegisterBackendServerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRegisterBackendServerExists("ksyun_lb_register_backend_server.foo", &val),
					testAccCheckRegisterBackendServerAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckRegisterBackendServerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("RegisterBackendServer id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		regist := make(map[string]interface{})
		regist["RegisterId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeBackendServers(&regist)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckRegisterBackendServerAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["BackendServerSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("RegisterBackendServer id is empty")
			}
		}
		return nil
	}
}
func testAccCheckRegisterBackendServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_register_backend_server" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		regist := make(map[string]interface{})
		regist["RegisterId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeBackendServers(&regist)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["BackendServerSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("RegisterBackendServer still exist")
			}
		}
	}

	return nil
}

const testAccRegisterBackendServerConfig = `
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
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"
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
  sriov_net_support="false"
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}

resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name="xuan-tf"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
}

resource "ksyun_lb_register_backend_server" "foo" {
  backend_server_group_id="${ksyun_lb_backend_server_group.default.id}"
  backend_server_ip="${ksyun_instance.default.private_ip_address}"
  backend_server_port="8081"
  weight=10
}
`

const testAccRegisterBackendServerUpdateConfig = `
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
resource "ksyun_instance" "default" {
  image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type="N3.2B"
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
  sriov_net_support="false"
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}

resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name="xuan-tf"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
}

resource "ksyun_lb_register_backend_server" "foo" {
  backend_server_group_id="${ksyun_lb_backend_server_group.default.id}"
  backend_server_ip="${ksyun_instance.default.private_ip_address}"
  backend_server_port="8081"
  weight=20
}`
