---
layout: "ksyun"
page_title: "Ksyun: ksyun_ks3"
sidebar_current: "docs-ksyun-resource-ks3"
description: |-
  Provides a KS3 resource.
---

# ks3_bucket_html

Provides a resource to create a KS3 bucket and set its attribution.

> **Note**  The bucket namespace is shared by all users of the KS3 system. Please set bucket name as unique as possible.

## Example Usage

Create Bucket
```hcl
resource "ksyun_ks3" "bucket-create" {
  bucket = "ks3-bucket-create"
}
```

Change Bucket ACL
```hcl
resource "ksyun_ks3" "bucket-acl" {
  bucket = "ks3-bucket-acl"
  acl    = "private"
}
```

Enable Bucket Logging
```hcl
resource "ksyun_ks3" "bucket-target" {
  bucket = "ks3-bucket-target"
  acl    = "public-read"
}

resource "ksyun_ks3" "bucket-logging" {
  bucket = "ks3-bucket-logging"

  logging {
    target_bucket = "${ksyun_ks3.bucket-target.id}"
  }
}
```

Set Bucket CORS
```hcl
resource "ksyun_ks3" "bucket-cors" {
  bucket = "ks3-bucket-cors"
  acl    = "public-read"

  cors_rule {
    allowed_header = ["*"]
    allowed_method = ["PUT", "POST"]
    allowed_origin = ["https://www.example.com"]
    expose_header  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, Forces new resource) The name of the bucket. The length should be between 3 ~ 63.
* `acl` - (Optional) The canned ACL to apply. Defaults to "private".
* `cors_rule` - (Optional) A rule of Cross-Origin Resource Sharing.
    * `allowed_header` - (Required) Specifies which headers are allowed.
    * `allowed_method` - (Required) Specifies which methods are allowed. Can be GET, PUT, POST, DELETE or HEAD.
    * `allowed_origin` - (Required) Specifies which origins are allowed.
    * `expose_header` - (Optional) Specifies expose header in the response.
    * `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.
* `logging` - (Optional) A settings of bucket logging.
    * `target_bucket` - (Required) The name of the bucket that will receive the log objects.
    * `target_prefix` - (Optional) To specify a key prefix for log objects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The name of the bucket.
* `acl` - The acl of the bucket.

## Import

KS3 bucket can be imported using the bucket name, e.g.

```
$ terraform import ksyun_ks3.example bucket-12345678
```