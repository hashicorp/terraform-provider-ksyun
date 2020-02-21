#ksyun_krds_security_groups 
Query security group information
## Example Usage
```
data "ksyun_sqlservers" "search-sqlservers"{
  output_file = "output_file"
  security_group_id = 123
}
```
## Argument Reference

The following arguments are supported:

* `output_file`- (Required) The filename of the content store will be returned
* `security_group_id`- (Optional) Security group ID

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