---
layout: "ksyun"
page_title: "Ksyun: ksyun_krds"
sidebar_current: "docs-ksyun-datasource-krds"
description: |-
  Provides a list of krds resources in the current region.
---

# ksyun_krds

Query HRDS and RDS-rr instance information

## Example Usage

```hcl
# Get  krds
data "ksyun_krds" "search-krds"{
  output_file = "output_file"
  db_instance_identifier = "***"
  db_instance_type = "HRDS,RR,TRDS"
  keyword = ""
  order = ""
  project_id = ""
  marker = ""
  max_records = ""
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Required) will return the file name of the content store
* `db_instance_identifier` - (Optional) instance ID (passed in the instance ID to get the details of the instance, otherwise get the list)
* `db_instance_type` - (Optional) hrds (highly available), RR (read-only), trds (temporary)
* `db_instance_status ` -(Optional) active / invalid (please renew)
* `keyword` -(Optional) fuzzy filter by name / VIP
* `order` - (Optional) case sensitive, value range: default (default sorting method), group (sorting by replication group, will rank read-only instances after their primary instances)
* `project_id` - (Optional) the default value is all projects
* `Marker` -(Optional) record start offset
* `MaxRecords` -(Optional) the maximum number of entries in the result of each page. Value range: 1-100

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
* `Eip`-  elastic IP address
* `EipPort`-  elastic IP port