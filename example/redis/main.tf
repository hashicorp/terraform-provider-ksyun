# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_vpc" "default" {
  vpc_name = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "ksyun_subnet" "default" {
  subnet_name = "${var.subnet_name}"
  cidr_block = "10.1.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.1.0.2"
  dhcp_ip_to = "10.1.0.253"
  vpc_id = "${ksyun_vpc.default.id}"
  gateway_ip = "10.1.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_redis_instance" "default" {
  available_zone = "${var.available_zone}"
  name = "MyRedisInstance1101"
  mode = 2
  capacity = 1
  net_type = 2
  vnet_id = "${ksyun_subnet.default.id}"
  vpc_id = "${ksyun_vpc.default.id}"
  bill_type = 5
  duration = ""
  duration_unit = ""
  pass_word = "Shiwo1101"
  iam_project_id = "0"
  slave_num = 0
  protocol = "${var.protocol}"
  reset_all_parameters = false
  parameters = {
    "appendonly" = "no",
    "appendfsync" = "everysec",
    "maxmemory-policy" = "volatile-lru",
    "hash-max-ziplist-entries" = "512",
    "zset-max-ziplist-entries" = "128",
    "list-max-ziplist-size" = "-2",
    "hash-max-ziplist-value" = "64",
    "notify-keyspace-events" = "",
    "zset-max-ziplist-value" = "64",
    "maxmemory-samples" = "5",
    "set-max-intset-entries" = "512",
    "timeout" = "600",
  }
}

resource "ksyun_redis_instance_node" "default" {
  cache_id = "${ksyun_redis_instance.default.id}"
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_instance_node" "node" {
  // creating multiple read-only nodes,
  // not concurrently, requires dependencies to synchronize the execution of creating multiple read-only nodes.
  // if only one read-only node is created, it is not required to fill in.
  pre_node_id = "${ksyun_redis_instance_node.default.id}"
  cache_id = "${ksyun_redis_instance.default.id}"
  available_zone = "${var.available_zone}"
}

resource "ksyun_redis_sec_rule" "default" {
  cache_id = "${ksyun_redis_instance.default.id}"
  available_zone = "${var.available_zone}"
  rules = ["10.0.3.4/32"]
}