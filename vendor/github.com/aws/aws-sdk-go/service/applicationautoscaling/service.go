// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package applicationautoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// With Application Auto Scaling, you can automatically scale your AWS resources.
// The experience similar to that of Auto Scaling (https://aws.amazon.com/autoscaling/).
// You can use Application Auto Scaling to accomplish the following tasks:
//
//    * Define scaling policies to automatically scale your AWS resources
//
//    * Scale your resources in response to CloudWatch alarms
//
//    * View the history of your scaling events
//
// Application Auto Scaling can scale the following AWS resources:
//
//    * Amazon ECS services. For more information, see Service Auto Scaling
//    (http://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-auto-scaling.html)
//    in the Amazon EC2 Container Service Developer Guide.
//
//    * Amazon EC2 Spot fleets. For more information, see Automatic Scaling
//    for Spot Fleet (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/fleet-auto-scaling.html)
//    in the Amazon EC2 User Guide.
//
//    * Amazon EMR clusters. For more information, see Using Automatic Scaling
//    in Amazon EMR (http://docs.aws.amazon.com/ElasticMapReduce/latest/ManagementGuide/emr-automatic-scaling.html)
//    in the Amazon EMR Management Guide.
//
//    * AppStream 2.0 fleets. For more information, see Autoscaling Amazon AppStream
//    2.0 Resources (http://docs.aws.amazon.com/appstream2/latest/developerguide/autoscaling.html)
//    in the Amazon AppStream 2.0 Developer Guide.
//
// For a list of supported regions, see AWS Regions and Endpoints: Application
// Auto Scaling (http://docs.aws.amazon.com/general/latest/gr/rande.html#as-app_region)
// in the AWS General Reference.
// The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/application-autoscaling-2016-02-06
type ApplicationAutoScaling struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "autoscaling"             // Service endpoint prefix API calls made to.
	EndpointsID = "application-autoscaling" // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the ApplicationAutoScaling client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a ApplicationAutoScaling client from just a session.
//     svc := applicationautoscaling.New(mySession)
//
//     // Create a ApplicationAutoScaling client with additional configuration
//     svc := applicationautoscaling.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *ApplicationAutoScaling {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) *ApplicationAutoScaling {
	if len(signingName) == 0 {
		signingName = "application-autoscaling"
	}
	svc := &ApplicationAutoScaling{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2016-02-06",
				JSONVersion:   "1.1",
				TargetPrefix:  "AnyScaleFrontendService",
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

// newRequest creates a new request for a ApplicationAutoScaling operation and runs any
// custom request initialization.
func (c *ApplicationAutoScaling) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
