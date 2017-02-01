package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["access_key"],
			},

			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["secret_key"],
			},

			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["profile"],
			},

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["shared_credentials_file"],
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OSC_REGION",
					"OSC_DEFAULT_REGION",
				}, nil),
				Description:  descriptions["region"],
				InputDefault: "eu-west-2",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     11,
				Description: descriptions["max_retries"],
			},

			"allowed_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"forbidden_account_ids"},
				Set:           schema.HashString,
			},

			"forbidden_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"allowed_account_ids"},
				Set:           schema.HashString,
			},

			"endpoints": endpointsSchema(),

			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},

			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_region_validation"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"osc_omi":                    dataSourceOscOmi(),
			"osc_availability_zone":      dataSourceOscAvailabilityZone(),
			"osc_availability_zones":     dataSourceOscAvailabilityZones(),
			"osc_ebs_snapshot":           dataSourceOscEbsSnapshot(),
			"osc_ebs_volume":             dataSourceOscEbsVolume(),
			"osc_eip":                    dataSourceOscEip(),
			"osc_lbu_hosted_zone_id":     dataSourceOscLbuHostedZoneId(),
			"osc_lbu_service_account":    dataSourceOscLbuServiceAccount(),
			"osc_eim_account_alias":      dataSourceOscEimAccountAlias(),
			"osc_eim_policy_document":    dataSourceOscEimPolicyDocument(),
			"osc_eim_server_certificate": dataSourceOscIAMServerCertificate(),
			"osc_instance":               dataSourceOscInstance(),
			"osc_ip_ranges":              dataSourceOscIPRanges(),
			"osc_prefix_list":            dataSourceOscPrefixList(),
			"osc_region":                 dataSourceOscRegion(),
			"osc_route_table":            dataSourceOscRouteTable(),
			"osc_osu_bucket_object":      dataSourceOscOsuBucketObject(),
			"osc_subnet":                 dataSourceOscSubnet(),
			"osc_security_group":         dataSourceOscSecurityGroup(),
			"osc_vpc":                    dataSourceOscVpc(),
			"osc_vpc_endpoint":           dataSourceOscVpcEndpoint(),
			"osc_vpc_endpoint_service":   dataSourceOscVpcEndpointService(),
			"osc_vpc_peering_connection": dataSourceOscVpcPeeringConnection(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"osc_omi":                                  resourceOscOmi(),
			"osc_omi_copy":                             resourceOscOmiCopy(),
			"osc_omi_from_instance":                    resourceOscOmiFromInstance(),
			"osc_omi_launch_permission":                resourceOscOmiLaunchPermission(),
			"osc_app_cookie_stickiness_policy":         resourceOscAppCookieStickinessPolicy(),
			"osc_customer_gateway":                     resourceOscCustomerGateway(),
			"osc_directory_service_directory":          resourceOscDirectoryServiceDirectory(),
			"osc_ebs_snapshot":                         resourceOscEbsSnapshot(),
			"osc_ebs_volume":                           resourceOscEbsVolume(),
			"osc_eip":                                  resourceOscEip(),
			"osc_eip_association":                      resourceOscEipAssociation(),
			"osc_elasticache_cluster":                  resourceOscElasticacheCluster(),
			"osc_lbu":                                  resourceOscLbu(),
			"osc_lbu_attachment":                       resourceOscLbuAttachment(),
			"osc_eim_access_key":                       resourceOscEimAccessKey(),
			"osc_eim_account_password_policy":          resourceOscEimAccountPasswordPolicy(),
			"osc_eim_group_policy":                     resourceOscEimGroupPolicy(),
			"osc_eim_group":                            resourceOscEimGroup(),
			"osc_eim_group_membership":                 resourceOscEimGroupMembership(),
			"osc_eim_group_policy_attachment":          resourceOscEimGroupPolicyAttachment(),
			"osc_eim_instance_profile":                 resourceOscEimInstanceProfile(),
			"osc_eim_policy":                           resourceOscEimPolicy(),
			"osc_eim_policy_attachment":                resourceOscEimPolicyAttachment(),
			"osc_eim_role_policy_attachment":           resourceOscEimRolePolicyAttachment(),
			"osc_eim_role_policy":                      resourceOscEimRolePolicy(),
			"osc_eim_role":                             resourceOscEimRole(),
			"osc_eim_saml_provider":                    resourceOscEimSamlProvider(),
			"osc_eim_server_certificate":               resourceOscIAMServerCertificate(),
			"osc_eim_user_policy_attachment":           resourceOscEimUserPolicyAttachment(),
			"osc_eim_user_policy":                      resourceOscEimUserPolicy(),
			"osc_eim_user_ssh_key":                     resourceOscEimUserSshKey(),
			"osc_eim_user":                             resourceOscEimUser(),
			"osc_eim_user_login_profile":               resourceOscEimUserLoginProfile(),
			"osc_instance":                             resourceOscInstance(),
			"osc_internet_gateway":                     resourceOscInternetGateway(),
			"osc_key_pair":                             resourceOscKeyPair(),
			"osc_lb_cookie_stickiness_policy":          resourceOscLBCookieStickinessPolicy(),
			"osc_load_balancer_policy":                 resourceOscLoadBalancerPolicy(),
			"osc_load_balancer_backend_server_policy":  resourceOscLoadBalancerBackendServerPolicies(),
			"osc_load_balancer_listener_policy":        resourceOscLoadBalancerListenerPolicies(),
			"osc_lb_ssl_negotiation_policy":            resourceOscLBSSLNegotiationPolicy(),
			"osc_main_route_table_association":         resourceOscMainRouteTableAssociation(),
			"osc_nat_gateway":                          resourceOscNatGateway(),
			"osc_network_acl":                          resourceOscNetworkAcl(),
			"osc_default_network_acl":                  resourceOscDefaultNetworkAcl(),
			"osc_default_route_table":                  resourceOscDefaultRouteTable(),
			"osc_network_acl_rule":                     resourceOscNetworkAclRule(),
			"osc_network_interface":                    resourceOscNetworkInterface(),
			"osc_placement_group":                      resourceOscPlacementGroup(),
			"osc_proxy_protocol_policy":                resourceOscProxyProtocolPolicy(),
			"osc_route":                                resourceOscRoute(),
			"osc_route_table":                          resourceOscRouteTable(),
			"osc_route_table_association":              resourceOscRouteTableAssociation(),
			"osc_default_security_group":               resourceOscDefaultSecurityGroup(),
			"osc_security_group":                       resourceOscSecurityGroup(),
			"osc_security_group_rule":                  resourceOscSecurityGroupRule(),
			"osc_snapshot_create_volume_permission":    resourceOscSnapshotCreateVolumePermission(),
			"osc_subnet":                               resourceOscSubnet(),
			"osc_volume_attachment":                    resourceOscVolumeAttachment(),
			"osc_vpc_dhcp_options_association":         resourceOscVpcDhcpOptionsAssociation(),
			"osc_vpc_dhcp_options":                     resourceOscVpcDhcpOptions(),
			"osc_vpc_peering_connection":               resourceOscVpcPeeringConnection(),
			"osc_vpc":                                  resourceOscVpc(),
			"osc_vpc_endpoint":                         resourceOscVpcEndpoint(),
			"osc_vpc_endpoint_route_table_association": resourceOscVpcEndpointRouteTableAssociation(),
			"osc_vpn_connection":                       resourceOscVpnConnection(),
			"osc_vpn_connection_route":                 resourceOscVpnConnectionRoute(),
			"osc_vpn_gateway":                          resourceOscVpnGateway(),
			"osc_vpn_gateway_attachment":               resourceOscVpnGatewayAttachment(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"region": "The region where OSC operations will take place. Examples\n" +
			"are eu-west-2, us-east-2 etc.",

		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the Cockpit Console ",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the Cockpit Console ",

		"profile": "The profile for API operations. If not set, the default profile\n" +
			"created with `aws configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set\n" +
			"this defaults to ~/.aws/credentials.",

		"max_retries": "The maximum number of times an OSC API request is\n" +
			"being executed. If the API request still fails, an error is\n" +
			"thrown.",

		"eim_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"fcu_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"lbu_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"osu_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"skip_region_validation": "Skip static validation of region name. " +
			"Used by users of alternative OSC-like APIs or users w/ access to regions that are not public (yet).",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey:            d.Get("access_key").(string),
		SecretKey:            d.Get("secret_key").(string),
		Profile:              d.Get("profile").(string),
		CredsFilename:        d.Get("shared_credentials_file").(string),
		Region:               d.Get("region").(string),
		MaxRetries:           d.Get("max_retries").(int),
		Insecure:             d.Get("insecure").(bool),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		config.AssumeRoleARN = assumeRole["role_arn"].(string)
		config.AssumeRoleSessionName = assumeRole["session_name"].(string)
		config.AssumeRoleExternalID = assumeRole["external_id"].(string)
		log.Printf("[INFO] assume_role configuration set: (ARN: %q, SessionID: %q, ExternalID: %q)",
			config.AssumeRoleARN, config.AssumeRoleSessionName, config.AssumeRoleExternalID)
	} else {
		log.Printf("[INFO] No assume_role block read from configuration")
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.EimEndpoint = endpoints["eim"].(string)
		config.FcuEndpoint = endpoints["fcu"].(string)
		config.LbuEndpoint = endpoints["lbu"].(string)
		config.OsuEndpoint = endpoints["osu"].(string)
	}

	if v, ok := d.GetOk("allowed_account_ids"); ok {
		config.AllowedAccountIds = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("forbidden_account_ids"); ok {
		config.ForbiddenAccountIds = v.(*schema.Set).List()
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var awsMutexKV = mutexkv.NewMutexKV()

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"eim": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eim_endpoint"],
				},

				"fcu": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fcu_endpoint"],
				},

				"lbu": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["lbu_endpoint"],
				},
				"osu": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["osu_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["eim"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fcu"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["lbu"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["osu"].(string)))

	return hashcode.String(buf.String())
}
