package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunHealthCheck_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_healthcheck.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthCheckConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHealthCheckExists("ksyun_healthcheck.foo", &val),
					testAccCheckHealthCheckAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunHealthCheck_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_healthcheck.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckHealthCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHealthCheckConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHealthCheckExists("ksyun_healthcheck.foo", &val),
					testAccCheckHealthCheckAttributes(&val),
				),
			},
			{
				Config: testAccHealthCheckUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHealthCheckExists("ksyun_healthcheck.foo", &val),
					testAccCheckHealthCheckAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckHealthCheckExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HealthCheck id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		healthCheck := make(map[string]interface{})
		healthCheck["HealthCheckId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeHealthChecks(&healthCheck)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HealthCheckSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckHealthCheckAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["HealthCheckSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("HealthCheck id is empty")
			}
		}
		return nil
	}
}
func testAccCheckHealthCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_healthcheck" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		healthCheck := make(map[string]interface{})
		healthCheck["HealthCheckId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeHealthChecks(&healthCheck)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HealthCheckSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("HealthCheck still exist")
			}
		}
	}

	return nil
}

const testAccHealthCheckConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
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
  listener_name = "tf-xun-2"
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
resource "ksyun_healthcheck" "foo" {
  listener_id = "${ksyun_lb_listener.default.id}"
  health_check_state = "stop"
  healthy_threshold = 2
  interval = 20
  timeout = 200
  unhealthy_threshold = 2
  url_path = "/monitor"
  is_default_host_name = true
  host_name = "www.ksyun.com"
}

`

const testAccHealthCheckUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
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
  listener_name = "tf-xun-2"
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
resource "ksyun_healthcheck" "foo" {
  listener_id = "${ksyun_lb_listener.default.id}"
  health_check_state = "stop"
  healthy_threshold = 10
  interval = 3
  timeout = 300
  unhealthy_threshold = 3
  url_path = "/monitor"
  is_default_host_name = true
  host_name = "www.ksyun.com"
}

`
