// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package applicationautoscaling provides the client and types for making API
// requests to Application Auto Scaling.
//
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
//    * AppStream 2.0 fleets. For more information, see Fleet Auto Scaling for
//    Amazon AppStream 2.0 (http://docs.aws.amazon.com/appstream2/latest/developerguide/autoscaling.html)
//    in the Amazon AppStream 2.0 Developer Guide.
//
//    * Provisioned read and write capacity for Amazon DynamoDB tables and global
//    secondary indexes. For more information, see Auto Scaling for DynamoDB
//    (http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/TargetTracking.html)
//    in the Amazon DynamoDB Developer Guide.
//
// For a list of supported regions, see AWS Regions and Endpoints: Application
// Auto Scaling (http://docs.aws.amazon.com/general/latest/gr/rande.html#as-app_region)
// in the AWS General Reference.
//
// See https://docs.aws.amazon.com/goto/WebAPI/application-autoscaling-2016-02-06 for more information on this service.
//
// See applicationautoscaling package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/applicationautoscaling/
//
// Using the Client
//
// To Application Auto Scaling with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the Application Auto Scaling client ApplicationAutoScaling for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/applicationautoscaling/#New
package applicationautoscaling
