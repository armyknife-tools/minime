// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package kms

const (

	// ErrCodeAlreadyExistsException for service response error code
	// "AlreadyExistsException".
	//
	// The request was rejected because it attempted to create a resource that already
	// exists.
	ErrCodeAlreadyExistsException = "AlreadyExistsException"

	// ErrCodeDependencyTimeoutException for service response error code
	// "DependencyTimeoutException".
	//
	// The system timed out while trying to fulfill the request. The request can
	// be retried.
	ErrCodeDependencyTimeoutException = "DependencyTimeoutException"

	// ErrCodeDisabledException for service response error code
	// "DisabledException".
	//
	// The request was rejected because the specified CMK is not enabled.
	ErrCodeDisabledException = "DisabledException"

	// ErrCodeExpiredImportTokenException for service response error code
	// "ExpiredImportTokenException".
	//
	// The request was rejected because the provided import token is expired. Use
	// GetParametersForImport to get a new import token and public key, use the
	// new public key to encrypt the key material, and then try the request again.
	ErrCodeExpiredImportTokenException = "ExpiredImportTokenException"

	// ErrCodeIncorrectKeyMaterialException for service response error code
	// "IncorrectKeyMaterialException".
	//
	// The request was rejected because the provided key material is invalid or
	// is not the same key material that was previously imported into this customer
	// master key (CMK).
	ErrCodeIncorrectKeyMaterialException = "IncorrectKeyMaterialException"

	// ErrCodeInternalException for service response error code
	// "InternalException".
	//
	// The request was rejected because an internal exception occurred. The request
	// can be retried.
	ErrCodeInternalException = "InternalException"

	// ErrCodeInvalidAliasNameException for service response error code
	// "InvalidAliasNameException".
	//
	// The request was rejected because the specified alias name is not valid.
	ErrCodeInvalidAliasNameException = "InvalidAliasNameException"

	// ErrCodeInvalidArnException for service response error code
	// "InvalidArnException".
	//
	// The request was rejected because a specified ARN was not valid.
	ErrCodeInvalidArnException = "InvalidArnException"

	// ErrCodeInvalidCiphertextException for service response error code
	// "InvalidCiphertextException".
	//
	// The request was rejected because the specified ciphertext, or additional
	// authenticated data incorporated into the ciphertext, such as the encryption
	// context, is corrupted, missing, or otherwise invalid.
	ErrCodeInvalidCiphertextException = "InvalidCiphertextException"

	// ErrCodeInvalidGrantIdException for service response error code
	// "InvalidGrantIdException".
	//
	// The request was rejected because the specified GrantId is not valid.
	ErrCodeInvalidGrantIdException = "InvalidGrantIdException"

	// ErrCodeInvalidGrantTokenException for service response error code
	// "InvalidGrantTokenException".
	//
	// The request was rejected because the specified grant token is not valid.
	ErrCodeInvalidGrantTokenException = "InvalidGrantTokenException"

	// ErrCodeInvalidImportTokenException for service response error code
	// "InvalidImportTokenException".
	//
	// The request was rejected because the provided import token is invalid or
	// is associated with a different customer master key (CMK).
	ErrCodeInvalidImportTokenException = "InvalidImportTokenException"

	// ErrCodeInvalidKeyUsageException for service response error code
	// "InvalidKeyUsageException".
	//
	// The request was rejected because the specified KeySpec value is not valid.
	ErrCodeInvalidKeyUsageException = "InvalidKeyUsageException"

	// ErrCodeInvalidMarkerException for service response error code
	// "InvalidMarkerException".
	//
	// The request was rejected because the marker that specifies where pagination
	// should next begin is not valid.
	ErrCodeInvalidMarkerException = "InvalidMarkerException"

	// ErrCodeInvalidStateException for service response error code
	// "InvalidStateException".
	//
	// The request was rejected because the state of the specified resource is not
	// valid for this request.
	//
	// For more information about how key state affects the use of a CMK, see How
	// Key State Affects Use of a Customer Master Key (http://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)
	// in the AWS Key Management Service Developer Guide.
	ErrCodeInvalidStateException = "InvalidStateException"

	// ErrCodeKeyUnavailableException for service response error code
	// "KeyUnavailableException".
	//
	// The request was rejected because the specified CMK was not available. The
	// request can be retried.
	ErrCodeKeyUnavailableException = "KeyUnavailableException"

	// ErrCodeLimitExceededException for service response error code
	// "LimitExceededException".
	//
	// The request was rejected because a limit was exceeded. For more information,
	// see Limits (http://docs.aws.amazon.com/kms/latest/developerguide/limits.html)
	// in the AWS Key Management Service Developer Guide.
	ErrCodeLimitExceededException = "LimitExceededException"

	// ErrCodeMalformedPolicyDocumentException for service response error code
	// "MalformedPolicyDocumentException".
	//
	// The request was rejected because the specified policy is not syntactically
	// or semantically correct.
	ErrCodeMalformedPolicyDocumentException = "MalformedPolicyDocumentException"

	// ErrCodeNotFoundException for service response error code
	// "NotFoundException".
	//
	// The request was rejected because the specified entity or resource could not
	// be found.
	ErrCodeNotFoundException = "NotFoundException"

	// ErrCodeTagException for service response error code
	// "TagException".
	//
	// The request was rejected because one or more tags are not valid.
	ErrCodeTagException = "TagException"

	// ErrCodeUnsupportedOperationException for service response error code
	// "UnsupportedOperationException".
	//
	// The request was rejected because a specified parameter is not supported or
	// a specified resource is not valid for this operation.
	ErrCodeUnsupportedOperationException = "UnsupportedOperationException"
)
