---
layout: "ksyun"
page_title: "Ksyun: ksyun_images"
sidebar_current: "docs-ksyun-datasource-images"
description: |-
  Provides a list of available image resources in the current region.
---

# ksyun_images

This data source providers a list of available image resources according to their availability zone, image ID and other fields.

## Example Usage

```h
# Get  ksyun_images
data "ksyun_images" "default" {
  output_file="output_result"
  is_public=true
  image_source="system"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of image IDs
* `name_regex` - (Optional) A regex string to filter resulting images by name. (Such as: `^CentOS 7.[1-2] 64` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system, "^Ubuntu 16.04 64" means Ubuntu 16.04 of 64-bit operating system).
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `most_recent` - (Optional, type: bool) If more than one result are returned, select the most recent one.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `creation_date` - Time of creation.
* `image_id` -  The ID of image.
* `image_source` -  Valid values are import, copy, share, extend, system.
* `image_state` - Status of the image.
* `is_public` - If ksyun provide the image. 
* `name` - Display name of the image.
* `platform` -  Platform type of the image system.
* `progress` - Progress of image creation.
* `sys_disk` - Size of the created disk.

