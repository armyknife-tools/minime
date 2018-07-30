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
//    associated with your AWS account.
//
//    * CreateRepository, which creates an AWS CodeCommit repository.
//
//    * DeleteRepository, which deletes an AWS CodeCommit repository.
//
//    * GetRepository, which returns information about a specified repository.
//
//    * ListRepositories, which lists all AWS CodeCommit repositories associated
//    with your AWS account.
//
//    * UpdateRepositoryDescription, which sets or updates the description of
//    the repository.
//
//    * UpdateRepositoryName, which changes the name of the repository. If you
//    change the name of a repository, no other users of that repository will
//    be able to access it until you send them the new HTTPS or SSH URL to use.
//
// Branches, by calling the following:
//
//    * CreateBranch, which creates a new branch in a specified repository.
//
//    * DeleteBranch, which deletes the specified branch in a repository unless
//    it is the default branch.
//
//    * GetBranch, which returns information about a specified branch.
//
//    * ListBranches, which lists all branches for a specified repository.
//
//    * UpdateDefaultBranch, which changes the default branch for a repository.
//
// Files, by calling the following:
//
//    * PutFile, which adds or modifies a file in a specified repository and
//    branch.
//
// Information about committed code in a repository, by calling the following:
//
//    * GetBlob, which returns the base-64 encoded content of an individual
//    Git blob object within a repository.
//
//    * GetCommit, which returns information about a commit, including commit
//    messages and author and committer information.
//
//    * GetDifferences, which returns information about the differences in a
//    valid commit specifier (such as a branch, tag, HEAD, commit ID or other
//    fully qualified reference).
//
// Pull requests, by calling the following:
//
//    * CreatePullRequest, which creates a pull request in a specified repository.
//
//    * DescribePullRequestEvents, which returns information about one or more
//    pull request events.
//
//    * GetCommentsForPullRequest, which returns information about comments
//    on a specified pull request.
//
//    * GetMergeConflicts, which returns information about merge conflicts between
//    the source and destination branch in a pull request.
//
//    * GetPullRequest, which returns information about a specified pull request.
//
//    * ListPullRequests, which lists all pull requests for a repository.
//
//    * MergePullRequestByFastForward, which merges the source destination branch
//    of a pull request into the specified destination branch for that pull
//    request using the fast-forward merge option.
//
//    * PostCommentForPullRequest, which posts a comment to a pull request at
//    the specified line, file, or request.
//
//    * UpdatePullRequestDescription, which updates the description of a pull
//    request.
//
//    * UpdatePullRequestStatus, which updates the status of a pull request.
//
//    * UpdatePullRequestTitle, which updates the title of a pull request.
//
// Information about comments in a repository, by calling the following:
//
//    * DeleteCommentContent, which deletes the content of a comment on a commit
//    in a repository.
//
//    * GetComment, which returns information about a comment on a commit.
//
//    * GetCommentsForComparedCommit, which returns information about comments
//    on the comparison between two commit specifiers in a repository.
//
//    * PostCommentForComparedCommit, which creates a comment on the comparison
//    between two commit specifiers in a repository.
//
//    * PostCommentReply, which creates a reply to a comment.
//
//    * UpdateComment, which updates the content of a comment on a commit in
//    a repository.
//
// Triggers, by calling the following:
//
//    * GetRepositoryTriggers, which returns information about triggers configured
//    for a repository.
//
//    * PutRepositoryTriggers, which replaces all triggers for a repository
//    and can be used to create or delete triggers.
//
//    * TestRepositoryTriggers, which tests the functionality of a repository
//    trigger by sending data to the trigger target.
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
// To contact AWS CodeCommit with the SDK use the New function to create
// a new service client. With that client you can make API requests to the service.
// These clients are safe to use concurrently.
//
// See the SDK's documentation for more information on how to use the SDK.
// https://docs.aws.amazon.com/sdk-for-go/api/
//
// See aws.Config documentation for more information on configuring SDK clients.
// https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
//
// See the AWS CodeCommit client CodeCommit for more
// information on creating client for this service.
// https://docs.aws.amazon.com/sdk-for-go/api/service/codecommit/#New
package codecommit
