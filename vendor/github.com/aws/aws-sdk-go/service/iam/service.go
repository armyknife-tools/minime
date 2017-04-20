// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/query"
)

// AWS Identity and Access Management (IAM) is a web service that you can use
// to manage users and user permissions under your AWS account. This guide provides
// descriptions of IAM actions that you can call programmatically. For general
// information about IAM, see AWS Identity and Access Management (IAM) (http://aws.amazon.com/iam/).
// For the user guide for IAM, see Using IAM (http://docs.aws.amazon.com/IAM/latest/UserGuide/).
//
// AWS provides SDKs that consist of libraries and sample code for various programming
// languages and platforms (Java, Ruby, .NET, iOS, Android, etc.). The SDKs
// provide a convenient way to create programmatic access to IAM and AWS. For
// example, the SDKs take care of tasks such as cryptographically signing requests
// (see below), managing errors, and retrying requests automatically. For information
// about the AWS SDKs, including how to download and install them, see the Tools
// for Amazon Web Services (http://aws.amazon.com/tools/) page.
//
// We recommend that you use the AWS SDKs to make programmatic API calls to
// IAM. However, you can also use the IAM Query API to make direct calls to
// the IAM web service. To learn more about the IAM Query API, see Making Query
// Requests (http://docs.aws.amazon.com/IAM/latest/UserGuide/IAM_UsingQueryAPI.html)
// in the Using IAM guide. IAM supports GET and POST requests for all actions.
// That is, the API does not require you to use GET for some actions and POST
// for others. However, GET requests are subject to the limitation size of a
// URL. Therefore, for operations that require larger sizes, use a POST request.
//
// Signing Requests
//
// Requests must be signed using an access key ID and a secret access key. We
// strongly recommend that you do not use your AWS account access key ID and
// secret access key for everyday work with IAM. You can use the access key
// ID and secret access key for an IAM user or you can use the AWS Security
// Token Service to generate temporary security credentials and use those to
// sign requests.
//
// To sign requests, we recommend that you use Signature Version 4 (http://docs.aws.amazon.com/general/latest/gr/signature-version-4.html).
// If you have an existing application that uses Signature Version 2, you do
// not have to update it to use Signature Version 4. However, some operations
// now require Signature Version 4. The documentation for operations that require
// version 4 indicate this requirement.
//
// Additional Resources
//
// For more information, see the following:
//
//    * AWS Security Credentials (http://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html).
//    This topic provides general information about the types of credentials
//    used for accessing AWS.
//
//    * IAM Best Practices (http://docs.aws.amazon.com/IAM/latest/UserGuide/IAMBestPractices.html).
//    This topic presents a list of suggestions for using the IAM service to
//    help secure your AWS resources.
//
//    * Signing AWS API Requests (http://docs.aws.amazon.com/general/latest/gr/signing_aws_api_requests.html).
//    This set of topics walk you through the process of signing a request using
//    an access key ID and secret access key.
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/iam-2010-05-08
type IAM struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "iam"       // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the IAM client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a IAM client from just a session.
//     svc := iam.New(mySession)
//
//     // Create a IAM client with additional configuration
//     svc := iam.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *IAM {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *IAM {
	svc := &IAM{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2010-05-08",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(query.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a IAM operation and runs any
// custom request initialization.
func (c *IAM) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
