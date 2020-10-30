package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_REGION", nil),
				Description: descriptions["region"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["insecure"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ksyun_lines":                         dataSourceKsyunLines(),
			"ksyun_eips":                          dataSourceKsyunEips(),
			"ksyun_slbs":                          dataSourceKsyunLbs(),
			"ksyun_lbs":                           dataSourceKsyunLbs(),
			"ksyun_listeners":                     dataSourceKsyunListeners(),
			"ksyun_health_checks":                 dataSourceKsyunHealthChecks(),
			"ksyun_listener_servers":              dataSourceKsyunLbListenerServers(),
			"ksyun_lb_listener_servers":           dataSourceKsyunLbListenerServers(),
			"ksyun_lb_acls":                       dataSourceKsyunSlbAcls(),
			"ksyun_availability_zones":            dataSourceKsyunAvailabilityZones(),
			"ksyun_network_interfaces":            dataSourceKsyunNetworkInterfaces(),
			"ksyun_vpcs":                          dataSourceKsyunVPCs(),
			"ksyun_subnets":                       dataSourceKsyunSubnets(),
			"ksyun_subnet_available_addresses":    dataSourceKsyunSubnetAvailableAddresses(),
			"ksyun_subnet_allocated_ip_addresses": dataSourceKsyunSubnetAllocatedIpAddresses(),
			"ksyun_security_groups":               dataSourceKsyunSecurityGroups(),
			"ksyun_instances":                     dataSourceKsyunInstances(),
			"ksyun_images":                        dataSourceKsyunImages(),
			"ksyun_sqlservers":                    dataSourceKsyunSqlServer(),
			"ksyun_krds":                          dataSourceKsyunKrds(),
			"ksyun_krds_security_groups":          dataSourceKsyunKrdsSecurityGroup(),
			"ksyun_certificates":                  dataSourceKsyunCertificates(),
			"ksyun_ssh_keys":                      dataSourceKsyunSSHKeys(),
			"ksyun_redis_instances":               dataSourceRedisInstances(),
			"ksyun_redis_security_groups":         dataSourceRedisSecurityGroups(),
			//	"ksyun_epcs":                          dataSourceKsyunEpcs(),
			"ksyun_volumes":                     dataSourceKsyunVolumes(),
			"ksyun_mongodbs":                    dataSourceKsyunMongodbs(),
			"ksyun_lb_host_headers":             dataSourceKsyunListenerHostHeaders(),
			"ksyun_lb_rules":                    dataSourceKsyunSlbRules(),
			"ksyun_lb_backend_server_groups":    dataSourceKsyunBackendServerGroups(),
			"ksyun_lb_register_backend_servers": dataSourceKsyunRegisterBackendServers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ksyun_eip":                       resourceKsyunEip(),
			"ksyun_eip_associate":             resourceKsyunEipAssociation(),
			"ksyun_lb":                        resourceKsyunLb(),
			"ksyun_healthcheck":               resourceKsyunHealthCheck(),
			"ksyun_lb_listener":               resourceKsyunListener(),
			"ksyun_lb_listener_server":        resourceKsyunInstancesWithListener(),
			"ksyun_lb_acl":                    resourceKsyunLoadBalancerAcl(),
			"ksyun_lb_acl_entry":              resourceKsyunLoadBalancerAclEntry(),
			"ksyun_lb_listener_associate_acl": resourceKsyunListenerLBAcl(),
			"ksyun_vpc":                       resourceKsyunVPC(),
			"ksyun_subnet":                    resourceKsyunSubnet(),
			"ksyun_security_group":            resourceKsyunSecurityGroup(),
			"ksyun_security_group_entry":      resourceKsyunSecurityGroupEntry(),
			"ksyun_instance":                  resourceKsyunInstance(),
			"ksyun_sqlserver":                 resourceKsyunSqlServer(),
			"ksyun_krds":                      resourceKsyunKrds(),
			"ksyun_krds_rr":                   resourceKsyunKrdsRr(),
			"ksyun_krds_security_group":       resourceKsyunKrdsSecurityGroup(),
			"ksyun_certificate":               resourceKsyunCertificate(),
			"ksyun_ssh_key":                   resourceKsyunSSHKey(),
			"ksyun_redis_instance":            resourceRedisInstance(),
			"ksyun_redis_instance_node":       resourceRedisInstanceNode(),
			"ksyun_redis_sec_group":           resourceRedisSecurityGroup(),
			"ksyun_redis_sec_group_rule":      resourceRedisSecurityGroupRule(),
			"ksyun_redis_sec_group_allocate":  resourceRedisSecurityGroupAllocate(),
			/*	"ksyun_epc":                       resourceKsyunEpc(),*/
			"ksyun_mongodb_instance":           resourceKsyunMongodbInstance(),
			"ksyun_mongodb_shard_instance":     resourceKsyunMongodbShardInstance(),
			"ksyun_mongodb_security_rule":      resourceKsyunMongodbSecurityRule(),
			"ksyun_volume":                     resourceKsyunVolume(),
			"ksyun_volume_attach":              resourceKsyunVolumeAttach(),
			"ksyun_lb_rule":                    resourceKsyunSlbRule(),
			"ksyun_lb_host_header":             resourceKsyunListenerHostHeader(),
			"ksyun_lb_backend_server_group":    resourceKsyunBackendServerGroup(),
			"ksyun_lb_register_backend_server": resourceKsyunRegisterBackendServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Region:    d.Get("region").(string),
		Insecure:  d.Get("insecure").(bool),
	}
	client, err := config.Client()
	return client, err
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "ak",
		"secret_key": "sk",
		"region":     "cn-beijing-6",
		"insecure":   "true",
	}
}
