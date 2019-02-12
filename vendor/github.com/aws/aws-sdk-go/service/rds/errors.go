// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package rds

const (

	// ErrCodeAuthorizationAlreadyExistsFault for service response error code
	// "AuthorizationAlreadyExists".
	//
	// The specified CIDRIP or Amazon EC2 security group is already authorized for
	// the specified DB security group.
	ErrCodeAuthorizationAlreadyExistsFault = "AuthorizationAlreadyExists"

	// ErrCodeAuthorizationNotFoundFault for service response error code
	// "AuthorizationNotFound".
	//
	// The specified CIDRIP or Amazon EC2 security group isn't authorized for the
	// specified DB security group.
	//
	// RDS also may not be authorized by using IAM to perform necessary actions
	// on your behalf.
	ErrCodeAuthorizationNotFoundFault = "AuthorizationNotFound"

	// ErrCodeAuthorizationQuotaExceededFault for service response error code
	// "AuthorizationQuotaExceeded".
	//
	// The DB security group authorization quota has been reached.
	ErrCodeAuthorizationQuotaExceededFault = "AuthorizationQuotaExceeded"

	// ErrCodeBackupPolicyNotFoundFault for service response error code
	// "BackupPolicyNotFoundFault".
	//
	// The backup policy was not found.
	ErrCodeBackupPolicyNotFoundFault = "BackupPolicyNotFoundFault"

	// ErrCodeCertificateNotFoundFault for service response error code
	// "CertificateNotFound".
	//
	// CertificateIdentifier doesn't refer to an existing certificate.
	ErrCodeCertificateNotFoundFault = "CertificateNotFound"

	// ErrCodeDBClusterAlreadyExistsFault for service response error code
	// "DBClusterAlreadyExistsFault".
	//
	// The user already has a DB cluster with the given identifier.
	ErrCodeDBClusterAlreadyExistsFault = "DBClusterAlreadyExistsFault"

	// ErrCodeDBClusterBacktrackNotFoundFault for service response error code
	// "DBClusterBacktrackNotFoundFault".
	//
	// BacktrackIdentifier doesn't refer to an existing backtrack.
	ErrCodeDBClusterBacktrackNotFoundFault = "DBClusterBacktrackNotFoundFault"

	// ErrCodeDBClusterEndpointAlreadyExistsFault for service response error code
	// "DBClusterEndpointAlreadyExistsFault".
	//
	// The specified custom endpoint can't be created because it already exists.
	ErrCodeDBClusterEndpointAlreadyExistsFault = "DBClusterEndpointAlreadyExistsFault"

	// ErrCodeDBClusterEndpointNotFoundFault for service response error code
	// "DBClusterEndpointNotFoundFault".
	//
	// The specified custom endpoint doesn't exist.
	ErrCodeDBClusterEndpointNotFoundFault = "DBClusterEndpointNotFoundFault"

	// ErrCodeDBClusterEndpointQuotaExceededFault for service response error code
	// "DBClusterEndpointQuotaExceededFault".
	//
	// The cluster already has the maximum number of custom endpoints.
	ErrCodeDBClusterEndpointQuotaExceededFault = "DBClusterEndpointQuotaExceededFault"

	// ErrCodeDBClusterNotFoundFault for service response error code
	// "DBClusterNotFoundFault".
	//
	// DBClusterIdentifier doesn't refer to an existing DB cluster.
	ErrCodeDBClusterNotFoundFault = "DBClusterNotFoundFault"

	// ErrCodeDBClusterParameterGroupNotFoundFault for service response error code
	// "DBClusterParameterGroupNotFound".
	//
	// DBClusterParameterGroupName doesn't refer to an existing DB cluster parameter
	// group.
	ErrCodeDBClusterParameterGroupNotFoundFault = "DBClusterParameterGroupNotFound"

	// ErrCodeDBClusterQuotaExceededFault for service response error code
	// "DBClusterQuotaExceededFault".
	//
	// The user attempted to create a new DB cluster and the user has already reached
	// the maximum allowed DB cluster quota.
	ErrCodeDBClusterQuotaExceededFault = "DBClusterQuotaExceededFault"

	// ErrCodeDBClusterRoleAlreadyExistsFault for service response error code
	// "DBClusterRoleAlreadyExists".
	//
	// The specified IAM role Amazon Resource Name (ARN) is already associated with
	// the specified DB cluster.
	ErrCodeDBClusterRoleAlreadyExistsFault = "DBClusterRoleAlreadyExists"

	// ErrCodeDBClusterRoleNotFoundFault for service response error code
	// "DBClusterRoleNotFound".
	//
	// The specified IAM role Amazon Resource Name (ARN) isn't associated with the
	// specified DB cluster.
	ErrCodeDBClusterRoleNotFoundFault = "DBClusterRoleNotFound"

	// ErrCodeDBClusterRoleQuotaExceededFault for service response error code
	// "DBClusterRoleQuotaExceeded".
	//
	// You have exceeded the maximum number of IAM roles that can be associated
	// with the specified DB cluster.
	ErrCodeDBClusterRoleQuotaExceededFault = "DBClusterRoleQuotaExceeded"

	// ErrCodeDBClusterSnapshotAlreadyExistsFault for service response error code
	// "DBClusterSnapshotAlreadyExistsFault".
	//
	// The user already has a DB cluster snapshot with the given identifier.
	ErrCodeDBClusterSnapshotAlreadyExistsFault = "DBClusterSnapshotAlreadyExistsFault"

	// ErrCodeDBClusterSnapshotNotFoundFault for service response error code
	// "DBClusterSnapshotNotFoundFault".
	//
	// DBClusterSnapshotIdentifier doesn't refer to an existing DB cluster snapshot.
	ErrCodeDBClusterSnapshotNotFoundFault = "DBClusterSnapshotNotFoundFault"

	// ErrCodeDBInstanceAlreadyExistsFault for service response error code
	// "DBInstanceAlreadyExists".
	//
	// The user already has a DB instance with the given identifier.
	ErrCodeDBInstanceAlreadyExistsFault = "DBInstanceAlreadyExists"

	// ErrCodeDBInstanceAutomatedBackupNotFoundFault for service response error code
	// "DBInstanceAutomatedBackupNotFound".
	//
	// No automated backup for this DB instance was found.
	ErrCodeDBInstanceAutomatedBackupNotFoundFault = "DBInstanceAutomatedBackupNotFound"

	// ErrCodeDBInstanceAutomatedBackupQuotaExceededFault for service response error code
	// "DBInstanceAutomatedBackupQuotaExceeded".
	//
	// The quota for retained automated backups was exceeded. This prevents you
	// from retaining any additional automated backups. The retained automated backups
	// quota is the same as your DB Instance quota.
	ErrCodeDBInstanceAutomatedBackupQuotaExceededFault = "DBInstanceAutomatedBackupQuotaExceeded"

	// ErrCodeDBInstanceNotFoundFault for service response error code
	// "DBInstanceNotFound".
	//
	// DBInstanceIdentifier doesn't refer to an existing DB instance.
	ErrCodeDBInstanceNotFoundFault = "DBInstanceNotFound"

	// ErrCodeDBInstanceRoleAlreadyExistsFault for service response error code
	// "DBInstanceRoleAlreadyExists".
	//
	// The specified RoleArn or FeatureName value is already associated with the
	// DB instance.
	ErrCodeDBInstanceRoleAlreadyExistsFault = "DBInstanceRoleAlreadyExists"

	// ErrCodeDBInstanceRoleNotFoundFault for service response error code
	// "DBInstanceRoleNotFound".
	//
	// The specified RoleArn value doesn't match the specifed feature for the DB
	// instance.
	ErrCodeDBInstanceRoleNotFoundFault = "DBInstanceRoleNotFound"

	// ErrCodeDBInstanceRoleQuotaExceededFault for service response error code
	// "DBInstanceRoleQuotaExceeded".
	//
	// You can't associate any more AWS Identity and Access Management (IAM) roles
	// with the DB instance because the quota has been reached.
	ErrCodeDBInstanceRoleQuotaExceededFault = "DBInstanceRoleQuotaExceeded"

	// ErrCodeDBLogFileNotFoundFault for service response error code
	// "DBLogFileNotFoundFault".
	//
	// LogFileName doesn't refer to an existing DB log file.
	ErrCodeDBLogFileNotFoundFault = "DBLogFileNotFoundFault"

	// ErrCodeDBParameterGroupAlreadyExistsFault for service response error code
	// "DBParameterGroupAlreadyExists".
	//
	// A DB parameter group with the same name exists.
	ErrCodeDBParameterGroupAlreadyExistsFault = "DBParameterGroupAlreadyExists"

	// ErrCodeDBParameterGroupNotFoundFault for service response error code
	// "DBParameterGroupNotFound".
	//
	// DBParameterGroupName doesn't refer to an existing DB parameter group.
	ErrCodeDBParameterGroupNotFoundFault = "DBParameterGroupNotFound"

	// ErrCodeDBParameterGroupQuotaExceededFault for service response error code
	// "DBParameterGroupQuotaExceeded".
	//
	// The request would result in the user exceeding the allowed number of DB parameter
	// groups.
	ErrCodeDBParameterGroupQuotaExceededFault = "DBParameterGroupQuotaExceeded"

	// ErrCodeDBSecurityGroupAlreadyExistsFault for service response error code
	// "DBSecurityGroupAlreadyExists".
	//
	// A DB security group with the name specified in DBSecurityGroupName already
	// exists.
	ErrCodeDBSecurityGroupAlreadyExistsFault = "DBSecurityGroupAlreadyExists"

	// ErrCodeDBSecurityGroupNotFoundFault for service response error code
	// "DBSecurityGroupNotFound".
	//
	// DBSecurityGroupName doesn't refer to an existing DB security group.
	ErrCodeDBSecurityGroupNotFoundFault = "DBSecurityGroupNotFound"

	// ErrCodeDBSecurityGroupNotSupportedFault for service response error code
	// "DBSecurityGroupNotSupported".
	//
	// A DB security group isn't allowed for this action.
	ErrCodeDBSecurityGroupNotSupportedFault = "DBSecurityGroupNotSupported"

	// ErrCodeDBSecurityGroupQuotaExceededFault for service response error code
	// "QuotaExceeded.DBSecurityGroup".
	//
	// The request would result in the user exceeding the allowed number of DB security
	// groups.
	ErrCodeDBSecurityGroupQuotaExceededFault = "QuotaExceeded.DBSecurityGroup"

	// ErrCodeDBSnapshotAlreadyExistsFault for service response error code
	// "DBSnapshotAlreadyExists".
	//
	// DBSnapshotIdentifier is already used by an existing snapshot.
	ErrCodeDBSnapshotAlreadyExistsFault = "DBSnapshotAlreadyExists"

	// ErrCodeDBSnapshotNotFoundFault for service response error code
	// "DBSnapshotNotFound".
	//
	// DBSnapshotIdentifier doesn't refer to an existing DB snapshot.
	ErrCodeDBSnapshotNotFoundFault = "DBSnapshotNotFound"

	// ErrCodeDBSubnetGroupAlreadyExistsFault for service response error code
	// "DBSubnetGroupAlreadyExists".
	//
	// DBSubnetGroupName is already used by an existing DB subnet group.
	ErrCodeDBSubnetGroupAlreadyExistsFault = "DBSubnetGroupAlreadyExists"

	// ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs for service response error code
	// "DBSubnetGroupDoesNotCoverEnoughAZs".
	//
	// Subnets in the DB subnet group should cover at least two Availability Zones
	// unless there is only one Availability Zone.
	ErrCodeDBSubnetGroupDoesNotCoverEnoughAZs = "DBSubnetGroupDoesNotCoverEnoughAZs"

	// ErrCodeDBSubnetGroupNotAllowedFault for service response error code
	// "DBSubnetGroupNotAllowedFault".
	//
	// The DBSubnetGroup shouldn't be specified while creating read replicas that
	// lie in the same region as the source instance.
	ErrCodeDBSubnetGroupNotAllowedFault = "DBSubnetGroupNotAllowedFault"

	// ErrCodeDBSubnetGroupNotFoundFault for service response error code
	// "DBSubnetGroupNotFoundFault".
	//
	// DBSubnetGroupName doesn't refer to an existing DB subnet group.
	ErrCodeDBSubnetGroupNotFoundFault = "DBSubnetGroupNotFoundFault"

	// ErrCodeDBSubnetGroupQuotaExceededFault for service response error code
	// "DBSubnetGroupQuotaExceeded".
	//
	// The request would result in the user exceeding the allowed number of DB subnet
	// groups.
	ErrCodeDBSubnetGroupQuotaExceededFault = "DBSubnetGroupQuotaExceeded"

	// ErrCodeDBSubnetQuotaExceededFault for service response error code
	// "DBSubnetQuotaExceededFault".
	//
	// The request would result in the user exceeding the allowed number of subnets
	// in a DB subnet groups.
	ErrCodeDBSubnetQuotaExceededFault = "DBSubnetQuotaExceededFault"

	// ErrCodeDBUpgradeDependencyFailureFault for service response error code
	// "DBUpgradeDependencyFailure".
	//
	// The DB upgrade failed because a resource the DB depends on can't be modified.
	ErrCodeDBUpgradeDependencyFailureFault = "DBUpgradeDependencyFailure"

	// ErrCodeDomainNotFoundFault for service response error code
	// "DomainNotFoundFault".
	//
	// Domain doesn't refer to an existing Active Directory domain.
	ErrCodeDomainNotFoundFault = "DomainNotFoundFault"

	// ErrCodeEventSubscriptionQuotaExceededFault for service response error code
	// "EventSubscriptionQuotaExceeded".
	//
	// You have reached the maximum number of event subscriptions.
	ErrCodeEventSubscriptionQuotaExceededFault = "EventSubscriptionQuotaExceeded"

	// ErrCodeGlobalClusterAlreadyExistsFault for service response error code
	// "GlobalClusterAlreadyExistsFault".
	ErrCodeGlobalClusterAlreadyExistsFault = "GlobalClusterAlreadyExistsFault"

	// ErrCodeGlobalClusterNotFoundFault for service response error code
	// "GlobalClusterNotFoundFault".
	ErrCodeGlobalClusterNotFoundFault = "GlobalClusterNotFoundFault"

	// ErrCodeGlobalClusterQuotaExceededFault for service response error code
	// "GlobalClusterQuotaExceededFault".
	ErrCodeGlobalClusterQuotaExceededFault = "GlobalClusterQuotaExceededFault"

	// ErrCodeInstanceQuotaExceededFault for service response error code
	// "InstanceQuotaExceeded".
	//
	// The request would result in the user exceeding the allowed number of DB instances.
	ErrCodeInstanceQuotaExceededFault = "InstanceQuotaExceeded"

	// ErrCodeInsufficientDBClusterCapacityFault for service response error code
	// "InsufficientDBClusterCapacityFault".
	//
	// The DB cluster doesn't have enough capacity for the current operation.
	ErrCodeInsufficientDBClusterCapacityFault = "InsufficientDBClusterCapacityFault"

	// ErrCodeInsufficientDBInstanceCapacityFault for service response error code
	// "InsufficientDBInstanceCapacity".
	//
	// The specified DB instance class isn't available in the specified Availability
	// Zone.
	ErrCodeInsufficientDBInstanceCapacityFault = "InsufficientDBInstanceCapacity"

	// ErrCodeInsufficientStorageClusterCapacityFault for service response error code
	// "InsufficientStorageClusterCapacity".
	//
	// There is insufficient storage available for the current action. You might
	// be able to resolve this error by updating your subnet group to use different
	// Availability Zones that have more storage available.
	ErrCodeInsufficientStorageClusterCapacityFault = "InsufficientStorageClusterCapacity"

	// ErrCodeInvalidDBClusterCapacityFault for service response error code
	// "InvalidDBClusterCapacityFault".
	//
	// Capacity isn't a valid Aurora Serverless DB cluster capacity. Valid capacity
	// values are 2, 4, 8, 16, 32, 64, 128, and 256.
	ErrCodeInvalidDBClusterCapacityFault = "InvalidDBClusterCapacityFault"

	// ErrCodeInvalidDBClusterEndpointStateFault for service response error code
	// "InvalidDBClusterEndpointStateFault".
	//
	// The requested operation can't be performed on the endpoint while the endpoint
	// is in this state.
	ErrCodeInvalidDBClusterEndpointStateFault = "InvalidDBClusterEndpointStateFault"

	// ErrCodeInvalidDBClusterSnapshotStateFault for service response error code
	// "InvalidDBClusterSnapshotStateFault".
	//
	// The supplied value isn't a valid DB cluster snapshot state.
	ErrCodeInvalidDBClusterSnapshotStateFault = "InvalidDBClusterSnapshotStateFault"

	// ErrCodeInvalidDBClusterStateFault for service response error code
	// "InvalidDBClusterStateFault".
	//
	// The requested operation can't be performed while the cluster is in this state.
	ErrCodeInvalidDBClusterStateFault = "InvalidDBClusterStateFault"

	// ErrCodeInvalidDBInstanceAutomatedBackupStateFault for service response error code
	// "InvalidDBInstanceAutomatedBackupState".
	//
	// The automated backup is in an invalid state. For example, this automated
	// backup is associated with an active instance.
	ErrCodeInvalidDBInstanceAutomatedBackupStateFault = "InvalidDBInstanceAutomatedBackupState"

	// ErrCodeInvalidDBInstanceStateFault for service response error code
	// "InvalidDBInstanceState".
	//
	// The DB instance isn't in a valid state.
	ErrCodeInvalidDBInstanceStateFault = "InvalidDBInstanceState"

	// ErrCodeInvalidDBParameterGroupStateFault for service response error code
	// "InvalidDBParameterGroupState".
	//
	// The DB parameter group is in use or is in an invalid state. If you are attempting
	// to delete the parameter group, you can't delete it when the parameter group
	// is in this state.
	ErrCodeInvalidDBParameterGroupStateFault = "InvalidDBParameterGroupState"

	// ErrCodeInvalidDBSecurityGroupStateFault for service response error code
	// "InvalidDBSecurityGroupState".
	//
	// The state of the DB security group doesn't allow deletion.
	ErrCodeInvalidDBSecurityGroupStateFault = "InvalidDBSecurityGroupState"

	// ErrCodeInvalidDBSnapshotStateFault for service response error code
	// "InvalidDBSnapshotState".
	//
	// The state of the DB snapshot doesn't allow deletion.
	ErrCodeInvalidDBSnapshotStateFault = "InvalidDBSnapshotState"

	// ErrCodeInvalidDBSubnetGroupFault for service response error code
	// "InvalidDBSubnetGroupFault".
	//
	// The DBSubnetGroup doesn't belong to the same VPC as that of an existing cross-region
	// read replica of the same source instance.
	ErrCodeInvalidDBSubnetGroupFault = "InvalidDBSubnetGroupFault"

	// ErrCodeInvalidDBSubnetGroupStateFault for service response error code
	// "InvalidDBSubnetGroupStateFault".
	//
	// The DB subnet group cannot be deleted because it's in use.
	ErrCodeInvalidDBSubnetGroupStateFault = "InvalidDBSubnetGroupStateFault"

	// ErrCodeInvalidDBSubnetStateFault for service response error code
	// "InvalidDBSubnetStateFault".
	//
	// The DB subnet isn't in the available state.
	ErrCodeInvalidDBSubnetStateFault = "InvalidDBSubnetStateFault"

	// ErrCodeInvalidEventSubscriptionStateFault for service response error code
	// "InvalidEventSubscriptionState".
	//
	// This error can occur if someone else is modifying a subscription. You should
	// retry the action.
	ErrCodeInvalidEventSubscriptionStateFault = "InvalidEventSubscriptionState"

	// ErrCodeInvalidGlobalClusterStateFault for service response error code
	// "InvalidGlobalClusterStateFault".
	ErrCodeInvalidGlobalClusterStateFault = "InvalidGlobalClusterStateFault"

	// ErrCodeInvalidOptionGroupStateFault for service response error code
	// "InvalidOptionGroupStateFault".
	//
	// The option group isn't in the available state.
	ErrCodeInvalidOptionGroupStateFault = "InvalidOptionGroupStateFault"

	// ErrCodeInvalidRestoreFault for service response error code
	// "InvalidRestoreFault".
	//
	// Cannot restore from VPC backup to non-VPC DB instance.
	ErrCodeInvalidRestoreFault = "InvalidRestoreFault"

	// ErrCodeInvalidS3BucketFault for service response error code
	// "InvalidS3BucketFault".
	//
	// The specified Amazon S3 bucket name can't be found or Amazon RDS isn't authorized
	// to access the specified Amazon S3 bucket. Verify the SourceS3BucketName and
	// S3IngestionRoleArn values and try again.
	ErrCodeInvalidS3BucketFault = "InvalidS3BucketFault"

	// ErrCodeInvalidSubnet for service response error code
	// "InvalidSubnet".
	//
	// The requested subnet is invalid, or multiple subnets were requested that
	// are not all in a common VPC.
	ErrCodeInvalidSubnet = "InvalidSubnet"

	// ErrCodeInvalidVPCNetworkStateFault for service response error code
	// "InvalidVPCNetworkStateFault".
	//
	// The DB subnet group doesn't cover all Availability Zones after it's created
	// because of users' change.
	ErrCodeInvalidVPCNetworkStateFault = "InvalidVPCNetworkStateFault"

	// ErrCodeKMSKeyNotAccessibleFault for service response error code
	// "KMSKeyNotAccessibleFault".
	//
	// An error occurred accessing an AWS KMS key.
	ErrCodeKMSKeyNotAccessibleFault = "KMSKeyNotAccessibleFault"

	// ErrCodeOptionGroupAlreadyExistsFault for service response error code
	// "OptionGroupAlreadyExistsFault".
	//
	// The option group you are trying to create already exists.
	ErrCodeOptionGroupAlreadyExistsFault = "OptionGroupAlreadyExistsFault"

	// ErrCodeOptionGroupNotFoundFault for service response error code
	// "OptionGroupNotFoundFault".
	//
	// The specified option group could not be found.
	ErrCodeOptionGroupNotFoundFault = "OptionGroupNotFoundFault"

	// ErrCodeOptionGroupQuotaExceededFault for service response error code
	// "OptionGroupQuotaExceededFault".
	//
	// The quota of 20 option groups was exceeded for this AWS account.
	ErrCodeOptionGroupQuotaExceededFault = "OptionGroupQuotaExceededFault"

	// ErrCodePointInTimeRestoreNotEnabledFault for service response error code
	// "PointInTimeRestoreNotEnabled".
	//
	// SourceDBInstanceIdentifier refers to a DB instance with BackupRetentionPeriod
	// equal to 0.
	ErrCodePointInTimeRestoreNotEnabledFault = "PointInTimeRestoreNotEnabled"

	// ErrCodeProvisionedIopsNotAvailableInAZFault for service response error code
	// "ProvisionedIopsNotAvailableInAZFault".
	//
	// Provisioned IOPS not available in the specified Availability Zone.
	ErrCodeProvisionedIopsNotAvailableInAZFault = "ProvisionedIopsNotAvailableInAZFault"

	// ErrCodeReservedDBInstanceAlreadyExistsFault for service response error code
	// "ReservedDBInstanceAlreadyExists".
	//
	// User already has a reservation with the given identifier.
	ErrCodeReservedDBInstanceAlreadyExistsFault = "ReservedDBInstanceAlreadyExists"

	// ErrCodeReservedDBInstanceNotFoundFault for service response error code
	// "ReservedDBInstanceNotFound".
	//
	// The specified reserved DB Instance not found.
	ErrCodeReservedDBInstanceNotFoundFault = "ReservedDBInstanceNotFound"

	// ErrCodeReservedDBInstanceQuotaExceededFault for service response error code
	// "ReservedDBInstanceQuotaExceeded".
	//
	// Request would exceed the user's DB Instance quota.
	ErrCodeReservedDBInstanceQuotaExceededFault = "ReservedDBInstanceQuotaExceeded"

	// ErrCodeReservedDBInstancesOfferingNotFoundFault for service response error code
	// "ReservedDBInstancesOfferingNotFound".
	//
	// Specified offering does not exist.
	ErrCodeReservedDBInstancesOfferingNotFoundFault = "ReservedDBInstancesOfferingNotFound"

	// ErrCodeResourceNotFoundFault for service response error code
	// "ResourceNotFoundFault".
	//
	// The specified resource ID was not found.
	ErrCodeResourceNotFoundFault = "ResourceNotFoundFault"

	// ErrCodeSNSInvalidTopicFault for service response error code
	// "SNSInvalidTopic".
	//
	// SNS has responded that there is a problem with the SND topic specified.
	ErrCodeSNSInvalidTopicFault = "SNSInvalidTopic"

	// ErrCodeSNSNoAuthorizationFault for service response error code
	// "SNSNoAuthorization".
	//
	// You do not have permission to publish to the SNS topic ARN.
	ErrCodeSNSNoAuthorizationFault = "SNSNoAuthorization"

	// ErrCodeSNSTopicArnNotFoundFault for service response error code
	// "SNSTopicArnNotFound".
	//
	// The SNS topic ARN does not exist.
	ErrCodeSNSTopicArnNotFoundFault = "SNSTopicArnNotFound"

	// ErrCodeSharedSnapshotQuotaExceededFault for service response error code
	// "SharedSnapshotQuotaExceeded".
	//
	// You have exceeded the maximum number of accounts that you can share a manual
	// DB snapshot with.
	ErrCodeSharedSnapshotQuotaExceededFault = "SharedSnapshotQuotaExceeded"

	// ErrCodeSnapshotQuotaExceededFault for service response error code
	// "SnapshotQuotaExceeded".
	//
	// The request would result in the user exceeding the allowed number of DB snapshots.
	ErrCodeSnapshotQuotaExceededFault = "SnapshotQuotaExceeded"

	// ErrCodeSourceNotFoundFault for service response error code
	// "SourceNotFound".
	//
	// The requested source could not be found.
	ErrCodeSourceNotFoundFault = "SourceNotFound"

	// ErrCodeStorageQuotaExceededFault for service response error code
	// "StorageQuotaExceeded".
	//
	// The request would result in the user exceeding the allowed amount of storage
	// available across all DB instances.
	ErrCodeStorageQuotaExceededFault = "StorageQuotaExceeded"

	// ErrCodeStorageTypeNotSupportedFault for service response error code
	// "StorageTypeNotSupported".
	//
	// Storage of the StorageType specified can't be associated with the DB instance.
	ErrCodeStorageTypeNotSupportedFault = "StorageTypeNotSupported"

	// ErrCodeSubnetAlreadyInUse for service response error code
	// "SubnetAlreadyInUse".
	//
	// The DB subnet is already in use in the Availability Zone.
	ErrCodeSubnetAlreadyInUse = "SubnetAlreadyInUse"

	// ErrCodeSubscriptionAlreadyExistFault for service response error code
	// "SubscriptionAlreadyExist".
	//
	// The supplied subscription name already exists.
	ErrCodeSubscriptionAlreadyExistFault = "SubscriptionAlreadyExist"

	// ErrCodeSubscriptionCategoryNotFoundFault for service response error code
	// "SubscriptionCategoryNotFound".
	//
	// The supplied category does not exist.
	ErrCodeSubscriptionCategoryNotFoundFault = "SubscriptionCategoryNotFound"

	// ErrCodeSubscriptionNotFoundFault for service response error code
	// "SubscriptionNotFound".
	//
	// The subscription name does not exist.
	ErrCodeSubscriptionNotFoundFault = "SubscriptionNotFound"
)
