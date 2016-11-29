// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

// Amazon EC2 Simple Systems Manager (SSM) enables you to remotely manage the
// configuration of your Amazon EC2 instances, virtual machines (VMs), or servers
// in your on-premises environment or in an environment provided by other cloud
// providers using scripts, commands, or the Amazon EC2 console. SSM includes
// an on-demand solution called Amazon EC2 Run Command and a lightweight instance
// configuration solution called SSM Config.
//
// This references is intended to be used with the EC2 Run Command User Guide
// for Linux (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/execute-remote-commands.html)
// or Windows (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/execute-remote-commands.html).
//
// You must register your on-premises servers and VMs through an activation
// process before you can configure them using Run Command. Registered servers
// and VMs are called managed instances. For more information, see Setting Up
// Run Command On Managed Instances (On-Premises Servers and VMs) on Linux (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/managed-instances.html)
// or Setting Up Run Command On Managed Instances (On-Premises Servers and VMs)
// on Windows (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/managed-instances.html).
//
// Run Command
//
// Run Command provides an on-demand experience for executing commands. You
// can use pre-defined SSM documents to perform the actions listed later in
// this section, or you can create your own documents. With these documents,
// you can remotely configure your instances by sending commands using the Commands
// page in the Amazon EC2 console (http://console.aws.amazon.com/ec2/), AWS
// Tools for Windows PowerShell (http://docs.aws.amazon.com/powershell/latest/reference/items/Amazon_Simple_Systems_Management_cmdlets.html),
// the AWS CLI (http://docs.aws.amazon.com/cli/latest/reference/ssm/index.html),
// or AWS SDKs.
//
// Run Command reports the status of the command execution for each instance
// targeted by a command. You can also audit the command execution to understand
// who executed commands, when, and what changes were made. By switching between
// different SSM documents, you can quickly configure your instances with different
// types of commands. To get started with Run Command, verify that your environment
// meets the prerequisites for remotely running commands on EC2 instances (Linux
// (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/remote-commands-prereq.html)
// or Windows (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/remote-commands-prereq.html)).
//
// SSM Config
//
// SSM Config is a lightweight instance configuration solution. SSM Config is
// currently only available for Windows instances. With SSM Config, you can
// specify a setup configuration for your instances. SSM Config is similar to
// EC2 User Data, which is another way of running one-time scripts or applying
// settings during instance launch. SSM Config is an extension of this capability.
// Using SSM documents, you can specify which actions the system should perform
// on your instances, including which applications to install, which AWS Directory
// Service directory to join, which Microsoft PowerShell modules to install,
// etc. If an instance is missing one or more of these configurations, the system
// makes those changes. By default, the system checks every five minutes to
// see if there is a new configuration to apply as defined in a new SSM document.
// If so, the system updates the instances accordingly. In this way, you can
// remotely maintain a consistent configuration baseline on your instances.
// SSM Config is available using the AWS CLI or the AWS Tools for Windows PowerShell.
// For more information, see Managing Windows Instance Configuration (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/ec2-configuration-manage.html).
//
// SSM Config and Run Command include the following pre-defined documents.
//
// Linux
//
// AWS-RunShellScript to run shell scripts
//
//    * AWS-UpdateSSMAgent to update the Amazon SSM agent
//
// Windows
//
//    * AWS-JoinDirectoryServiceDomain to join an AWS Directory
//
//    * AWS-RunPowerShellScript to run PowerShell commands or scripts
//
//    * AWS-UpdateEC2Config to update the EC2Config service
//
//    * AWS-ConfigureWindowsUpdate to configure Windows Update settings
//
//    * AWS-InstallApplication to install, repair, or uninstall software using
//    an MSI package
//
//    * AWS-InstallPowerShellModule to install PowerShell modules
//
//    * AWS-ConfigureCloudWatch to configure Amazon CloudWatch Logs to monitor
//    applications and systems
//
//    * AWS-ListWindowsInventory to collect information about an EC2 instance
//    running in Windows.
//
//    * AWS-FindWindowsUpdates to scan an instance and determines which updates
//    are missing.
//
//    * AWS-InstallMissingWindowsUpdates to install missing updates on your
//    EC2 instance.
//
//    * AWS-InstallSpecificWindowsUpdates to install one or more specific updates.
//
// The commands or scripts specified in SSM documents run with administrative
//    privilege on your instances because the Amazon SSM agent runs as root
//    on Linux and the EC2Config service runs in the Local System account on
//    Windows. If a user has permission to execute any of the pre-defined SSM
//    documents (any document that begins with AWS-*) then that user also has
//    administrator access to the instance. Delegate access to Run Command and
//    SSM Config judiciously. This becomes extremely important if you create
//    your own SSM documents. Amazon Web Services does not provide guidance
//    about how to create secure SSM documents. You create SSM documents and
//    delegate access to Run Command at your own risk. As a security best practice,
//    we recommend that you assign access to "AWS-*" documents, especially the
//    AWS-RunShellScript document on Linux and the AWS-RunPowerShellScript document
//    on Windows, to trusted administrators only. You can create SSM documents
//    for specific tasks and delegate access to non-administrators.
//
// For information about creating and sharing SSM documents, see the following
//    topics in the SSM User Guide:
//
//    * Creating SSM Documents (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/create-ssm-doc.html)
//    and    * Sharing SSM Documents (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ssm-sharing.html)
//    (Linux)
//
//    * Creating SSM Documents (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/create-ssm-doc.html)
//    and    * Sharing SSM Documents (http://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/ssm-sharing.html)
//    (Windows)
//The service client's operations are safe to be used concurrently.
// It is not safe to mutate any of the client's properties though.
type SSM struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "ssm"

// New creates a new instance of the SSM client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a SSM client from just a session.
//     svc := ssm.New(mySession)
//
//     // Create a SSM client with additional configuration
//     svc := ssm.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) *SSM {
	c := p.ClientConfig(ServiceName, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion string) *SSM {
	svc := &SSM{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2014-11-06",
				JSONVersion:   "1.1",
				TargetPrefix:  "AmazonSSM",
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

// newRequest creates a new request for a SSM operation and runs any
// custom request initialization.
func (c *SSM) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
