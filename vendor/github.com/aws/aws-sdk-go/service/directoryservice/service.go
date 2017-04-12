// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package directoryservice

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// AWS Directory Service is a web service that makes it easy for you to setup
// and run directories in the AWS cloud, or connect your AWS resources with
// an existing on-premises Microsoft Active Directory. This guide provides detailed
// information about AWS Directory Service operations, data types, parameters,
// and errors. For information about AWS Directory Services features, see AWS
// Directory Service (https://aws.amazon.com/directoryservice/) and the AWS
// Directory Service Administration Guide (http://docs.aws.amazon.com/directoryservice/latest/admin-guide/what_is.html).
//
// AWS provides SDKs that consist of libraries and sample code for various programming
// languages and platforms (Java, Ruby, .Net, iOS, Android, etc.). The SDKs
// provide a convenient way to create programmatic access to AWS Directory Service
// and other AWS services. For more information about the AWS SDKs, including
// how to download and install them, see Tools for Amazon Web Services (http://aws.amazon.com/tools/).
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ds-2015-04-16
type DirectoryService struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "ds"        // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the DirectoryService client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a DirectoryService client from just a session.
//     svc := directoryservice.New(mySession)
//
//     // Create a DirectoryService client with additional configuration
//     svc := directoryservice.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *DirectoryService {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *DirectoryService {
	svc := &DirectoryService{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2015-04-16",
				JSONVersion:   "1.1",
				TargetPrefix:  "DirectoryService_20150416",
			},
			handlers,
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	svc.Handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc
}

// newRequest creates a new request for a DirectoryService operation and runs any
// custom request initialization.
func (c *DirectoryService) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
