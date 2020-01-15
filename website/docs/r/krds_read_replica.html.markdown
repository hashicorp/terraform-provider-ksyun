#ksyun_krds_rr
Provides an RDS Read Only instance resource. A DB read only instance is an isolated database environment in the cloud. 
 
## Example Usage
»Create a RDS Read Only MySQL instance

```h
resource "ksyun_krds_rr" "my_rds_rr"{
  output_file = "output_file"
  db_instance_identifier= "******"
  db_instance_class= "db.ram.2|db.disk.50"
  db_instance_name = "houbin_terraform_888_rr_1"
  bill_type = "DAY"
  security_group_id = "******"

  parameters {
    name = "auto_increment_increment"
    value = "7"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `output_file`- (Required)  will return the file name of the content store
* `db_instance_identifier`- (Required) passes in the instance ID of the RDS highly available instance. A RDS highly available instance can have at most three read-only instances
* `db_instance_class`-（Required）-this value regex db.ram.d{1,3}|db.disk.d{1,5} , db.ram is rds random access memory size, db.disk is disk size
* `db_instance_name`- (Required)instance name
* `bill_type`- (Required)Bill type, year'month (monthly package), day (daily billing), default: year'month
* `duration`- (Optional) purchase duration in months
* `security_group_id`- (Optional)security group ID
* `project_id`- (Optional) subproject ID
* `parameters`- (Optional) database parameters

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* ` DBInstanceClass`- instance specification
* `  Vcpus`-  number of CPUs
* `  Disk`-   hard disk size
* `  Ram `-   memory size
* `DBInstanceIdentifier`-  instance ID
* `DBInstanceName`-    instance name
* `DBInstanceStatus `- instance status
* `DBInstanceType `-  instance type
* `DBParameterGroupId `-  parameter group ID
* `GroupId `-  group ID
* `SecurityGroupId`-  security group ID
* `Vip`-  virtual IP
* `Port `- port number
* `Engine `-  Database Engine
* `EngineVersion`-   database engine version
* `InstanceCreateTime `- instance creation time
* `MasterUserName `-  primary account user name
* `DatastoreVersionId `- database version
* `Region `- region
* `VpcId `-virtual private network ID
* `ReadReplicaDBInstanceIdentifiers`-  read only instance
* `BillType `- Bill type
* `MultiAvailabilityZone`-  Multi availability zone
* `ProductId`- Product ID
* `DiskUsed`-  hard disk usage
* `ProjectId`-  Project ID
* `ProjectName `- project name


NOTE:RDS RR do not support modify

»Attributes Reference
The following attributes are exported:
```
id - The RDS instance ID.
port - RDS database connection port.
```
»Timeouts
NOTE: Available in 1.52.1+.
```
The timeouts block allows you to specify timeouts for certain actions:

create - (Defaults to 30 mins) Used when creating the db instance (until it reaches the initial Running status).
update - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial Running status).
delete - (Defaults to 10 mins) Used when terminating the db instance.
```