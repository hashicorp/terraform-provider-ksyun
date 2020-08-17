package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunLbRule_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_rule.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbRuleExists("ksyun_lb_rule.foo", &val),
					testAccCheckLbRuleAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunLbRule_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_rule.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLbRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbRuleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbRuleExists("ksyun_lb_rule.foo", &val),
					testAccCheckLbRuleAttributes(&val),
				),
			},
			{
				Config: testAccLbRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLbRuleExists("ksyun_lb_rule.foo", &val),
					testAccCheckLbRuleAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckLbRuleExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("LbRule id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		lbRule := make(map[string]interface{})
		lbRule["RuleId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeRules(&lbRule)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RuleSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckLbRuleAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["RuleSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("LbRule id is empty")
			}
		}
		return nil
	}
}
func testAccCheckLbRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_rule" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		lbRule := make(map[string]interface{})
		lbRule["RuleId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeRules(&lbRule)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["RuleSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("LbRule still exist")
			}
		}
	}

	return nil
}

const testAccLbRuleConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name="xuan-tf"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
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
  listener_state = "start"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "start"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
resource "ksyun_lb_host_header" "default" {
  listener_id = "${ksyun_lb_listener.default.id}"
  host_header = "tf-xuan"
  certificate_id = ""
}

resource "ksyun_lb_rule" "foo" {
  path = "/tfxun"
  host_header_id = "${ksyun_lb_host_header.default.id}"
  backend_server_group_id="${ksyun_lb_backend_server_group.default.id}"
  listener_sync="on"
  method="RoundRobin"
  session {
    session_state = "start"
    session_persistence_period = 1000
    cookie_type = "ImplantCookie"
  #  cookie_name = "cookiexunqq"
  }
  health_check{
    health_check_state = "start"
    healthy_threshold = 2
    interval = 200
    timeout = 2000
    unhealthy_threshold = 2
    url_path = "/monitor"
    host_name = "www.ksyun.com"
  }
}
`

const testAccLbRuleUpdateConfig = `
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.5.0.0/21"
}
resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name="xuan-tf"
  vpc_id="${ksyun_vpc.default.id}"
  backend_server_group_type=""
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
  listener_state = "start"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "start"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
resource "ksyun_lb_host_header" "default" {
  listener_id = "${ksyun_lb_listener.default.id}"
  host_header = "tf-xuan"
  certificate_id = ""
}

resource "ksyun_lb_rule" "foo" {
  path = "/tfxun/update"
  host_header_id = "${ksyun_lb_host_header.default.id}"
  backend_server_group_id="${ksyun_lb_backend_server_group.default.id}"
  listener_sync="on"
  method="RoundRobin"
  session {
    session_state = "start"
    session_persistence_period = 1000
    cookie_type = "ImplantCookie"
  #  cookie_name = "cookiexunqq"
  }
  health_check{
    health_check_state = "start"
    healthy_threshold = 2
    interval = 20
    timeout = 20
    unhealthy_threshold = 2
    url_path = "/monitor"
    host_name = "www.ksyun.com"
  }
}`
