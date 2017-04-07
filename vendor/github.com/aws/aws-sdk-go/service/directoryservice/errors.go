// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package directoryservice

const (

	// ErrCodeAuthenticationFailedException for service response error code
	// "AuthenticationFailedException".
	//
	// An authentication error occurred.
	ErrCodeAuthenticationFailedException = "AuthenticationFailedException"

	// ErrCodeClientException for service response error code
	// "ClientException".
	//
	// A client exception has occurred.
	ErrCodeClientException = "ClientException"

	// ErrCodeDirectoryLimitExceededException for service response error code
	// "DirectoryLimitExceededException".
	//
	// The maximum number of directories in the region has been reached. You can
	// use the GetDirectoryLimits operation to determine your directory limits in
	// the region.
	ErrCodeDirectoryLimitExceededException = "DirectoryLimitExceededException"

	// ErrCodeDirectoryUnavailableException for service response error code
	// "DirectoryUnavailableException".
	//
	// The specified directory is unavailable or could not be found.
	ErrCodeDirectoryUnavailableException = "DirectoryUnavailableException"

	// ErrCodeEntityAlreadyExistsException for service response error code
	// "EntityAlreadyExistsException".
	//
	// The specified entity already exists.
	ErrCodeEntityAlreadyExistsException = "EntityAlreadyExistsException"

	// ErrCodeEntityDoesNotExistException for service response error code
	// "EntityDoesNotExistException".
	//
	// The specified entity could not be found.
	ErrCodeEntityDoesNotExistException = "EntityDoesNotExistException"

	// ErrCodeInsufficientPermissionsException for service response error code
	// "InsufficientPermissionsException".
	//
	// The account does not have sufficient permission to perform the operation.
	ErrCodeInsufficientPermissionsException = "InsufficientPermissionsException"

	// ErrCodeInvalidNextTokenException for service response error code
	// "InvalidNextTokenException".
	//
	// The NextToken value is not valid.
	ErrCodeInvalidNextTokenException = "InvalidNextTokenException"

	// ErrCodeInvalidParameterException for service response error code
	// "InvalidParameterException".
	//
	// One or more parameters are not valid.
	ErrCodeInvalidParameterException = "InvalidParameterException"

	// ErrCodeIpRouteLimitExceededException for service response error code
	// "IpRouteLimitExceededException".
	//
	// The maximum allowed number of IP addresses was exceeded. The default limit
	// is 100 IP address blocks.
	ErrCodeIpRouteLimitExceededException = "IpRouteLimitExceededException"

	// ErrCodeServiceException for service response error code
	// "ServiceException".
	//
	// An exception has occurred in AWS Directory Service.
	ErrCodeServiceException = "ServiceException"

	// ErrCodeSnapshotLimitExceededException for service response error code
	// "SnapshotLimitExceededException".
	//
	// The maximum number of manual snapshots for the directory has been reached.
	// You can use the GetSnapshotLimits operation to determine the snapshot limits
	// for a directory.
	ErrCodeSnapshotLimitExceededException = "SnapshotLimitExceededException"

	// ErrCodeTagLimitExceededException for service response error code
	// "TagLimitExceededException".
	//
	// The maximum allowed number of tags was exceeded.
	ErrCodeTagLimitExceededException = "TagLimitExceededException"

	// ErrCodeUnsupportedOperationException for service response error code
	// "UnsupportedOperationException".
	//
	// The operation is not supported.
	ErrCodeUnsupportedOperationException = "UnsupportedOperationException"
)
