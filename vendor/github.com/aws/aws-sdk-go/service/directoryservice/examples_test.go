// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package directoryservice_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/directoryservice"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleDirectoryService_ConnectDirectory() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.ConnectDirectoryInput{
		ConnectSettings: &directoryservice.DirectoryConnectSettings{ // Required
			CustomerDnsIps: []*string{ // Required
				aws.String("IpAddr"), // Required
				// More values...
			},
			CustomerUserName: aws.String("UserName"), // Required
			SubnetIds: []*string{ // Required
				aws.String("SubnetId"), // Required
				// More values...
			},
			VpcId: aws.String("VpcId"), // Required
		},
		Name:        aws.String("DirectoryName"),   // Required
		Password:    aws.String("ConnectPassword"), // Required
		Size:        aws.String("DirectorySize"),   // Required
		Description: aws.String("Description"),
		ShortName:   aws.String("DirectoryShortName"),
	}
	resp, err := svc.ConnectDirectory(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateAlias() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateAliasInput{
		Alias:       aws.String("AliasName"),   // Required
		DirectoryId: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.CreateAlias(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateComputer() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateComputerInput{
		ComputerName: aws.String("ComputerName"),     // Required
		DirectoryId:  aws.String("DirectoryId"),      // Required
		Password:     aws.String("ComputerPassword"), // Required
		ComputerAttributes: []*directoryservice.Attribute{
			{ // Required
				Name:  aws.String("AttributeName"),
				Value: aws.String("AttributeValue"),
			},
			// More values...
		},
		OrganizationalUnitDistinguishedName: aws.String("OrganizationalUnitDN"),
	}
	resp, err := svc.CreateComputer(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateDirectory() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateDirectoryInput{
		Name:        aws.String("DirectoryName"), // Required
		Password:    aws.String("Password"),      // Required
		Size:        aws.String("DirectorySize"), // Required
		Description: aws.String("Description"),
		ShortName:   aws.String("DirectoryShortName"),
		VpcSettings: &directoryservice.DirectoryVpcSettings{
			SubnetIds: []*string{ // Required
				aws.String("SubnetId"), // Required
				// More values...
			},
			VpcId: aws.String("VpcId"), // Required
		},
	}
	resp, err := svc.CreateDirectory(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateMicrosoftAD() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateMicrosoftADInput{
		Name:     aws.String("DirectoryName"), // Required
		Password: aws.String("Password"),      // Required
		VpcSettings: &directoryservice.DirectoryVpcSettings{ // Required
			SubnetIds: []*string{ // Required
				aws.String("SubnetId"), // Required
				// More values...
			},
			VpcId: aws.String("VpcId"), // Required
		},
		Description: aws.String("Description"),
		ShortName:   aws.String("DirectoryShortName"),
	}
	resp, err := svc.CreateMicrosoftAD(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateSnapshot() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateSnapshotInput{
		DirectoryId: aws.String("DirectoryId"), // Required
		Name:        aws.String("SnapshotName"),
	}
	resp, err := svc.CreateSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_CreateTrust() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.CreateTrustInput{
		DirectoryId:      aws.String("DirectoryId"),      // Required
		RemoteDomainName: aws.String("RemoteDomainName"), // Required
		TrustDirection:   aws.String("TrustDirection"),   // Required
		TrustPassword:    aws.String("TrustPassword"),    // Required
		TrustType:        aws.String("TrustType"),
	}
	resp, err := svc.CreateTrust(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DeleteDirectory() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DeleteDirectoryInput{
		DirectoryId: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.DeleteDirectory(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DeleteSnapshot() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DeleteSnapshotInput{
		SnapshotId: aws.String("SnapshotId"), // Required
	}
	resp, err := svc.DeleteSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DeleteTrust() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DeleteTrustInput{
		TrustId: aws.String("TrustId"), // Required
	}
	resp, err := svc.DeleteTrust(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DescribeDirectories() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DescribeDirectoriesInput{
		DirectoryIds: []*string{
			aws.String("DirectoryId"), // Required
			// More values...
		},
		Limit:     aws.Int64(1),
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.DescribeDirectories(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DescribeSnapshots() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DescribeSnapshotsInput{
		DirectoryId: aws.String("DirectoryId"),
		Limit:       aws.Int64(1),
		NextToken:   aws.String("NextToken"),
		SnapshotIds: []*string{
			aws.String("SnapshotId"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeSnapshots(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DescribeTrusts() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DescribeTrustsInput{
		DirectoryId: aws.String("DirectoryId"),
		Limit:       aws.Int64(1),
		NextToken:   aws.String("NextToken"),
		TrustIds: []*string{
			aws.String("TrustId"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeTrusts(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DisableRadius() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DisableRadiusInput{
		DirectoryId: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.DisableRadius(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_DisableSso() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.DisableSsoInput{
		DirectoryId: aws.String("DirectoryId"), // Required
		Password:    aws.String("ConnectPassword"),
		UserName:    aws.String("UserName"),
	}
	resp, err := svc.DisableSso(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_EnableRadius() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.EnableRadiusInput{
		DirectoryId: aws.String("DirectoryId"), // Required
		RadiusSettings: &directoryservice.RadiusSettings{ // Required
			AuthenticationProtocol: aws.String("RadiusAuthenticationProtocol"),
			DisplayLabel:           aws.String("RadiusDisplayLabel"),
			RadiusPort:             aws.Int64(1),
			RadiusRetries:          aws.Int64(1),
			RadiusServers: []*string{
				aws.String("Server"), // Required
				// More values...
			},
			RadiusTimeout:   aws.Int64(1),
			SharedSecret:    aws.String("RadiusSharedSecret"),
			UseSameUsername: aws.Bool(true),
		},
	}
	resp, err := svc.EnableRadius(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_EnableSso() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.EnableSsoInput{
		DirectoryId: aws.String("DirectoryId"), // Required
		Password:    aws.String("ConnectPassword"),
		UserName:    aws.String("UserName"),
	}
	resp, err := svc.EnableSso(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_GetDirectoryLimits() {
	svc := directoryservice.New(session.New())

	var params *directoryservice.GetDirectoryLimitsInput
	resp, err := svc.GetDirectoryLimits(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_GetSnapshotLimits() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.GetSnapshotLimitsInput{
		DirectoryId: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.GetSnapshotLimits(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_RestoreFromSnapshot() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.RestoreFromSnapshotInput{
		SnapshotId: aws.String("SnapshotId"), // Required
	}
	resp, err := svc.RestoreFromSnapshot(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_UpdateRadius() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.UpdateRadiusInput{
		DirectoryId: aws.String("DirectoryId"), // Required
		RadiusSettings: &directoryservice.RadiusSettings{ // Required
			AuthenticationProtocol: aws.String("RadiusAuthenticationProtocol"),
			DisplayLabel:           aws.String("RadiusDisplayLabel"),
			RadiusPort:             aws.Int64(1),
			RadiusRetries:          aws.Int64(1),
			RadiusServers: []*string{
				aws.String("Server"), // Required
				// More values...
			},
			RadiusTimeout:   aws.Int64(1),
			SharedSecret:    aws.String("RadiusSharedSecret"),
			UseSameUsername: aws.Bool(true),
		},
	}
	resp, err := svc.UpdateRadius(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func ExampleDirectoryService_VerifyTrust() {
	svc := directoryservice.New(session.New())

	params := &directoryservice.VerifyTrustInput{
		TrustId: aws.String("TrustId"), // Required
	}
	resp, err := svc.VerifyTrust(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}
