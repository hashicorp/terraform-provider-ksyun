provider "ksyun"{
  region = "cn-shanghai-2"
}

resource "ksyun_krds_security_group" "krds_sec_group_14" {
  security_group_name = "terraform_security_group_14"
  security_group_description = "terraform-security-group-14"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "wtf"
  }

}