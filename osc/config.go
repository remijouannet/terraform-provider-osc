package osc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/terraform"
)

type Config struct {
	AccessKey     string
	SecretKey     string
	CredsFilename string
	Profile       string
	Token         string
	Region        string
	MaxRetries    int

	AssumeRoleARN         string
	AssumeRoleExternalID  string
	AssumeRoleSessionName string
	AssumeRolePolicy      string

	AllowedAccountIds   []interface{}
	ForbiddenAccountIds []interface{}

	Ec2Endpoint string
	IamEndpoint string
	ElbEndpoint string
	S3Endpoint  string
	Insecure    bool

	SkipRegionValidation bool
	SkipMetadataApiCheck bool
	S3ForcePathStyle     bool
}

type AWSClient struct {
	ec2conn    *ec2.EC2
	elbconn    *elb.ELB
	elbv2conn  *elbv2.ELBV2
	apigateway *apigateway.APIGateway
	s3conn     *s3.S3
	partition  string
	accountid  string
	region     string
	iamconn    *iam.IAM
}

// Client configures and returns a fully initialized AWSClient
func (c *Config) Client() (interface{}, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if c.SkipRegionValidation {
		log.Println("[INFO] Skipping region validation")
	} else {
		log.Println("[INFO] Building AWS region structure")
		err := c.ValidateRegion()
		if err != nil {
			return nil, err
		}
	}

	var client AWSClient
	// store AWS region in client struct, for region specific operations such as
	// bucket storage in S3
	client.region = c.Region

	log.Println("[INFO] Building AWS auth structure")
	creds, err := GetCredentials(c)
	if err != nil {
		return nil, err
	}
	// Call Get to check for credential provider. If nothing found, we'll get an
	// error, and we can present it nicely to the user
	cp, err := creds.Get()
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoCredentialProviders" {
			return nil, errors.New(`No valid credential sources found for AWS Provider.
  Please see https://terraform.io/docs/providers/aws/index.html for more information on
  providing credentials for the AWS Provider`)
		}

		return nil, fmt.Errorf("Error loading credentials for AWS Provider: %s", err)
	}

	log.Printf("[INFO] AWS Auth provider used: %q", cp.ProviderName)

	awsConfig := &aws.Config{
		Credentials:      creds,
		Region:           aws.String(c.Region),
		MaxRetries:       aws.Int(c.MaxRetries),
		HTTPClient:       cleanhttp.DefaultClient(),
		S3ForcePathStyle: aws.Bool(c.S3ForcePathStyle),
	}

	if logging.IsDebugOrHigher() {
		awsConfig.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
		awsConfig.Logger = awsLogger{}
	}

	if c.Insecure {
		transport := awsConfig.HTTPClient.Transport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Set up base session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, errwrap.Wrapf("Error creating AWS session: {{err}}", err)
	}

	// Removes the SDK Version handler, so we only have the provider User-Agent
	// Ex: "User-Agent: APN/1.0 HashiCorp/1.0 Terraform/0.7.9-dev"
	sess.Handlers.Build.Remove(request.NamedHandler{Name: "core.SDKVersionUserAgentHandler"})
	sess.Handlers.Build.PushFrontNamed(addTerraformVersionToUserAgent)

	if extraDebug := os.Getenv("TERRAFORM_AWS_AUTHFAILURE_DEBUG"); extraDebug != "" {
		sess.Handlers.UnmarshalError.PushFrontNamed(debugAuthFailure)
	}

	// Some services have user-configurable endpoints
	awsEc2Sess := sess.Copy(&aws.Config{Endpoint: aws.String(c.Ec2Endpoint)})
	awsElbSess := sess.Copy(&aws.Config{Endpoint: aws.String(c.ElbEndpoint)})
	awsIamSess := sess.Copy(&aws.Config{Endpoint: aws.String(c.IamEndpoint)})
	awsS3Sess := sess.Copy(&aws.Config{Endpoint: aws.String(c.S3Endpoint)})

	// These two services need to be set up early so we can check on AccountID
	client.iamconn = iam.New(awsIamSess)

	partition, accountId, err := GetAccountInfo(client.iamconn, nil, cp.ProviderName)
	if err == nil {
		client.partition = partition
		client.accountid = accountId
	}

	authErr := c.ValidateAccountId(client.accountid)
	if authErr != nil {
		return nil, authErr
	}

	client.ec2conn = ec2.New(awsEc2Sess)
	client.elbconn = elb.New(awsElbSess)
	client.elbv2conn = elbv2.New(awsElbSess)
	client.s3conn = s3.New(awsS3Sess)

	return &client, nil
}

// ValidateRegion returns an error if the configured region is not a
// valid aws region and nil otherwise.
func (c *Config) ValidateRegion() error {
	var regions = []string{
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"cn-north-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-gov-west-1",
		"us-west-1",
		"us-west-2",
	}

	for _, valid := range regions {
		if c.Region == valid {
			return nil
		}
	}
	return fmt.Errorf("Not a valid region: %s", c.Region)
}

// ValidateAccountId returns a context-specific error if the configured account
// id is explicitly forbidden or not authorised; and nil if it is authorised.
func (c *Config) ValidateAccountId(accountId string) error {
	if c.AllowedAccountIds == nil && c.ForbiddenAccountIds == nil {
		return nil
	}

	log.Println("[INFO] Validating account ID")

	if c.ForbiddenAccountIds != nil {
		for _, id := range c.ForbiddenAccountIds {
			if id == accountId {
				return fmt.Errorf("Forbidden account ID (%s)", id)
			}
		}
	}

	if c.AllowedAccountIds != nil {
		for _, id := range c.AllowedAccountIds {
			if id == accountId {
				return nil
			}
		}
		return fmt.Errorf("Account ID not allowed (%s)", accountId)
	}

	return nil
}

// addTerraformVersionToUserAgent is a named handler that will add Terraform's
// version information to requests made by the AWS SDK.
var addTerraformVersionToUserAgent = request.NamedHandler{
	Name: "terraform.TerraformVersionUserAgentHandler",
	Fn: request.MakeAddToUserAgentHandler(
		"APN/1.0 HashiCorp/1.0 terraform-provider-osc/v0.7.1 Terraform", terraform.VersionString()),
}

var debugAuthFailure = request.NamedHandler{
	Name: "terraform.AuthFailureAdditionalDebugHandler",
	Fn: func(req *request.Request) {
		if isAWSErr(req.Error, "AuthFailure", "AWS was not able to validate the provided access credentials") {
			log.Printf("[INFO] Additional AuthFailure Debugging Context")
			log.Printf("[INFO] Current system UTC time: %s", time.Now().UTC())
			log.Printf("[INFO] Request object: %s", spew.Sdump(req))
		}
	},
}

type awsLogger struct{}

func (l awsLogger) Log(args ...interface{}) {
	tokens := make([]string, 0, len(args))
	for _, arg := range args {
		if token, ok := arg.(string); ok {
			tokens = append(tokens, token)
		}
	}
	log.Printf("[DEBUG] [aws-sdk-go] %s", strings.Join(tokens, " "))
}
