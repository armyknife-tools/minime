// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package elbv2 provides the client and types for making API
// requests to Elastic Load Balancing.
//
// A load balancer distributes incoming traffic across targets, such as your
// EC2 instances. This enables you to increase the availability of your application.
// The load balancer also monitors the health of its registered targets and
// ensures that it routes traffic only to healthy targets. You configure your
// load balancer to accept incoming traffic by specifying one or more listeners,
// which are configured with a protocol and port number for connections from
// clients to the load balancer. You configure a target group with a protocol
// and port number for connections from the load balancer to the targets, and
// with health check settings to be used when checking the health status of
// the targets.
//
// Elastic Load Balancing supports the following types of load balancers: Application
// Load Balancers, Network Load Balancers, and Classic Load Balancers.
//
// An Application Load Balancer makes routing and load balancing decisions at
// the application layer (HTTP/HTTPS). A Network Load Balancer makes routing
// and load balancing decisions at the transport layer (TCP/TLS). Both Application
// Load Balancers and Network Load Balancers can route requests to one or more
// ports on each EC2 instance or container instance in your virtual private
// cloud (VPC).
//
// A Classic Load Balancer makes routing and load balancing decisions either
// at the transport layer (TCP/SSL) or the application layer (HTTP/HTTPS), and
// supports either EC2-Classic or a VPC. For more information, see the Elastic
// Load Balancing User Guide (http://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/).
//
// This reference covers the 2015-12-01 API, which supports Application Load
// Balancers and Network Load Balancers. The 2012-06-01 API supports Classic
// Load Balancers.
//
// To get started, complete the following tasks:
//
// Create a load balancer using CreateLoadBalancer.
//
// Create a target group using CreateTargetGroup.
//
// Register targets for the target group using RegisterTargets.
//
// Create one or more listeners for your load balancer using CreateListener.
//
// To delete a load balancer and its related resources, complete the following
// tasks:
//
// Delete the load balancer using DeleteLoadBalancer.
//
// Delete the target group using DeleteTargetGroup.
//
// All Elastic Load Balancing operations are idempotent, which means that they
// complete at most one time. If you repeat an operation, it succeeds.
//
// See https://docs.aws.amazon.com/goto/WebAPI/elasticloadbalancingv2-2015-12-01 for more information on this service.
//
// See elbv2 package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/
//
// Using the Client
//
// To contact Elastic Load Balancing with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the Elastic Load Balancing client ELBV2 for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#New
package elbv2
