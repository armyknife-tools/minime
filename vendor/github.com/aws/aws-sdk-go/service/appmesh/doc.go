// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package appmesh provides the client and types for making API
// requests to AWS App Mesh.
//
// AWS App Mesh is a service mesh based on the Envoy proxy that makes it easy
// to monitor and control containerized microservices. App Mesh standardizes
// how your microservices communicate, giving you end-to-end visibility and
// helping to ensure high-availability for your applications.
//
// App Mesh gives you consistent visibility and network traffic controls for
// every microservice in an application. You can use App Mesh with Amazon ECS
// (using the Amazon EC2 launch type), Amazon EKS, and Kubernetes on AWS.
//
// App Mesh supports containerized microservice applications that use service
// discovery naming for their components. To use App Mesh, you must have a containerized
// application running on Amazon EC2 instances, hosted in either Amazon ECS,
// Amazon EKS, or Kubernetes on AWS. For more information about service discovery
// on Amazon ECS, see Service Discovery (http://docs.aws.amazon.com/AmazonECS/latest/developerguideservice-discovery.html)
// in the Amazon Elastic Container Service Developer Guide. Kubernetes kube-dns
// is supported. For more information, see DNS for Services and Pods (https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/)
// in the Kubernetes documentation.
//
// See https://docs.aws.amazon.com/goto/WebAPI/appmesh-2018-10-01 for more information on this service.
//
// See appmesh package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/appmesh/
//
// Using the Client
//
// To contact AWS App Mesh with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the AWS App Mesh client AppMesh for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/appmesh/#New
package appmesh
