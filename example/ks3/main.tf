# Specify the provider and access details
provider "ksyun" {
  access_key = "你的ak"
  secret_key = "你的sk"
  region = "ks3-cn-beijing.ksyun.com"
}

#Create Bucket
resource "ksyun_ks3" "bucket-create" {
  bucket = "ks3-bucket-create"
}

#Change Bucket ACL
resource "ksyun_ks3" "bucket-acl" {
  bucket = "ks3-bucket-acl"
  acl = "private"
}

#Enable Bucket Logging
resource "ksyun_ks3" "bucket-target" {
  bucket = "ks3-bucket-target"
  acl = "public-read"
}

resource "ksyun_ks3" "bucket-logging" {
  bucket = "ks3-bucket-logging"

  logging {
    target_bucket = "${ksyun_ks3.bucket-target.id}"
  }
}

#Set Bucket CORS
resource "ksyun_ks3" "bucket-cors" {
  bucket = "ks3-bucket-cors"
  acl = "public-read"

  cors_rule {
    allowed_header = ["*"]
    allowed_method = ["PUT", "POST"]
    allowed_origin = ["https://www.example.com"]
    expose_header = ["ETag"]
    max_age_seconds = 3000
  }
}
