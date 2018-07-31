// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package eks

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
)

const opCreateCluster = "CreateCluster"

// CreateClusterRequest generates a "aws/request.Request" representing the
// client's request for the CreateCluster operation. The "output" return
// value will be populated with the request's response once the request completes
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See CreateCluster for more information on using the CreateCluster
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the CreateClusterRequest method.
//    req, resp := client.CreateClusterRequest(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/CreateCluster
func (c *EKS) CreateClusterRequest(input *CreateClusterInput) (req *request.Request, output *CreateClusterOutput) {
	op := &request.Operation{
		Name:       opCreateCluster,
		HTTPMethod: "POST",
		HTTPPath:   "/clusters",
	}

	if input == nil {
		input = &CreateClusterInput{}
	}

	output = &CreateClusterOutput{}
	req = c.newRequest(op, input, output)
	return
}

// CreateCluster API operation for Amazon Elastic Container Service for Kubernetes.
//
// Creates an Amazon EKS control plane.
//
// The Amazon EKS control plane consists of control plane instances that run
// the Kubernetes software, like etcd and the API server. The control plane
// runs in an account managed by AWS, and the Kubernetes API is exposed via
// the Amazon EKS API server endpoint.
//
// Amazon EKS worker nodes run in your AWS account and connect to your cluster's
// control plane via the Kubernetes API server endpoint and a certificate file
// that is created for your cluster.
//
// The cluster control plane is provisioned across multiple Availability Zones
// and fronted by an Elastic Load Balancing Network Load Balancer. Amazon EKS
// also provisions elastic network interfaces in your VPC subnets to provide
// connectivity from the control plane instances to the worker nodes (for example,
// to support kubectl exec, logs, and proxy data flows).
//
// After you create an Amazon EKS cluster, you must configure your Kubernetes
// tooling to communicate with the API server and launch worker nodes into your
// cluster. For more information, see Managing Cluster Authentication (http://docs.aws.amazon.com/eks/latest/userguide/managing-auth.html)
// and Launching Amazon EKS Worker Nodes (http://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html)in
// the Amazon EKS User Guide.
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for Amazon Elastic Container Service for Kubernetes's
// API operation CreateCluster for usage and error information.
//
// Returned Error Codes:
//   * ErrCodeResourceInUseException "ResourceInUseException"
//   The specified resource is in use.
//
//   * ErrCodeResourceLimitExceededException "ResourceLimitExceededException"
//   You have encountered a service limit on the specified resource.
//
//   * ErrCodeInvalidParameterException "InvalidParameterException"
//   The specified parameter is invalid. Review the available parameters for the
//   API request.
//
//   * ErrCodeClientException "ClientException"
//   These errors are usually caused by a client action, such as using an action
//   or resource on behalf of a user that doesn't have permissions to use the
//   action or resource, or specifying an identifier that is not valid.
//
//   * ErrCodeServerException "ServerException"
//   These errors are usually caused by a server-side issue.
//
//   * ErrCodeServiceUnavailableException "ServiceUnavailableException"
//   The service is unavailable, back off and retry the operation.
//
//   * ErrCodeUnsupportedAvailabilityZoneException "UnsupportedAvailabilityZoneException"
//   At least one of your specified cluster subnets is in an Availability Zone
//   that does not support Amazon EKS. The exception output will specify the supported
//   Availability Zones for your account, from which you can choose subnets for
//   your cluster.
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/CreateCluster
func (c *EKS) CreateCluster(input *CreateClusterInput) (*CreateClusterOutput, error) {
	req, out := c.CreateClusterRequest(input)
	return out, req.Send()
}

// CreateClusterWithContext is the same as CreateCluster with the addition of
// the ability to pass a context and additional request options.
//
// See CreateCluster for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *EKS) CreateClusterWithContext(ctx aws.Context, input *CreateClusterInput, opts ...request.Option) (*CreateClusterOutput, error) {
	req, out := c.CreateClusterRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opDeleteCluster = "DeleteCluster"

// DeleteClusterRequest generates a "aws/request.Request" representing the
// client's request for the DeleteCluster operation. The "output" return
// value will be populated with the request's response once the request completes
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See DeleteCluster for more information on using the DeleteCluster
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the DeleteClusterRequest method.
//    req, resp := client.DeleteClusterRequest(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/DeleteCluster
func (c *EKS) DeleteClusterRequest(input *DeleteClusterInput) (req *request.Request, output *DeleteClusterOutput) {
	op := &request.Operation{
		Name:       opDeleteCluster,
		HTTPMethod: "DELETE",
		HTTPPath:   "/clusters/{name}",
	}

	if input == nil {
		input = &DeleteClusterInput{}
	}

	output = &DeleteClusterOutput{}
	req = c.newRequest(op, input, output)
	return
}

// DeleteCluster API operation for Amazon Elastic Container Service for Kubernetes.
//
// Deletes the Amazon EKS cluster control plane.
//
// If you have active services in your cluster that are associated with a load
// balancer, you must delete those services before deleting the cluster so that
// the load balancers are deleted properly. Otherwise, you can have orphaned
// resources in your VPC that prevent you from being able to delete the VPC.
// For more information, see Deleting a Cluster (http://docs.aws.amazon.com/eks/latest/userguide/delete-cluster.html)
// in the Amazon EKS User Guide.
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for Amazon Elastic Container Service for Kubernetes's
// API operation DeleteCluster for usage and error information.
//
// Returned Error Codes:
//   * ErrCodeResourceInUseException "ResourceInUseException"
//   The specified resource is in use.
//
//   * ErrCodeResourceNotFoundException "ResourceNotFoundException"
//   The specified resource could not be found. You can view your available clusters
//   with ListClusters. Amazon EKS clusters are region-specific.
//
//   * ErrCodeClientException "ClientException"
//   These errors are usually caused by a client action, such as using an action
//   or resource on behalf of a user that doesn't have permissions to use the
//   action or resource, or specifying an identifier that is not valid.
//
//   * ErrCodeServerException "ServerException"
//   These errors are usually caused by a server-side issue.
//
//   * ErrCodeServiceUnavailableException "ServiceUnavailableException"
//   The service is unavailable, back off and retry the operation.
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/DeleteCluster
func (c *EKS) DeleteCluster(input *DeleteClusterInput) (*DeleteClusterOutput, error) {
	req, out := c.DeleteClusterRequest(input)
	return out, req.Send()
}

// DeleteClusterWithContext is the same as DeleteCluster with the addition of
// the ability to pass a context and additional request options.
//
// See DeleteCluster for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *EKS) DeleteClusterWithContext(ctx aws.Context, input *DeleteClusterInput, opts ...request.Option) (*DeleteClusterOutput, error) {
	req, out := c.DeleteClusterRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opDescribeCluster = "DescribeCluster"

// DescribeClusterRequest generates a "aws/request.Request" representing the
// client's request for the DescribeCluster operation. The "output" return
// value will be populated with the request's response once the request completes
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See DescribeCluster for more information on using the DescribeCluster
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the DescribeClusterRequest method.
//    req, resp := client.DescribeClusterRequest(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/DescribeCluster
func (c *EKS) DescribeClusterRequest(input *DescribeClusterInput) (req *request.Request, output *DescribeClusterOutput) {
	op := &request.Operation{
		Name:       opDescribeCluster,
		HTTPMethod: "GET",
		HTTPPath:   "/clusters/{name}",
	}

	if input == nil {
		input = &DescribeClusterInput{}
	}

	output = &DescribeClusterOutput{}
	req = c.newRequest(op, input, output)
	return
}

// DescribeCluster API operation for Amazon Elastic Container Service for Kubernetes.
//
// Returns descriptive information about an Amazon EKS cluster.
//
// The API server endpoint and certificate authority data returned by this operation
// are required for kubelet and kubectl to communicate with your Kubernetes
// API server. For more information, see Create a kubeconfig for Amazon EKS
// (http://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html).
//
// The API server endpoint and certificate authority data are not available
// until the cluster reaches the ACTIVE state.
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for Amazon Elastic Container Service for Kubernetes's
// API operation DescribeCluster for usage and error information.
//
// Returned Error Codes:
//   * ErrCodeResourceNotFoundException "ResourceNotFoundException"
//   The specified resource could not be found. You can view your available clusters
//   with ListClusters. Amazon EKS clusters are region-specific.
//
//   * ErrCodeClientException "ClientException"
//   These errors are usually caused by a client action, such as using an action
//   or resource on behalf of a user that doesn't have permissions to use the
//   action or resource, or specifying an identifier that is not valid.
//
//   * ErrCodeServerException "ServerException"
//   These errors are usually caused by a server-side issue.
//
//   * ErrCodeServiceUnavailableException "ServiceUnavailableException"
//   The service is unavailable, back off and retry the operation.
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/DescribeCluster
func (c *EKS) DescribeCluster(input *DescribeClusterInput) (*DescribeClusterOutput, error) {
	req, out := c.DescribeClusterRequest(input)
	return out, req.Send()
}

// DescribeClusterWithContext is the same as DescribeCluster with the addition of
// the ability to pass a context and additional request options.
//
// See DescribeCluster for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *EKS) DescribeClusterWithContext(ctx aws.Context, input *DescribeClusterInput, opts ...request.Option) (*DescribeClusterOutput, error) {
	req, out := c.DescribeClusterRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

const opListClusters = "ListClusters"

// ListClustersRequest generates a "aws/request.Request" representing the
// client's request for the ListClusters operation. The "output" return
// value will be populated with the request's response once the request completes
// successfuly.
//
// Use "Send" method on the returned Request to send the API call to the service.
// the "output" return value is not valid until after Send returns without error.
//
// See ListClusters for more information on using the ListClusters
// API call, and error handling.
//
// This method is useful when you want to inject custom logic or configuration
// into the SDK's request lifecycle. Such as custom headers, or retry logic.
//
//
//    // Example sending a request using the ListClustersRequest method.
//    req, resp := client.ListClustersRequest(params)
//
//    err := req.Send()
//    if err == nil { // resp is now filled
//        fmt.Println(resp)
//    }
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/ListClusters
func (c *EKS) ListClustersRequest(input *ListClustersInput) (req *request.Request, output *ListClustersOutput) {
	op := &request.Operation{
		Name:       opListClusters,
		HTTPMethod: "GET",
		HTTPPath:   "/clusters",
	}

	if input == nil {
		input = &ListClustersInput{}
	}

	output = &ListClustersOutput{}
	req = c.newRequest(op, input, output)
	return
}

// ListClusters API operation for Amazon Elastic Container Service for Kubernetes.
//
// Lists the Amazon EKS clusters in your AWS account in the specified region.
//
// Returns awserr.Error for service API and SDK errors. Use runtime type assertions
// with awserr.Error's Code and Message methods to get detailed information about
// the error.
//
// See the AWS API reference guide for Amazon Elastic Container Service for Kubernetes's
// API operation ListClusters for usage and error information.
//
// Returned Error Codes:
//   * ErrCodeInvalidParameterException "InvalidParameterException"
//   The specified parameter is invalid. Review the available parameters for the
//   API request.
//
//   * ErrCodeClientException "ClientException"
//   These errors are usually caused by a client action, such as using an action
//   or resource on behalf of a user that doesn't have permissions to use the
//   action or resource, or specifying an identifier that is not valid.
//
//   * ErrCodeServerException "ServerException"
//   These errors are usually caused by a server-side issue.
//
//   * ErrCodeServiceUnavailableException "ServiceUnavailableException"
//   The service is unavailable, back off and retry the operation.
//
// See also, https://docs.aws.amazon.com/goto/WebAPI/eks-2017-11-01/ListClusters
func (c *EKS) ListClusters(input *ListClustersInput) (*ListClustersOutput, error) {
	req, out := c.ListClustersRequest(input)
	return out, req.Send()
}

// ListClustersWithContext is the same as ListClusters with the addition of
// the ability to pass a context and additional request options.
//
// See ListClusters for details on how to use this API operation.
//
// The context must be non-nil and will be used for request cancellation. If
// the context is nil a panic will occur. In the future the SDK may create
// sub-contexts for http.Requests. See https://golang.org/pkg/context/
// for more information on using Contexts.
func (c *EKS) ListClustersWithContext(ctx aws.Context, input *ListClustersInput, opts ...request.Option) (*ListClustersOutput, error) {
	req, out := c.ListClustersRequest(input)
	req.SetContext(ctx)
	req.ApplyOptions(opts...)
	return out, req.Send()
}

// An object representing the certificate-authority-data for your cluster.
type Certificate struct {
	_ struct{} `type:"structure"`

	// The base64 encoded certificate data required to communicate with your cluster.
	// Add this to the certificate-authority-data section of the kubeconfig file
	// for your cluster.
	Data *string `locationName:"data" type:"string"`
}

// String returns the string representation
func (s Certificate) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Certificate) GoString() string {
	return s.String()
}

// SetData sets the Data field's value.
func (s *Certificate) SetData(v string) *Certificate {
	s.Data = &v
	return s
}

// An object representing an Amazon EKS cluster.
type Cluster struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) of the cluster.
	Arn *string `locationName:"arn" type:"string"`

	// The certificate-authority-data for your cluster.
	CertificateAuthority *Certificate `locationName:"certificateAuthority" type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request.
	ClientRequestToken *string `locationName:"clientRequestToken" type:"string"`

	// The Unix epoch time stamp in seconds for when the cluster was created.
	CreatedAt *time.Time `locationName:"createdAt" type:"timestamp"`

	// The endpoint for your Kubernetes API server.
	Endpoint *string `locationName:"endpoint" type:"string"`

	// The name of the cluster.
	Name *string `locationName:"name" type:"string"`

	// The VPC subnets and security groups used by the cluster control plane. Amazon
	// EKS VPC resources have specific requirements to work properly with Kubernetes.
	// For more information, see Cluster VPC Considerations (http://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html)
	// and Cluster Security Group Considerations (http://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html)
	// in the Amazon EKS User Guide.
	ResourcesVpcConfig *VpcConfigResponse `locationName:"resourcesVpcConfig" type:"structure"`

	// The Amazon Resource Name (ARN) of the IAM role that provides permissions
	// for the Kubernetes control plane to make calls to AWS API operations on your
	// behalf.
	RoleArn *string `locationName:"roleArn" type:"string"`

	// The current status of the cluster.
	Status *string `locationName:"status" type:"string" enum:"ClusterStatus"`

	// The Kubernetes server version for the cluster.
	Version *string `locationName:"version" type:"string"`
}

// String returns the string representation
func (s Cluster) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Cluster) GoString() string {
	return s.String()
}

// SetArn sets the Arn field's value.
func (s *Cluster) SetArn(v string) *Cluster {
	s.Arn = &v
	return s
}

// SetCertificateAuthority sets the CertificateAuthority field's value.
func (s *Cluster) SetCertificateAuthority(v *Certificate) *Cluster {
	s.CertificateAuthority = v
	return s
}

// SetClientRequestToken sets the ClientRequestToken field's value.
func (s *Cluster) SetClientRequestToken(v string) *Cluster {
	s.ClientRequestToken = &v
	return s
}

// SetCreatedAt sets the CreatedAt field's value.
func (s *Cluster) SetCreatedAt(v time.Time) *Cluster {
	s.CreatedAt = &v
	return s
}

// SetEndpoint sets the Endpoint field's value.
func (s *Cluster) SetEndpoint(v string) *Cluster {
	s.Endpoint = &v
	return s
}

// SetName sets the Name field's value.
func (s *Cluster) SetName(v string) *Cluster {
	s.Name = &v
	return s
}

// SetResourcesVpcConfig sets the ResourcesVpcConfig field's value.
func (s *Cluster) SetResourcesVpcConfig(v *VpcConfigResponse) *Cluster {
	s.ResourcesVpcConfig = v
	return s
}

// SetRoleArn sets the RoleArn field's value.
func (s *Cluster) SetRoleArn(v string) *Cluster {
	s.RoleArn = &v
	return s
}

// SetStatus sets the Status field's value.
func (s *Cluster) SetStatus(v string) *Cluster {
	s.Status = &v
	return s
}

// SetVersion sets the Version field's value.
func (s *Cluster) SetVersion(v string) *Cluster {
	s.Version = &v
	return s
}

type CreateClusterInput struct {
	_ struct{} `type:"structure"`

	// Unique, case-sensitive identifier you provide to ensure the idempotency of
	// the request.
	ClientRequestToken *string `locationName:"clientRequestToken" type:"string" idempotencyToken:"true"`

	// The unique name to give to your cluster.
	//
	// Name is a required field
	Name *string `locationName:"name" min:"1" type:"string" required:"true"`

	// The VPC subnets and security groups used by the cluster control plane. Amazon
	// EKS VPC resources have specific requirements to work properly with Kubernetes.
	// For more information, see Cluster VPC Considerations (http://docs.aws.amazon.com/eks/latest/userguide/network_reqs.html)
	// and Cluster Security Group Considerations (http://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html)
	// in the Amazon EKS User Guide.
	//
	// ResourcesVpcConfig is a required field
	ResourcesVpcConfig *VpcConfigRequest `locationName:"resourcesVpcConfig" type:"structure" required:"true"`

	// The Amazon Resource Name (ARN) of the IAM role that provides permissions
	// for Amazon EKS to make calls to other AWS API operations on your behalf.
	// For more information, see Amazon EKS Service IAM Role (http://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html)
	// in the Amazon EKS User Guide
	//
	// RoleArn is a required field
	RoleArn *string `locationName:"roleArn" type:"string" required:"true"`

	// The desired Kubernetes version for your cluster. If you do not specify a
	// value here, the latest version available in Amazon EKS is used.
	Version *string `locationName:"version" type:"string"`
}

// String returns the string representation
func (s CreateClusterInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateClusterInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateClusterInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "CreateClusterInput"}
	if s.Name == nil {
		invalidParams.Add(request.NewErrParamRequired("Name"))
	}
	if s.Name != nil && len(*s.Name) < 1 {
		invalidParams.Add(request.NewErrParamMinLen("Name", 1))
	}
	if s.ResourcesVpcConfig == nil {
		invalidParams.Add(request.NewErrParamRequired("ResourcesVpcConfig"))
	}
	if s.RoleArn == nil {
		invalidParams.Add(request.NewErrParamRequired("RoleArn"))
	}
	if s.ResourcesVpcConfig != nil {
		if err := s.ResourcesVpcConfig.Validate(); err != nil {
			invalidParams.AddNested("ResourcesVpcConfig", err.(request.ErrInvalidParams))
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetClientRequestToken sets the ClientRequestToken field's value.
func (s *CreateClusterInput) SetClientRequestToken(v string) *CreateClusterInput {
	s.ClientRequestToken = &v
	return s
}

// SetName sets the Name field's value.
func (s *CreateClusterInput) SetName(v string) *CreateClusterInput {
	s.Name = &v
	return s
}

// SetResourcesVpcConfig sets the ResourcesVpcConfig field's value.
func (s *CreateClusterInput) SetResourcesVpcConfig(v *VpcConfigRequest) *CreateClusterInput {
	s.ResourcesVpcConfig = v
	return s
}

// SetRoleArn sets the RoleArn field's value.
func (s *CreateClusterInput) SetRoleArn(v string) *CreateClusterInput {
	s.RoleArn = &v
	return s
}

// SetVersion sets the Version field's value.
func (s *CreateClusterInput) SetVersion(v string) *CreateClusterInput {
	s.Version = &v
	return s
}

type CreateClusterOutput struct {
	_ struct{} `type:"structure"`

	// The full description of your new cluster.
	Cluster *Cluster `locationName:"cluster" type:"structure"`
}

// String returns the string representation
func (s CreateClusterOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s CreateClusterOutput) GoString() string {
	return s.String()
}

// SetCluster sets the Cluster field's value.
func (s *CreateClusterOutput) SetCluster(v *Cluster) *CreateClusterOutput {
	s.Cluster = v
	return s
}

type DeleteClusterInput struct {
	_ struct{} `type:"structure"`

	// The name of the cluster to delete.
	//
	// Name is a required field
	Name *string `location:"uri" locationName:"name" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteClusterInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteClusterInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteClusterInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DeleteClusterInput"}
	if s.Name == nil {
		invalidParams.Add(request.NewErrParamRequired("Name"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetName sets the Name field's value.
func (s *DeleteClusterInput) SetName(v string) *DeleteClusterInput {
	s.Name = &v
	return s
}

type DeleteClusterOutput struct {
	_ struct{} `type:"structure"`

	// The full description of the cluster to delete.
	Cluster *Cluster `locationName:"cluster" type:"structure"`
}

// String returns the string representation
func (s DeleteClusterOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteClusterOutput) GoString() string {
	return s.String()
}

// SetCluster sets the Cluster field's value.
func (s *DeleteClusterOutput) SetCluster(v *Cluster) *DeleteClusterOutput {
	s.Cluster = v
	return s
}

type DescribeClusterInput struct {
	_ struct{} `type:"structure"`

	// The name of the cluster to describe.
	//
	// Name is a required field
	Name *string `location:"uri" locationName:"name" type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeClusterInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeClusterInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeClusterInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "DescribeClusterInput"}
	if s.Name == nil {
		invalidParams.Add(request.NewErrParamRequired("Name"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetName sets the Name field's value.
func (s *DescribeClusterInput) SetName(v string) *DescribeClusterInput {
	s.Name = &v
	return s
}

type DescribeClusterOutput struct {
	_ struct{} `type:"structure"`

	// The full description of your specified cluster.
	Cluster *Cluster `locationName:"cluster" type:"structure"`
}

// String returns the string representation
func (s DescribeClusterOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeClusterOutput) GoString() string {
	return s.String()
}

// SetCluster sets the Cluster field's value.
func (s *DescribeClusterOutput) SetCluster(v *Cluster) *DescribeClusterOutput {
	s.Cluster = v
	return s
}

type ListClustersInput struct {
	_ struct{} `type:"structure"`

	// The maximum number of cluster results returned by ListClusters in paginated
	// output. When this parameter is used, ListClusters only returns maxResults
	// results in a single page along with a nextToken response element. The remaining
	// results of the initial request can be seen by sending another ListClusters
	// request with the returned nextToken value. This value can be between 1 and
	// 100. If this parameter is not used, then ListClusters returns up to 100 results
	// and a nextToken value if applicable.
	MaxResults *int64 `location:"querystring" locationName:"maxResults" min:"1" type:"integer"`

	// The nextToken value returned from a previous paginated ListClusters request
	// where maxResults was used and the results exceeded the value of that parameter.
	// Pagination continues from the end of the previous results that returned the
	// nextToken value.
	//
	// This token should be treated as an opaque identifier that is only used to
	// retrieve the next items in a list and not for other programmatic purposes.
	NextToken *string `location:"querystring" locationName:"nextToken" type:"string"`
}

// String returns the string representation
func (s ListClustersInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListClustersInput) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *ListClustersInput) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "ListClustersInput"}
	if s.MaxResults != nil && *s.MaxResults < 1 {
		invalidParams.Add(request.NewErrParamMinValue("MaxResults", 1))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetMaxResults sets the MaxResults field's value.
func (s *ListClustersInput) SetMaxResults(v int64) *ListClustersInput {
	s.MaxResults = &v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *ListClustersInput) SetNextToken(v string) *ListClustersInput {
	s.NextToken = &v
	return s
}

type ListClustersOutput struct {
	_ struct{} `type:"structure"`

	// A list of all of the clusters for your account in the specified region.
	Clusters []*string `locationName:"clusters" type:"list"`

	// The nextToken value to include in a future ListClusters request. When the
	// results of a ListClusters request exceed maxResults, this value can be used
	// to retrieve the next page of results. This value is null when there are no
	// more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`
}

// String returns the string representation
func (s ListClustersOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListClustersOutput) GoString() string {
	return s.String()
}

// SetClusters sets the Clusters field's value.
func (s *ListClustersOutput) SetClusters(v []*string) *ListClustersOutput {
	s.Clusters = v
	return s
}

// SetNextToken sets the NextToken field's value.
func (s *ListClustersOutput) SetNextToken(v string) *ListClustersOutput {
	s.NextToken = &v
	return s
}

// An object representing an Amazon EKS cluster VPC configuration request.
type VpcConfigRequest struct {
	_ struct{} `type:"structure"`

	// Specify one or more security groups for the cross-account elastic network
	// interfaces that Amazon EKS creates to use to allow communication between
	// your worker nodes and the Kubernetes control plane.
	SecurityGroupIds []*string `locationName:"securityGroupIds" type:"list"`

	// Specify subnets for your Amazon EKS worker nodes. Amazon EKS creates cross-account
	// elastic network interfaces in these subnets to allow communication between
	// your worker nodes and the Kubernetes control plane.
	//
	// SubnetIds is a required field
	SubnetIds []*string `locationName:"subnetIds" type:"list" required:"true"`
}

// String returns the string representation
func (s VpcConfigRequest) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VpcConfigRequest) GoString() string {
	return s.String()
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *VpcConfigRequest) Validate() error {
	invalidParams := request.ErrInvalidParams{Context: "VpcConfigRequest"}
	if s.SubnetIds == nil {
		invalidParams.Add(request.NewErrParamRequired("SubnetIds"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// SetSecurityGroupIds sets the SecurityGroupIds field's value.
func (s *VpcConfigRequest) SetSecurityGroupIds(v []*string) *VpcConfigRequest {
	s.SecurityGroupIds = v
	return s
}

// SetSubnetIds sets the SubnetIds field's value.
func (s *VpcConfigRequest) SetSubnetIds(v []*string) *VpcConfigRequest {
	s.SubnetIds = v
	return s
}

// An object representing an Amazon EKS cluster VPC configuration response.
type VpcConfigResponse struct {
	_ struct{} `type:"structure"`

	// The security groups associated with the cross-account elastic network interfaces
	// that are used to allow communication between your worker nodes and the Kubernetes
	// control plane.
	SecurityGroupIds []*string `locationName:"securityGroupIds" type:"list"`

	// The subnets associated with your cluster.
	SubnetIds []*string `locationName:"subnetIds" type:"list"`

	// The VPC associated with your cluster.
	VpcId *string `locationName:"vpcId" type:"string"`
}

// String returns the string representation
func (s VpcConfigResponse) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s VpcConfigResponse) GoString() string {
	return s.String()
}

// SetSecurityGroupIds sets the SecurityGroupIds field's value.
func (s *VpcConfigResponse) SetSecurityGroupIds(v []*string) *VpcConfigResponse {
	s.SecurityGroupIds = v
	return s
}

// SetSubnetIds sets the SubnetIds field's value.
func (s *VpcConfigResponse) SetSubnetIds(v []*string) *VpcConfigResponse {
	s.SubnetIds = v
	return s
}

// SetVpcId sets the VpcId field's value.
func (s *VpcConfigResponse) SetVpcId(v string) *VpcConfigResponse {
	s.VpcId = &v
	return s
}

const (
	// ClusterStatusCreating is a ClusterStatus enum value
	ClusterStatusCreating = "CREATING"

	// ClusterStatusActive is a ClusterStatus enum value
	ClusterStatusActive = "ACTIVE"

	// ClusterStatusDeleting is a ClusterStatus enum value
	ClusterStatusDeleting = "DELETING"

	// ClusterStatusFailed is a ClusterStatus enum value
	ClusterStatusFailed = "FAILED"
)
