---
layout: "ksyun"
page_title: "Ksyun: ksyun_volumes"
sidebar_current: "docs-ksyun-datasource-volumes"
description: |-
  Provides a list of volume resources in the current region.
---

# ksyun_volumes

This data source providers a list of volume resources according to their availability zone.

## Example Usage

```h
data "ksyun_volumes" "default" {
  output_file="output_result"
  ids=[]
  volume_category=""
  volume_status=""
  volume_type=""
  availability_zone=""
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of volume IDs.Without this parameter, query all volumes' information.
* `volume_category` - (Optional) Volume classification,  system disk "system" or data disk "data".
* `volume_status` - (Optional) The status of volumes, “creating|available|attaching|in-use|detaching|extending|deleting|error|recycling”.
* `volume_type` - (Optional) The type of volumes. "SSD" or "SATA".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`)

