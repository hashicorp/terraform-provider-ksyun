# ksyun_krds_security_group
Provide RDS security group function
## Example Usage
»Create a RDS Security GroupL instance

```h
»Create a RDS Security Group instance
resource "ksyun_krds_security_group" "krds_sec_group_13" {
  output_file = "output_file"
  security_group_name = "terraform_security_group_13"
  security_group_description = "terraform-security-group-13"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "asdf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "asdf2"
  }

}
```
## Argument Reference

The following arguments are supported:

* `output_file`- (Required)  will return the file name of the content store
* `security_group_name `-(Required)  the name of the security group
* `security_group_description`-（Optional）description of security group
* `security_group_rule`- (Optional)security group rule
* `security_group_rule_protocol`- (Required)  0.0.0.0/32 format
* `security_group_rule_name`- (Required) no more than 256 bytes, only Chinese, uppercase and lowercase letters, numbers, minus signs and underscores are supported

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `SecurityGroupId`- Security group ID
* `SecurityGroupName`- security group name
* `SecurityGroupDescription`- Security Group Description
* `Instances`- corresponding instance 
* `DBInstanceIdentifier`- instance ID
* `DBInstanceName`-instance name
* `Vip`- instance virtual IP
* `SecurityGroupRules`- security group rules
* `SecurityGroupRuleId`-rule ID
* `SecurityGroupRuleName`-rule name
* `SecurityGroupRuleProtocol`- rule protocol

»Timeouts
NOTE: Available in 1.52.1+.
```
The timeouts block allows you to specify timeouts for certain actions:

create - (Defaults to 10 mins) Used when creating the db instance (until it reaches the initial Running status).
update - (Defaults to 10 mins) Used when updating the db instance (until it reaches the initial Running status).
delete - (Defaults to 10 mins) Used when terminating the db instance.
```