package osc

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

			"assume_role": assumeRoleSchema(),

			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["shared_credentials_file"],
			},

			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["token"],
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				Description:  descriptions["region"],
				InputDefault: "us-east-1",
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

			"skip_metadata_api_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_metadata_api_check"],
			},

			"s3_force_path_style": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["s3_force_path_style"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"osc_ami":                     dataSourceAwsAmi(),
			"osc_availability_zone":       dataSourceAwsAvailabilityZone(),
			"osc_availability_zones":      dataSourceAwsAvailabilityZones(),
			"osc_billing_service_account": dataSourceAwsBillingServiceAccount(),
			"osc_caller_identity":         dataSourceAwsCallerIdentity(),
			"osc_canonical_user_id":       dataSourceAwsCanonicalUserId(),
			"osc_ebs_snapshot":            dataSourceAwsEbsSnapshot(),
			"osc_ebs_volume":              dataSourceAwsEbsVolume(),
			"osc_eip":                     dataSourceAwsEip(),
			"osc_elb_hosted_zone_id":      dataSourceAwsElbHostedZoneId(),
			"osc_elb_service_account":     dataSourceAwsElbServiceAccount(),
			"osc_iam_account_alias":       dataSourceAwsIamAccountAlias(),
			"osc_iam_policy_document":     dataSourceAwsIamPolicyDocument(),
			"osc_iam_server_certificate":  dataSourceAwsIAMServerCertificate(),
			"osc_instance":                dataSourceAwsInstance(),
			"osc_ip_ranges":               dataSourceAwsIPRanges(),
			"osc_partition":               dataSourceAwsPartition(),
			"osc_prefix_list":             dataSourceAwsPrefixList(),
			"osc_region":                  dataSourceAwsRegion(),
			"osc_route_table":             dataSourceAwsRouteTable(),
			"osc_s3_bucket_object":        dataSourceAwsS3BucketObject(),
			"osc_subnet":                  dataSourceAwsSubnet(),
			"osc_security_group":          dataSourceAwsSecurityGroup(),
			"osc_vpc":                     dataSourceAwsVpc(),
			"osc_vpc_endpoint":            dataSourceAwsVpcEndpoint(),
			"osc_vpc_endpoint_service":    dataSourceAwsVpcEndpointService(),
			"osc_vpc_peering_connection":  dataSourceAwsVpcPeeringConnection(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"osc_ami":                                 resourceAwsAmi(),
			"osc_ami_copy":                            resourceAwsAmiCopy(),
			"osc_ami_from_instance":                   resourceAwsAmiFromInstance(),
			"osc_ami_launch_permission":               resourceAwsAmiLaunchPermission(),
			"osc_api_gateway_account":                 resourceAwsApiGatewayAccount(),
			"osc_api_gateway_api_key":                 resourceAwsApiGatewayApiKey(),
			"osc_api_gateway_authorizer":              resourceAwsApiGatewayAuthorizer(),
			"osc_api_gateway_base_path_mapping":       resourceAwsApiGatewayBasePathMapping(),
			"osc_api_gateway_client_certificate":      resourceAwsApiGatewayClientCertificate(),
			"osc_api_gateway_integration":             resourceAwsApiGatewayIntegration(),
			"osc_api_gateway_integration_response":    resourceAwsApiGatewayIntegrationResponse(),
			"osc_api_gateway_method":                  resourceAwsApiGatewayMethod(),
			"osc_api_gateway_method_response":         resourceAwsApiGatewayMethodResponse(),
			"osc_api_gateway_model":                   resourceAwsApiGatewayModel(),
			"osc_api_gateway_resource":                resourceAwsApiGatewayResource(),
			"osc_api_gateway_rest_api":                resourceAwsApiGatewayRestApi(),
			"osc_app_cookie_stickiness_policy":        resourceAwsAppCookieStickinessPolicy(),
			"osc_customer_gateway":                    resourceAwsCustomerGateway(),
			"osc_ebs_snapshot":                        resourceAwsEbsSnapshot(),
			"osc_ebs_volume":                          resourceAwsEbsVolume(),
			"osc_eip":                                 resourceAwsEip(),
			"osc_eip_association":                     resourceAwsEipAssociation(),
			"osc_elb":                                 resourceAwsElb(),
			"osc_elb_attachment":                      resourceAwsElbAttachment(),
			"osc_flow_log":                            resourceAwsFlowLog(),
			"osc_iam_access_key":                      resourceAwsIamAccessKey(),
			"osc_iam_account_password_policy":         resourceAwsIamAccountPasswordPolicy(),
			"osc_iam_group_policy":                    resourceAwsIamGroupPolicy(),
			"osc_iam_group":                           resourceAwsIamGroup(),
			"osc_iam_group_membership":                resourceAwsIamGroupMembership(),
			"osc_iam_group_policy_attachment":         resourceAwsIamGroupPolicyAttachment(),
			"osc_iam_instance_profile":                resourceAwsIamInstanceProfile(),
			"osc_iam_policy":                          resourceAwsIamPolicy(),
			"osc_iam_policy_attachment":               resourceAwsIamPolicyAttachment(),
			"osc_iam_role_policy_attachment":          resourceAwsIamRolePolicyAttachment(),
			"osc_iam_role_policy":                     resourceAwsIamRolePolicy(),
			"osc_iam_role":                            resourceAwsIamRole(),
			"osc_iam_saml_provider":                   resourceAwsIamSamlProvider(),
			"osc_iam_server_certificate":              resourceAwsIAMServerCertificate(),
			"osc_iam_user_policy_attachment":          resourceAwsIamUserPolicyAttachment(),
			"osc_iam_user_policy":                     resourceAwsIamUserPolicy(),
			"osc_iam_user_ssh_key":                    resourceAwsIamUserSshKey(),
			"osc_iam_user":                            resourceAwsIamUser(),
			"osc_iam_user_login_profile":              resourceAwsIamUserLoginProfile(),
			"osc_instance":                            resourceAwsInstance(),
			"osc_internet_gateway":                    resourceAwsInternetGateway(),
			"osc_key_pair":                            resourceAwsKeyPair(),
			"osc_lb_cookie_stickiness_policy":         resourceAwsLBCookieStickinessPolicy(),
			"osc_load_balancer_policy":                resourceAwsLoadBalancerPolicy(),
			"osc_load_balancer_backend_server_policy": resourceAwsLoadBalancerBackendServerPolicies(),
			"osc_load_balancer_listener_policy":       resourceAwsLoadBalancerListenerPolicies(),
			"osc_lb_ssl_negotiation_policy":           resourceAwsLBSSLNegotiationPolicy(),
			"osc_main_route_table_association":        resourceAwsMainRouteTableAssociation(),
			"osc_nat_gateway":                         resourceAwsNatGateway(),
			"osc_network_acl":                         resourceAwsNetworkAcl(),
			"osc_default_network_acl":                 resourceAwsDefaultNetworkAcl(),
			"osc_default_route_table":                 resourceAwsDefaultRouteTable(),
			"osc_network_acl_rule":                    resourceAwsNetworkAclRule(),
			"osc_network_interface":                   resourceAwsNetworkInterface(),
			"osc_placement_group":                     resourceAwsPlacementGroup(),
			"osc_proxy_protocol_policy":               resourceAwsProxyProtocolPolicy(),
			"osc_route":                               resourceAwsRoute(),
			"osc_route_table":                         resourceAwsRouteTable(),
			"osc_route_table_association":             resourceAwsRouteTableAssociation(),
			"osc_s3_bucket":                           resourceAwsS3Bucket(),
			"osc_s3_bucket_policy":                    resourceAwsS3BucketPolicy(),
			"osc_s3_bucket_object":                    resourceAwsS3BucketObject(),
			"osc_s3_bucket_notification":              resourceAwsS3BucketNotification(),
			"osc_default_security_group":              resourceAwsDefaultSecurityGroup(),
			"osc_security_group":                      resourceAwsSecurityGroup(),
			"osc_security_group_rule":                 resourceAwsSecurityGroupRule(),
			"osc_snapshot_create_volume_permission":   resourceAwsSnapshotCreateVolumePermission(),
			"osc_subnet":                              resourceAwsSubnet(),
			"osc_volume_attachment":                   resourceAwsVolumeAttachment(),
			"osc_vpc_dhcp_options_association":        resourceAwsVpcDhcpOptionsAssociation(),
			"osc_vpc_dhcp_options":                    resourceAwsVpcDhcpOptions(),
			"osc_vpc_peering_connection":              resourceAwsVpcPeeringConnection(),
			"osc_vpc_peering_connection_accepter":     resourceAwsVpcPeeringConnectionAccepter(),
			"osc_vpc":                                  resourceAwsVpc(),
			"osc_vpc_endpoint":                         resourceAwsVpcEndpoint(),
			"osc_vpc_endpoint_route_table_association": resourceAwsVpcEndpointRouteTableAssociation(),
			"osc_vpn_connection":                       resourceAwsVpnConnection(),
			"osc_vpn_connection_route":                 resourceAwsVpnConnectionRoute(),
			"osc_vpn_gateway":                          resourceAwsVpnGateway(),
			"osc_vpn_gateway_attachment":               resourceAwsVpnGatewayAttachment(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"region": "The region where AWS operations will take place. Examples\n" +
			"are us-east-1, us-west-2, etc.",

		"access_key": "The access key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"secret_key": "The secret key for API operations. You can retrieve this\n" +
			"from the 'Security & Credentials' section of the AWS console.",

		"profile": "The profile for API operations. If not set, the default profile\n" +
			"created with `aws configure` will be used.",

		"shared_credentials_file": "The path to the shared credentials file. If not set\n" +
			"this defaults to ~/.aws/credentials.",

		"token": "session token. A session token is only required if you are\n" +
			"using temporary security credentials.",

		"max_retries": "The maximum number of times an AWS API request is\n" +
			"being executed. If the API request still fails, an error is\n" +
			"thrown.",

		"iam_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"ec2_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"elb_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"s3_endpoint": "Use this to override the default endpoint URL constructed from the `region`.\n",

		"insecure": "Explicitly allow the provider to perform \"insecure\" SSL requests. If omitted," +
			"default value is `false`",

		"skip_region_validation": "Skip static validation of region name. " +
			"Used by users of alternative AWS-like APIs or users w/ access to regions that are not public (yet).",

		"skip_medatadata_api_check": "Skip the AWS Metadata API check. " +
			"Used for AWS API implementations that do not have a metadata api endpoint.",

		"s3_force_path_style": "Set this to true to force the request to use path-style addressing,\n" +
			"i.e., http://s3.amazonaws.com/BUCKET/KEY. By default, the S3 client will\n" +
			"use virtual hosted bucket addressing when possible\n" +
			"(http://BUCKET.s3.amazonaws.com/KEY). Specific to the Amazon S3 service.",

		"assume_role_role_arn": "The ARN of an IAM role to assume prior to making API calls.",

		"assume_role_session_name": "The session name to use when assuming the role. If omitted," +
			" no session name is passed to the AssumeRole call.",

		"assume_role_external_id": "The external ID to use when assuming the role. If omitted," +
			" no external ID is passed to the AssumeRole call.",

		"assume_role_policy": "The permissions applied when assuming a role. You cannot use," +
			" this policy to grant further permissions that are in excess to those of the, " +
			" role that is being assumed.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey:               d.Get("access_key").(string),
		SecretKey:               d.Get("secret_key").(string),
		Profile:                 d.Get("profile").(string),
		CredsFilename:           d.Get("shared_credentials_file").(string),
		Token:                   d.Get("token").(string),
		Region:                  d.Get("region").(string),
		MaxRetries:              d.Get("max_retries").(int),
		Insecure:                d.Get("insecure").(bool),
		SkipRegionValidation:    d.Get("skip_region_validation").(bool),
		SkipMetadataApiCheck:    d.Get("skip_metadata_api_check").(bool),
		S3ForcePathStyle:        d.Get("s3_force_path_style").(bool),
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		config.AssumeRoleARN = assumeRole["role_arn"].(string)
		config.AssumeRoleSessionName = assumeRole["session_name"].(string)
		config.AssumeRoleExternalID = assumeRole["external_id"].(string)

		if v := assumeRole["policy"].(string); v != "" {
			config.AssumeRolePolicy = v
		}

		log.Printf("[INFO] assume_role configuration set: (ARN: %q, SessionID: %q, ExternalID: %q, Policy: %q)",
			config.AssumeRoleARN, config.AssumeRoleSessionName, config.AssumeRoleExternalID, config.AssumeRolePolicy)
	} else {
		log.Printf("[INFO] No assume_role block read from configuration")
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.IamEndpoint = endpoints["iam"].(string)
		config.Ec2Endpoint = endpoints["ec2"].(string)
		config.ElbEndpoint = endpoints["elb"].(string)
		config.S3Endpoint = endpoints["s3"].(string)
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

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_role_arn"],
				},

				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
				},

				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_external_id"],
				},

				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
			},
		},
		Set: assumeRoleToHash,
	}
}

func assumeRoleToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["role_arn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["session_name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["external_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["policy"].(string)))
	return hashcode.String(buf.String())
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"iam": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["iam_endpoint"],
				},

				"ec2": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ec2_endpoint"],
				},

				"elb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elb_endpoint"],
				},
				"s3": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["s3_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["iam"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ec2"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["s3"].(string)))

	return hashcode.String(buf.String())
}
