output "all_images" {
  value = "${data.ksyun_images.default.total_count}"
}
output "centos-7.5" {
  value = "${data.ksyun_images.centos-7_5.total_count}"
}