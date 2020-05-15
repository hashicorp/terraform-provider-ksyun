package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccKsyunKrdsSecrityGroup_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_krds_security_group.krds_sec_group_234",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKrdsSecrityGroupDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccKrdsSecrityGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testCheckKrdsSecrityGroupExists("ksyun_krds_security_group.krds_sec_group_234", &val),
				),
			},
		},
	})
}

func testCheckKrdsSecrityGroupExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found : %s", n)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("instance is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"SecurityGroupId": res.Primary.ID,
		}
		resp, err := client.krdsconn.DescribeSecurityGroup(&req)
		if err != nil {
			return err
		}
		*val = *resp
		return nil
	}
}

func testAccCheckKrdsSecrityGroupDestroy(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "ksyun_krds" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		req := map[string]interface{}{
			"DBInstanceIdentifier": res.Primary.ID,
		}
		_, err := client.krdsconn.DescribeDBInstances(&req)
		if err != nil {
			if err.(awserr.Error).Code() == "NOT_FOUND" {
				return nil
			}
			return err
		}
	}

	return nil
}

const testAccKrdsSecrityGroupConfig = `


resource "ksyun_krds_security_group" "krds_sec_group_234" {
  output_file = "output_file"
  security_group_name = "terraform_security_group_234"
  security_group_description = "terraform-security-group-234"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "asdf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "asdf2"
  }
}



`
