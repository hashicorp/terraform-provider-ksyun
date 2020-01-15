# ksyun_sqlserver
Provides an SqlServer instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

## Example Usage
»Create a RDS SqlServer instance

```h
provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = ""
  secret_key = ""
}

variable "available_zone" {
  default = "cn-shanghai-3a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_sqlserver" "sqlserver-1"{
  output_file = "output_file"
  dbinstanceclass= "db.ram.2|db.disk.20"
  dbinstancename = "ksyun_sqlserver_1"
  dbinstancetype = "HRDS_SS"
  engine = "SQLServer"
  engineversion = "2008r2"
  masterusername = "admin"
  masteruserpassword = "123qweASD"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  billtype = "DAY"
}
```
## Argument Reference

The following arguments are supported:

* `output_file` -(Required) will return the file name of the content store
* `db_instance_class`-（Required）-this value regex db.ram.d{1,3}|db.disk.d{1,5} , db.ram is rds random access memory size, db.disk is disk size
* `db_instance_name`-(Required)instance name
* `db_instance_type`- (Required)instance type supports HRDS_SS
* `engine (Required)`-engine is db type, only support SQLServer
* `engine_version`- (Required)db engine version only support 2008r2,2012,2016
* `master_user_name`- (Required)database primary account name
* `master_user_password `-(Required) master account password
* `vpc_id `- (Required)ID of virtual private network
* `subnet_id`- (Required)subnet ID
* `bill_type`- (Required) Bill type, year'month (monthly package), day (daily billing), default: year'month
* `duration`- (Optional) purchase duration in months
* `security_group_id `-(Optional) security group ID
* `preferred_backup_time`- (Optional) backup time
* `availability_zone_1`- (Optional) zone 1
* `availability_zone_2`- (Optional) zone 2
* `project_id`- (Optional)  subproject ID
* `parameters`- (Optional) database parameters
* `port `-(Optional) port number

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

NOTE: SQLServer not support modify

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