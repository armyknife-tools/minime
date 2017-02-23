// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package codecommit

const (

	// ErrCodeBlobIdDoesNotExistException for service response error code
	// "BlobIdDoesNotExistException".
	//
	// The specified blob does not exist.
	ErrCodeBlobIdDoesNotExistException = "BlobIdDoesNotExistException"

	// ErrCodeBlobIdRequiredException for service response error code
	// "BlobIdRequiredException".
	//
	// A blob ID is required but was not specified.
	ErrCodeBlobIdRequiredException = "BlobIdRequiredException"

	// ErrCodeBranchDoesNotExistException for service response error code
	// "BranchDoesNotExistException".
	//
	// The specified branch does not exist.
	ErrCodeBranchDoesNotExistException = "BranchDoesNotExistException"

	// ErrCodeBranchNameExistsException for service response error code
	// "BranchNameExistsException".
	//
	// The specified branch name already exists.
	ErrCodeBranchNameExistsException = "BranchNameExistsException"

	// ErrCodeBranchNameRequiredException for service response error code
	// "BranchNameRequiredException".
	//
	// A branch name is required but was not specified.
	ErrCodeBranchNameRequiredException = "BranchNameRequiredException"

	// ErrCodeCommitDoesNotExistException for service response error code
	// "CommitDoesNotExistException".
	//
	// The specified commit does not exist or no commit was specified, and the specified
	// repository has no default branch.
	ErrCodeCommitDoesNotExistException = "CommitDoesNotExistException"

	// ErrCodeCommitIdDoesNotExistException for service response error code
	// "CommitIdDoesNotExistException".
	//
	// The specified commit ID does not exist.
	ErrCodeCommitIdDoesNotExistException = "CommitIdDoesNotExistException"

	// ErrCodeCommitIdRequiredException for service response error code
	// "CommitIdRequiredException".
	//
	// A commit ID was not specified.
	ErrCodeCommitIdRequiredException = "CommitIdRequiredException"

	// ErrCodeCommitRequiredException for service response error code
	// "CommitRequiredException".
	//
	// A commit was not specified.
	ErrCodeCommitRequiredException = "CommitRequiredException"

	// ErrCodeEncryptionIntegrityChecksFailedException for service response error code
	// "EncryptionIntegrityChecksFailedException".
	//
	// An encryption integrity check failed.
	ErrCodeEncryptionIntegrityChecksFailedException = "EncryptionIntegrityChecksFailedException"

	// ErrCodeEncryptionKeyAccessDeniedException for service response error code
	// "EncryptionKeyAccessDeniedException".
	//
	// An encryption key could not be accessed.
	ErrCodeEncryptionKeyAccessDeniedException = "EncryptionKeyAccessDeniedException"

	// ErrCodeEncryptionKeyDisabledException for service response error code
	// "EncryptionKeyDisabledException".
	//
	// The encryption key is disabled.
	ErrCodeEncryptionKeyDisabledException = "EncryptionKeyDisabledException"

	// ErrCodeEncryptionKeyNotFoundException for service response error code
	// "EncryptionKeyNotFoundException".
	//
	// No encryption key was found.
	ErrCodeEncryptionKeyNotFoundException = "EncryptionKeyNotFoundException"

	// ErrCodeEncryptionKeyUnavailableException for service response error code
	// "EncryptionKeyUnavailableException".
	//
	// The encryption key is not available.
	ErrCodeEncryptionKeyUnavailableException = "EncryptionKeyUnavailableException"

	// ErrCodeFileTooLargeException for service response error code
	// "FileTooLargeException".
	//
	// The specified file exceeds the file size limit for AWS CodeCommit. For more
	// information about limits in AWS CodeCommit, see AWS CodeCommit User Guide
	// (http://docs.aws.amazon.com/codecommit/latest/userguide/limits.html).
	ErrCodeFileTooLargeException = "FileTooLargeException"

	// ErrCodeInvalidBlobIdException for service response error code
	// "InvalidBlobIdException".
	//
	// The specified blob is not valid.
	ErrCodeInvalidBlobIdException = "InvalidBlobIdException"

	// ErrCodeInvalidBranchNameException for service response error code
	// "InvalidBranchNameException".
	//
	// The specified branch name is not valid.
	ErrCodeInvalidBranchNameException = "InvalidBranchNameException"

	// ErrCodeInvalidCommitException for service response error code
	// "InvalidCommitException".
	//
	// The specified commit is not valid.
	ErrCodeInvalidCommitException = "InvalidCommitException"

	// ErrCodeInvalidCommitIdException for service response error code
	// "InvalidCommitIdException".
	//
	// The specified commit ID is not valid.
	ErrCodeInvalidCommitIdException = "InvalidCommitIdException"

	// ErrCodeInvalidContinuationTokenException for service response error code
	// "InvalidContinuationTokenException".
	//
	// The specified continuation token is not valid.
	ErrCodeInvalidContinuationTokenException = "InvalidContinuationTokenException"

	// ErrCodeInvalidMaxResultsException for service response error code
	// "InvalidMaxResultsException".
	//
	// The specified number of maximum results is not valid.
	ErrCodeInvalidMaxResultsException = "InvalidMaxResultsException"

	// ErrCodeInvalidOrderException for service response error code
	// "InvalidOrderException".
	//
	// The specified sort order is not valid.
	ErrCodeInvalidOrderException = "InvalidOrderException"

	// ErrCodeInvalidPathException for service response error code
	// "InvalidPathException".
	//
	// The specified path is not valid.
	ErrCodeInvalidPathException = "InvalidPathException"

	// ErrCodeInvalidRepositoryDescriptionException for service response error code
	// "InvalidRepositoryDescriptionException".
	//
	// The specified repository description is not valid.
	ErrCodeInvalidRepositoryDescriptionException = "InvalidRepositoryDescriptionException"

	// ErrCodeInvalidRepositoryNameException for service response error code
	// "InvalidRepositoryNameException".
	//
	// At least one specified repository name is not valid.
	//
	// This exception only occurs when a specified repository name is not valid.
	// Other exceptions occur when a required repository parameter is missing, or
	// when a specified repository does not exist.
	ErrCodeInvalidRepositoryNameException = "InvalidRepositoryNameException"

	// ErrCodeInvalidRepositoryTriggerBranchNameException for service response error code
	// "InvalidRepositoryTriggerBranchNameException".
	//
	// One or more branch names specified for the trigger is not valid.
	ErrCodeInvalidRepositoryTriggerBranchNameException = "InvalidRepositoryTriggerBranchNameException"

	// ErrCodeInvalidRepositoryTriggerCustomDataException for service response error code
	// "InvalidRepositoryTriggerCustomDataException".
	//
	// The custom data provided for the trigger is not valid.
	ErrCodeInvalidRepositoryTriggerCustomDataException = "InvalidRepositoryTriggerCustomDataException"

	// ErrCodeInvalidRepositoryTriggerDestinationArnException for service response error code
	// "InvalidRepositoryTriggerDestinationArnException".
	//
	// The Amazon Resource Name (ARN) for the trigger is not valid for the specified
	// destination. The most common reason for this error is that the ARN does not
	// meet the requirements for the service type.
	ErrCodeInvalidRepositoryTriggerDestinationArnException = "InvalidRepositoryTriggerDestinationArnException"

	// ErrCodeInvalidRepositoryTriggerEventsException for service response error code
	// "InvalidRepositoryTriggerEventsException".
	//
	// One or more events specified for the trigger is not valid. Check to make
	// sure that all events specified match the requirements for allowed events.
	ErrCodeInvalidRepositoryTriggerEventsException = "InvalidRepositoryTriggerEventsException"

	// ErrCodeInvalidRepositoryTriggerNameException for service response error code
	// "InvalidRepositoryTriggerNameException".
	//
	// The name of the trigger is not valid.
	ErrCodeInvalidRepositoryTriggerNameException = "InvalidRepositoryTriggerNameException"

	// ErrCodeInvalidRepositoryTriggerRegionException for service response error code
	// "InvalidRepositoryTriggerRegionException".
	//
	// The region for the trigger target does not match the region for the repository.
	// Triggers must be created in the same region as the target for the trigger.
	ErrCodeInvalidRepositoryTriggerRegionException = "InvalidRepositoryTriggerRegionException"

	// ErrCodeInvalidSortByException for service response error code
	// "InvalidSortByException".
	//
	// The specified sort by value is not valid.
	ErrCodeInvalidSortByException = "InvalidSortByException"

	// ErrCodeMaximumBranchesExceededException for service response error code
	// "MaximumBranchesExceededException".
	//
	// The number of branches for the trigger was exceeded.
	ErrCodeMaximumBranchesExceededException = "MaximumBranchesExceededException"

	// ErrCodeMaximumRepositoryNamesExceededException for service response error code
	// "MaximumRepositoryNamesExceededException".
	//
	// The maximum number of allowed repository names was exceeded. Currently, this
	// number is 25.
	ErrCodeMaximumRepositoryNamesExceededException = "MaximumRepositoryNamesExceededException"

	// ErrCodeMaximumRepositoryTriggersExceededException for service response error code
	// "MaximumRepositoryTriggersExceededException".
	//
	// The number of triggers allowed for the repository was exceeded.
	ErrCodeMaximumRepositoryTriggersExceededException = "MaximumRepositoryTriggersExceededException"

	// ErrCodePathDoesNotExistException for service response error code
	// "PathDoesNotExistException".
	//
	// The specified path does not exist.
	ErrCodePathDoesNotExistException = "PathDoesNotExistException"

	// ErrCodeRepositoryDoesNotExistException for service response error code
	// "RepositoryDoesNotExistException".
	//
	// The specified repository does not exist.
	ErrCodeRepositoryDoesNotExistException = "RepositoryDoesNotExistException"

	// ErrCodeRepositoryLimitExceededException for service response error code
	// "RepositoryLimitExceededException".
	//
	// A repository resource limit was exceeded.
	ErrCodeRepositoryLimitExceededException = "RepositoryLimitExceededException"

	// ErrCodeRepositoryNameExistsException for service response error code
	// "RepositoryNameExistsException".
	//
	// The specified repository name already exists.
	ErrCodeRepositoryNameExistsException = "RepositoryNameExistsException"

	// ErrCodeRepositoryNameRequiredException for service response error code
	// "RepositoryNameRequiredException".
	//
	// A repository name is required but was not specified.
	ErrCodeRepositoryNameRequiredException = "RepositoryNameRequiredException"

	// ErrCodeRepositoryNamesRequiredException for service response error code
	// "RepositoryNamesRequiredException".
	//
	// A repository names object is required but was not specified.
	ErrCodeRepositoryNamesRequiredException = "RepositoryNamesRequiredException"

	// ErrCodeRepositoryTriggerBranchNameListRequiredException for service response error code
	// "RepositoryTriggerBranchNameListRequiredException".
	//
	// At least one branch name is required but was not specified in the trigger
	// configuration.
	ErrCodeRepositoryTriggerBranchNameListRequiredException = "RepositoryTriggerBranchNameListRequiredException"

	// ErrCodeRepositoryTriggerDestinationArnRequiredException for service response error code
	// "RepositoryTriggerDestinationArnRequiredException".
	//
	// A destination ARN for the target service for the trigger is required but
	// was not specified.
	ErrCodeRepositoryTriggerDestinationArnRequiredException = "RepositoryTriggerDestinationArnRequiredException"

	// ErrCodeRepositoryTriggerEventsListRequiredException for service response error code
	// "RepositoryTriggerEventsListRequiredException".
	//
	// At least one event for the trigger is required but was not specified.
	ErrCodeRepositoryTriggerEventsListRequiredException = "RepositoryTriggerEventsListRequiredException"

	// ErrCodeRepositoryTriggerNameRequiredException for service response error code
	// "RepositoryTriggerNameRequiredException".
	//
	// A name for the trigger is required but was not specified.
	ErrCodeRepositoryTriggerNameRequiredException = "RepositoryTriggerNameRequiredException"

	// ErrCodeRepositoryTriggersListRequiredException for service response error code
	// "RepositoryTriggersListRequiredException".
	//
	// The list of triggers for the repository is required but was not specified.
	ErrCodeRepositoryTriggersListRequiredException = "RepositoryTriggersListRequiredException"
)
