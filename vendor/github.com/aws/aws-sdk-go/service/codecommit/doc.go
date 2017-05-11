// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

// Package codecommit provides the client and types for making API
// requests to AWS CodeCommit.
//
// This is the AWS CodeCommit API Reference. This reference provides descriptions
// of the operations and data types for AWS CodeCommit API along with usage
// examples.
//
// You can use the AWS CodeCommit API to work with the following objects:
//
// Repositories, by calling the following:
//
//    * BatchGetRepositories, which returns information about one or more repositories
//    associated with your AWS account
//
//    * CreateRepository, which creates an AWS CodeCommit repository
//
//    * DeleteRepository, which deletes an AWS CodeCommit repository
//
//    * GetRepository, which returns information about a specified repository
//
//    * ListRepositories, which lists all AWS CodeCommit repositories associated
//    with your AWS account
//
//    * UpdateRepositoryDescription, which sets or updates the description of
//    the repository
//
//    * UpdateRepositoryName, which changes the name of the repository. If you
//    change the name of a repository, no other users of that repository will
//    be able to access it until you send them the new HTTPS or SSH URL to use.
//
// Branches, by calling the following:
//
//    * CreateBranch, which creates a new branch in a specified repository
//
//    * GetBranch, which returns information about a specified branch
//
//    * ListBranches, which lists all branches for a specified repository
//
//    * UpdateDefaultBranch, which changes the default branch for a repository
//
// Information about committed code in a repository, by calling the following:
//
//    * GetBlob, which returns the base-64 encoded content of an individual
//    Git blob object within a repository
//
//    * GetCommit, which returns information about a commit, including commit
//    messages and author and committer information
//
//    * GetDifferences, which returns information about the differences in a
//    valid commit specifier (such as a branch, tag, HEAD, commit ID or other
//    fully qualified reference)
//
// Triggers, by calling the following:
//
//    * GetRepositoryTriggers, which returns information about triggers configured
//    for a repository
//
//    * PutRepositoryTriggers, which replaces all triggers for a repository
//    and can be used to create or delete triggers
//
//    * TestRepositoryTriggers, which tests the functionality of a repository
//    trigger by sending data to the trigger target
//
// For information about how to use AWS CodeCommit, see the AWS CodeCommit User
// Guide (http://docs.aws.amazon.com/codecommit/latest/userguide/welcome.html).
//
// See https://docs.aws.amazon.com/goto/WebAPI/codecommit-2015-04-13 for more information on this service.
//
// See codecommit package documentation for more information.
// https://docs.aws.amazon.com/sdk-for-go/api/service/codecommit/
//
// Using the Client
//
// To use the client for AWS CodeCommit you will first need
// to create a new instance of it.
//
// When creating a client for an AWS service you'll first need to have a Session
// already created. The Session provides configuration that can be shared
// between multiple service clients. Additional configuration can be applied to
// the Session and service's client when they are constructed. The aws package's
// Config type contains several fields such as Region for the AWS Region the
// client should make API requests too. The optional Config value can be provided
// as the variadic argument for Sessions and client creation.
//
// Once the service's client is created you can use it to make API requests the
// AWS service. These clients are safe to use concurrently.
//
//   // Create a session to share configuration, and load external configuration.
//   sess := session.Must(session.NewSession())
//
//   // Create the service's client with the session.
//   svc := codecommit.New(sess)
//
// See the SDK's documentation for more information on how to use service clients.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws package's Config type for more information on configuration options.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the AWS CodeCommit client CodeCommit for more
// information on creating the service's client.
// https://docs.aws.amazon.com/sdk-for-go/api/service/codecommit/#New
//
// Once the client is created you can make an API request to the service.
// Each API method takes a input parameter, and returns the service response
// and an error.
//
// The API method will document which error codes the service can be returned
// by the operation if the service models the API operation's errors. These
// errors will also be available as const strings prefixed with "ErrCode".
//
//   result, err := svc.BatchGetRepositories(params)
//   if err != nil {
//       // Cast err to awserr.Error to handle specific error codes.
//       aerr, ok := err.(awserr.Error)
//       if ok && aerr.Code() == <error code to check for> {
//           // Specific error code handling
//       }
//       return err
//   }
//
//   fmt.Println("BatchGetRepositories result:")
//   fmt.Println(result)
//
// Using the Client with Context
//
// The service's client also provides methods to make API requests with a Context
// value. This allows you to control the timeout, and cancellation of pending
// requests. These methods also take request Option as variadic parameter to apply
// additional configuration to the API request.
//
//   ctx := context.Background()
//
//   result, err := svc.BatchGetRepositoriesWithContext(ctx, params)
//
// See the request package documentation for more information on using Context pattern
// with the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/request/
package codecommit
