// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package cloudtrail

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go/private/signer/v4"
)

// This is the CloudTrail API Reference. It provides descriptions of actions,
// data types, common parameters, and common errors for CloudTrail.
//
// CloudTrail is a web service that records AWS API calls for your AWS account
// and delivers log files to an Amazon S3 bucket. The recorded information includes
// the identity of the user, the start time of the AWS API call, the source
// IP address, the request parameters, and the response elements returned by
// the service.
//
//  As an alternative to using the API, you can use one of the AWS SDKs, which
// consist of libraries and sample code for various programming languages and
// platforms (Java, Ruby, .NET, iOS, Android, etc.). The SDKs provide a convenient
// way to create programmatic access to AWSCloudTrail. For example, the SDKs
// take care of cryptographically signing requests, managing errors, and retrying
// requests automatically. For information about the AWS SDKs, including how
// to download and install them, see the Tools for Amazon Web Services page
// (http://aws.amazon.com/tools/).  See the CloudTrail User Guide for information
// about the data that is included with each AWS API call listed in the log
// files.
//The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
type CloudTrail struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "cloudtrail"

// New creates a new instance of the CloudTrail client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a CloudTrail client from just a session.
//     svc := cloudtrail.New(mySession)
//
//     // Create a CloudTrail client with additional configuration
//     svc := cloudtrail.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *CloudTrail {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion string) *CloudTrail {
	svc := &CloudTrail{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2013-11-01",
				JSONVersion:   "1.1",
				TargetPrefix:  "com.amazonaws.cloudtrail.v20131101.CloudTrail_20131101",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBack(v4.Sign)
	svc.Handlers.Build.PushBack(jsonrpc.Build)
	svc.Handlers.Unmarshal.PushBack(jsonrpc.Unmarshal)
	svc.Handlers.UnmarshalMeta.PushBack(jsonrpc.UnmarshalMeta)
	svc.Handlers.UnmarshalError.PushBack(jsonrpc.UnmarshalError)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a CloudTrail operation and runs any
// custom request initialization.
func (c *CloudTrail) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
