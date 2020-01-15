package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunListenerServer_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_listener_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckListenerServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerServerExists("ksyun_lb_listener_server.foo", &val),
					testAccCheckListenerServerAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunListenerServer_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_listener_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckListenerServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerServerExists("ksyun_lb_listener_server.foo", &val),
					testAccCheckListenerServerAttributes(&val),
				),
			},
			{
				Config: testAccListenerServerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerServerExists("ksyun_lb_listener_server.foo", &val),
					testAccCheckListenerServerAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckListenerServerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ListenerServer id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		listenerServer := make(map[string]interface{})
		listenerServer["RegisterId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeInstancesWithListener(&listenerServer)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RealServerSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckListenerServerAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["RealServerSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("ListenerServer id is empty")
			}
		}
		return nil
	}
}
func testAccCheckListenerServerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_listener_server" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		listenerServer := make(map[string]interface{})
		listenerServer["RegisterId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeInstancesWithListener(&listenerServer)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RealServerSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("ListenerServer still exist")
			}
		}
	}

	return nil
}

const testAccListenerServerConfig = `
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
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "default" {
  listener_name = "tf-xun-update"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
resource "ksyun_lb_listener_server" "foo" {
  listener_id = "${ksyun_lb_listener.default.id}"
  real_server_ip = "10.0.77.20"
  real_server_port = 8000
  real_server_type = "host"
  instance_id = "${ksyun_instance.default.id}"
  weight = 10
}
`

const testAccListenerServerUpdateConfig = `
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
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  force_delete=true
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "default" {
  listener_name = "tf-xun-update"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
resource "ksyun_lb_listener_server" "foo" {
  listener_id = "${ksyun_lb_listener.default.id}"
  real_server_ip = "10.0.77.20"
  real_server_port = 80
  real_server_type = "host"
  instance_id = "${ksyun_instance.default.id}"
  weight = 2
}
`
