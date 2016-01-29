// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package redshift_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleRedshift_AuthorizeClusterSecurityGroupIngress() {
	svc := redshift.New(session.New())

	params := &redshift.AuthorizeClusterSecurityGroupIngressInput{
		ClusterSecurityGroupName: aws.String("String"), // Required
		CIDRIP:                  aws.String("String"),
		EC2SecurityGroupName:    aws.String("String"),
		EC2SecurityGroupOwnerId: aws.String("String"),
	}
	resp, err := svc.AuthorizeClusterSecurityGroupIngress(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_AuthorizeSnapshotAccess() {
	svc := redshift.New(session.New())

	params := &redshift.AuthorizeSnapshotAccessInput{
		AccountWithRestoreAccess:  aws.String("String"), // Required
		SnapshotIdentifier:        aws.String("String"), // Required
		SnapshotClusterIdentifier: aws.String("String"),
	}
	resp, err := svc.AuthorizeSnapshotAccess(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CopyClusterSnapshot() {
	svc := redshift.New(session.New())

	params := &redshift.CopyClusterSnapshotInput{
		SourceSnapshotIdentifier:        aws.String("String"), // Required
		TargetSnapshotIdentifier:        aws.String("String"), // Required
		SourceSnapshotClusterIdentifier: aws.String("String"),
	}
	resp, err := svc.CopyClusterSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateCluster() {
	svc := redshift.New(session.New())

	params := &redshift.CreateClusterInput{
		ClusterIdentifier:                aws.String("String"), // Required
		MasterUserPassword:               aws.String("String"), // Required
		MasterUsername:                   aws.String("String"), // Required
		NodeType:                         aws.String("String"), // Required
		AllowVersionUpgrade:              aws.Bool(true),
		AutomatedSnapshotRetentionPeriod: aws.Int64(1),
		AvailabilityZone:                 aws.String("String"),
		ClusterParameterGroupName:        aws.String("String"),
		ClusterSecurityGroups: []*string{
			aws.String("String"), // Required
			// More values...
		},
		ClusterSubnetGroupName:         aws.String("String"),
		ClusterType:                    aws.String("String"),
		ClusterVersion:                 aws.String("String"),
		DBName:                         aws.String("String"),
		ElasticIp:                      aws.String("String"),
		Encrypted:                      aws.Bool(true),
		HsmClientCertificateIdentifier: aws.String("String"),
		HsmConfigurationIdentifier:     aws.String("String"),
		KmsKeyId:                       aws.String("String"),
		NumberOfNodes:                  aws.Int64(1),
		Port:                           aws.Int64(1),
		PreferredMaintenanceWindow: aws.String("String"),
		PubliclyAccessible:         aws.Bool(true),
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
		VpcSecurityGroupIds: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.CreateCluster(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateClusterParameterGroup() {
	svc := redshift.New(session.New())

	params := &redshift.CreateClusterParameterGroupInput{
		Description:          aws.String("String"), // Required
		ParameterGroupFamily: aws.String("String"), // Required
		ParameterGroupName:   aws.String("String"), // Required
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateClusterParameterGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateClusterSecurityGroup() {
	svc := redshift.New(session.New())

	params := &redshift.CreateClusterSecurityGroupInput{
		ClusterSecurityGroupName: aws.String("String"), // Required
		Description:              aws.String("String"), // Required
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateClusterSecurityGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateClusterSnapshot() {
	svc := redshift.New(session.New())

	params := &redshift.CreateClusterSnapshotInput{
		ClusterIdentifier:  aws.String("String"), // Required
		SnapshotIdentifier: aws.String("String"), // Required
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateClusterSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateClusterSubnetGroup() {
	svc := redshift.New(session.New())

	params := &redshift.CreateClusterSubnetGroupInput{
		ClusterSubnetGroupName: aws.String("String"), // Required
		Description:            aws.String("String"), // Required
		SubnetIds: []*string{ // Required
			aws.String("String"), // Required
			// More values...
		},
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateClusterSubnetGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateEventSubscription() {
	svc := redshift.New(session.New())

	params := &redshift.CreateEventSubscriptionInput{
		SnsTopicArn:      aws.String("String"), // Required
		SubscriptionName: aws.String("String"), // Required
		Enabled:          aws.Bool(true),
		EventCategories: []*string{
			aws.String("String"), // Required
			// More values...
		},
		Severity: aws.String("String"),
		SourceIds: []*string{
			aws.String("String"), // Required
			// More values...
		},
		SourceType: aws.String("String"),
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateEventSubscription(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateHsmClientCertificate() {
	svc := redshift.New(session.New())

	params := &redshift.CreateHsmClientCertificateInput{
		HsmClientCertificateIdentifier: aws.String("String"), // Required
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateHsmClientCertificate(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateHsmConfiguration() {
	svc := redshift.New(session.New())

	params := &redshift.CreateHsmConfigurationInput{
		Description:                aws.String("String"), // Required
		HsmConfigurationIdentifier: aws.String("String"), // Required
		HsmIpAddress:               aws.String("String"), // Required
		HsmPartitionName:           aws.String("String"), // Required
		HsmPartitionPassword:       aws.String("String"), // Required
		HsmServerPublicCertificate: aws.String("String"), // Required
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateHsmConfiguration(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateSnapshotCopyGrant() {
	svc := redshift.New(session.New())

	params := &redshift.CreateSnapshotCopyGrantInput{
		SnapshotCopyGrantName: aws.String("String"), // Required
		KmsKeyId:              aws.String("String"),
		Tags: []*redshift.Tag{
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateSnapshotCopyGrant(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_CreateTags() {
	svc := redshift.New(session.New())

	params := &redshift.CreateTagsInput{
		ResourceName: aws.String("String"), // Required
		Tags: []*redshift.Tag{ // Required
			{ // Required
				Key:   aws.String("String"),
				Value: aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.CreateTags(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteCluster() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteClusterInput{
		ClusterIdentifier:              aws.String("String"), // Required
		FinalClusterSnapshotIdentifier: aws.String("String"),
		SkipFinalClusterSnapshot:       aws.Bool(true),
	}
	resp, err := svc.DeleteCluster(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteClusterParameterGroup() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteClusterParameterGroupInput{
		ParameterGroupName: aws.String("String"), // Required
	}
	resp, err := svc.DeleteClusterParameterGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteClusterSecurityGroup() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteClusterSecurityGroupInput{
		ClusterSecurityGroupName: aws.String("String"), // Required
	}
	resp, err := svc.DeleteClusterSecurityGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteClusterSnapshot() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteClusterSnapshotInput{
		SnapshotIdentifier:        aws.String("String"), // Required
		SnapshotClusterIdentifier: aws.String("String"),
	}
	resp, err := svc.DeleteClusterSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteClusterSubnetGroup() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteClusterSubnetGroupInput{
		ClusterSubnetGroupName: aws.String("String"), // Required
	}
	resp, err := svc.DeleteClusterSubnetGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteEventSubscription() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteEventSubscriptionInput{
		SubscriptionName: aws.String("String"), // Required
	}
	resp, err := svc.DeleteEventSubscription(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteHsmClientCertificate() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteHsmClientCertificateInput{
		HsmClientCertificateIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DeleteHsmClientCertificate(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteHsmConfiguration() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteHsmConfigurationInput{
		HsmConfigurationIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DeleteHsmConfiguration(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteSnapshotCopyGrant() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteSnapshotCopyGrantInput{
		SnapshotCopyGrantName: aws.String("String"), // Required
	}
	resp, err := svc.DeleteSnapshotCopyGrant(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DeleteTags() {
	svc := redshift.New(session.New())

	params := &redshift.DeleteTagsInput{
		ResourceName: aws.String("String"), // Required
		TagKeys: []*string{ // Required
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DeleteTags(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterParameterGroups() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterParameterGroupsInput{
		Marker:             aws.String("String"),
		MaxRecords:         aws.Int64(1),
		ParameterGroupName: aws.String("String"),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeClusterParameterGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterParameters() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterParametersInput{
		ParameterGroupName: aws.String("String"), // Required
		Marker:             aws.String("String"),
		MaxRecords:         aws.Int64(1),
		Source:             aws.String("String"),
	}
	resp, err := svc.DescribeClusterParameters(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterSecurityGroups() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterSecurityGroupsInput{
		ClusterSecurityGroupName: aws.String("String"),
		Marker:     aws.String("String"),
		MaxRecords: aws.Int64(1),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeClusterSecurityGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterSnapshots() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterSnapshotsInput{
		ClusterIdentifier:  aws.String("String"),
		EndTime:            aws.Time(time.Now()),
		Marker:             aws.String("String"),
		MaxRecords:         aws.Int64(1),
		OwnerAccount:       aws.String("String"),
		SnapshotIdentifier: aws.String("String"),
		SnapshotType:       aws.String("String"),
		StartTime:          aws.Time(time.Now()),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeClusterSnapshots(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterSubnetGroups() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterSubnetGroupsInput{
		ClusterSubnetGroupName: aws.String("String"),
		Marker:                 aws.String("String"),
		MaxRecords:             aws.Int64(1),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeClusterSubnetGroups(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusterVersions() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClusterVersionsInput{
		ClusterParameterGroupFamily: aws.String("String"),
		ClusterVersion:              aws.String("String"),
		Marker:                      aws.String("String"),
		MaxRecords:                  aws.Int64(1),
	}
	resp, err := svc.DescribeClusterVersions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeClusters() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeClustersInput{
		ClusterIdentifier: aws.String("String"),
		Marker:            aws.String("String"),
		MaxRecords:        aws.Int64(1),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeClusters(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeDefaultClusterParameters() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeDefaultClusterParametersInput{
		ParameterGroupFamily: aws.String("String"), // Required
		Marker:               aws.String("String"),
		MaxRecords:           aws.Int64(1),
	}
	resp, err := svc.DescribeDefaultClusterParameters(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeEventCategories() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeEventCategoriesInput{
		SourceType: aws.String("String"),
	}
	resp, err := svc.DescribeEventCategories(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeEventSubscriptions() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeEventSubscriptionsInput{
		Marker:           aws.String("String"),
		MaxRecords:       aws.Int64(1),
		SubscriptionName: aws.String("String"),
	}
	resp, err := svc.DescribeEventSubscriptions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeEvents() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeEventsInput{
		Duration:         aws.Int64(1),
		EndTime:          aws.Time(time.Now()),
		Marker:           aws.String("String"),
		MaxRecords:       aws.Int64(1),
		SourceIdentifier: aws.String("String"),
		SourceType:       aws.String("SourceType"),
		StartTime:        aws.Time(time.Now()),
	}
	resp, err := svc.DescribeEvents(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeHsmClientCertificates() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeHsmClientCertificatesInput{
		HsmClientCertificateIdentifier: aws.String("String"),
		Marker:     aws.String("String"),
		MaxRecords: aws.Int64(1),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeHsmClientCertificates(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeHsmConfigurations() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeHsmConfigurationsInput{
		HsmConfigurationIdentifier: aws.String("String"),
		Marker:     aws.String("String"),
		MaxRecords: aws.Int64(1),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeHsmConfigurations(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeLoggingStatus() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeLoggingStatusInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DescribeLoggingStatus(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeOrderableClusterOptions() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeOrderableClusterOptionsInput{
		ClusterVersion: aws.String("String"),
		Marker:         aws.String("String"),
		MaxRecords:     aws.Int64(1),
		NodeType:       aws.String("String"),
	}
	resp, err := svc.DescribeOrderableClusterOptions(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeReservedNodeOfferings() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeReservedNodeOfferingsInput{
		Marker:                 aws.String("String"),
		MaxRecords:             aws.Int64(1),
		ReservedNodeOfferingId: aws.String("String"),
	}
	resp, err := svc.DescribeReservedNodeOfferings(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeReservedNodes() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeReservedNodesInput{
		Marker:         aws.String("String"),
		MaxRecords:     aws.Int64(1),
		ReservedNodeId: aws.String("String"),
	}
	resp, err := svc.DescribeReservedNodes(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeResize() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeResizeInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DescribeResize(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeSnapshotCopyGrants() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeSnapshotCopyGrantsInput{
		Marker:                aws.String("String"),
		MaxRecords:            aws.Int64(1),
		SnapshotCopyGrantName: aws.String("String"),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeSnapshotCopyGrants(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DescribeTags() {
	svc := redshift.New(session.New())

	params := &redshift.DescribeTagsInput{
		Marker:       aws.String("String"),
		MaxRecords:   aws.Int64(1),
		ResourceName: aws.String("String"),
		ResourceType: aws.String("String"),
		TagKeys: []*string{
			aws.String("String"), // Required
			// More values...
		},
		TagValues: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeTags(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DisableLogging() {
	svc := redshift.New(session.New())

	params := &redshift.DisableLoggingInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DisableLogging(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_DisableSnapshotCopy() {
	svc := redshift.New(session.New())

	params := &redshift.DisableSnapshotCopyInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.DisableSnapshotCopy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_EnableLogging() {
	svc := redshift.New(session.New())

	params := &redshift.EnableLoggingInput{
		BucketName:        aws.String("String"), // Required
		ClusterIdentifier: aws.String("String"), // Required
		S3KeyPrefix:       aws.String("String"),
	}
	resp, err := svc.EnableLogging(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_EnableSnapshotCopy() {
	svc := redshift.New(session.New())

	params := &redshift.EnableSnapshotCopyInput{
		ClusterIdentifier:     aws.String("String"), // Required
		DestinationRegion:     aws.String("String"), // Required
		RetentionPeriod:       aws.Int64(1),
		SnapshotCopyGrantName: aws.String("String"),
	}
	resp, err := svc.EnableSnapshotCopy(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ModifyCluster() {
	svc := redshift.New(session.New())

	params := &redshift.ModifyClusterInput{
		ClusterIdentifier:                aws.String("String"), // Required
		AllowVersionUpgrade:              aws.Bool(true),
		AutomatedSnapshotRetentionPeriod: aws.Int64(1),
		ClusterParameterGroupName:        aws.String("String"),
		ClusterSecurityGroups: []*string{
			aws.String("String"), // Required
			// More values...
		},
		ClusterType:                    aws.String("String"),
		ClusterVersion:                 aws.String("String"),
		HsmClientCertificateIdentifier: aws.String("String"),
		HsmConfigurationIdentifier:     aws.String("String"),
		MasterUserPassword:             aws.String("String"),
		NewClusterIdentifier:           aws.String("String"),
		NodeType:                       aws.String("String"),
		NumberOfNodes:                  aws.Int64(1),
		PreferredMaintenanceWindow:     aws.String("String"),
		VpcSecurityGroupIds: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.ModifyCluster(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ModifyClusterParameterGroup() {
	svc := redshift.New(session.New())

	params := &redshift.ModifyClusterParameterGroupInput{
		ParameterGroupName: aws.String("String"), // Required
		Parameters: []*redshift.Parameter{ // Required
			{ // Required
				AllowedValues:        aws.String("String"),
				ApplyType:            aws.String("ParameterApplyType"),
				DataType:             aws.String("String"),
				Description:          aws.String("String"),
				IsModifiable:         aws.Bool(true),
				MinimumEngineVersion: aws.String("String"),
				ParameterName:        aws.String("String"),
				ParameterValue:       aws.String("String"),
				Source:               aws.String("String"),
			},
			// More values...
		},
	}
	resp, err := svc.ModifyClusterParameterGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ModifyClusterSubnetGroup() {
	svc := redshift.New(session.New())

	params := &redshift.ModifyClusterSubnetGroupInput{
		ClusterSubnetGroupName: aws.String("String"), // Required
		SubnetIds: []*string{ // Required
			aws.String("String"), // Required
			// More values...
		},
		Description: aws.String("String"),
	}
	resp, err := svc.ModifyClusterSubnetGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ModifyEventSubscription() {
	svc := redshift.New(session.New())

	params := &redshift.ModifyEventSubscriptionInput{
		SubscriptionName: aws.String("String"), // Required
		Enabled:          aws.Bool(true),
		EventCategories: []*string{
			aws.String("String"), // Required
			// More values...
		},
		Severity:    aws.String("String"),
		SnsTopicArn: aws.String("String"),
		SourceIds: []*string{
			aws.String("String"), // Required
			// More values...
		},
		SourceType: aws.String("String"),
	}
	resp, err := svc.ModifyEventSubscription(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ModifySnapshotCopyRetentionPeriod() {
	svc := redshift.New(session.New())

	params := &redshift.ModifySnapshotCopyRetentionPeriodInput{
		ClusterIdentifier: aws.String("String"), // Required
		RetentionPeriod:   aws.Int64(1),         // Required
	}
	resp, err := svc.ModifySnapshotCopyRetentionPeriod(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_PurchaseReservedNodeOffering() {
	svc := redshift.New(session.New())

	params := &redshift.PurchaseReservedNodeOfferingInput{
		ReservedNodeOfferingId: aws.String("String"), // Required
		NodeCount:              aws.Int64(1),
	}
	resp, err := svc.PurchaseReservedNodeOffering(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_RebootCluster() {
	svc := redshift.New(session.New())

	params := &redshift.RebootClusterInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.RebootCluster(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_ResetClusterParameterGroup() {
	svc := redshift.New(session.New())

	params := &redshift.ResetClusterParameterGroupInput{
		ParameterGroupName: aws.String("String"), // Required
		Parameters: []*redshift.Parameter{
			{ // Required
				AllowedValues:        aws.String("String"),
				ApplyType:            aws.String("ParameterApplyType"),
				DataType:             aws.String("String"),
				Description:          aws.String("String"),
				IsModifiable:         aws.Bool(true),
				MinimumEngineVersion: aws.String("String"),
				ParameterName:        aws.String("String"),
				ParameterValue:       aws.String("String"),
				Source:               aws.String("String"),
			},
			// More values...
		},
		ResetAllParameters: aws.Bool(true),
	}
	resp, err := svc.ResetClusterParameterGroup(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_RestoreFromClusterSnapshot() {
	svc := redshift.New(session.New())

	params := &redshift.RestoreFromClusterSnapshotInput{
		ClusterIdentifier:                aws.String("String"), // Required
		SnapshotIdentifier:               aws.String("String"), // Required
		AllowVersionUpgrade:              aws.Bool(true),
		AutomatedSnapshotRetentionPeriod: aws.Int64(1),
		AvailabilityZone:                 aws.String("String"),
		ClusterParameterGroupName:        aws.String("String"),
		ClusterSecurityGroups: []*string{
			aws.String("String"), // Required
			// More values...
		},
		ClusterSubnetGroupName:         aws.String("String"),
		ElasticIp:                      aws.String("String"),
		HsmClientCertificateIdentifier: aws.String("String"),
		HsmConfigurationIdentifier:     aws.String("String"),
		KmsKeyId:                       aws.String("String"),
		NodeType:                       aws.String("String"),
		OwnerAccount:                   aws.String("String"),
		Port:                           aws.Int64(1),
		PreferredMaintenanceWindow: aws.String("String"),
		PubliclyAccessible:         aws.Bool(true),
		SnapshotClusterIdentifier:  aws.String("String"),
		VpcSecurityGroupIds: []*string{
			aws.String("String"), // Required
			// More values...
		},
	}
	resp, err := svc.RestoreFromClusterSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_RevokeClusterSecurityGroupIngress() {
	svc := redshift.New(session.New())

	params := &redshift.RevokeClusterSecurityGroupIngressInput{
		ClusterSecurityGroupName: aws.String("String"), // Required
		CIDRIP:                  aws.String("String"),
		EC2SecurityGroupName:    aws.String("String"),
		EC2SecurityGroupOwnerId: aws.String("String"),
	}
	resp, err := svc.RevokeClusterSecurityGroupIngress(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_RevokeSnapshotAccess() {
	svc := redshift.New(session.New())

	params := &redshift.RevokeSnapshotAccessInput{
		AccountWithRestoreAccess:  aws.String("String"), // Required
		SnapshotIdentifier:        aws.String("String"), // Required
		SnapshotClusterIdentifier: aws.String("String"),
	}
	resp, err := svc.RevokeSnapshotAccess(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleRedshift_RotateEncryptionKey() {
	svc := redshift.New(session.New())

	params := &redshift.RotateEncryptionKeyInput{
		ClusterIdentifier: aws.String("String"), // Required
	}
	resp, err := svc.RotateEncryptionKey(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
