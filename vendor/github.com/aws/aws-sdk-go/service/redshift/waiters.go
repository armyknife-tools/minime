// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package redshift

import (
	"github.com/aws/aws-sdk-go/private/waiter"
)

// WaitUntilClusterAvailable uses the Amazon Redshift API operation
// DescribeClusters to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *Redshift) WaitUntilClusterAvailable(input *DescribeClustersInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeClusters",
		Delay:       60,
		MaxAttempts: 30,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "pathAll",
				Argument: "Clusters[].ClusterStatus",
				Expected: "available",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Clusters[].ClusterStatus",
				Expected: "deleting",
			},
			{
				State:    "retry",
				Matcher:  "error",
				Argument: "",
				Expected: "ClusterNotFound",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

// WaitUntilClusterDeleted uses the Amazon Redshift API operation
// DescribeClusters to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *Redshift) WaitUntilClusterDeleted(input *DescribeClustersInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeClusters",
		Delay:       60,
		MaxAttempts: 30,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "error",
				Argument: "",
				Expected: "ClusterNotFound",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Clusters[].ClusterStatus",
				Expected: "creating",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Clusters[].ClusterStatus",
				Expected: "modifying",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

// WaitUntilClusterRestored uses the Amazon Redshift API operation
// DescribeClusters to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *Redshift) WaitUntilClusterRestored(input *DescribeClustersInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeClusters",
		Delay:       60,
		MaxAttempts: 30,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "pathAll",
				Argument: "Clusters[].RestoreStatus.Status",
				Expected: "completed",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Clusters[].ClusterStatus",
				Expected: "deleting",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}

// WaitUntilSnapshotAvailable uses the Amazon Redshift API operation
// DescribeClusterSnapshots to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *Redshift) WaitUntilSnapshotAvailable(input *DescribeClusterSnapshotsInput) error {
	waiterCfg := waiter.Config{
		Operation:   "DescribeClusterSnapshots",
		Delay:       15,
		MaxAttempts: 20,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "pathAll",
				Argument: "Snapshots[].Status",
				Expected: "available",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Snapshots[].Status",
				Expected: "failed",
			},
			{
				State:    "failure",
				Matcher:  "pathAny",
				Argument: "Snapshots[].Status",
				Expected: "deleted",
			},
		},
	}

	w := waiter.Waiter{
		Client: c,
		Input:  input,
		Config: waiterCfg,
	}
	return w.Wait()
}
