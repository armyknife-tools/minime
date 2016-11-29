// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package cloudfront

import (
	"github.com/aws/aws-sdk-go/private/waiter"
)

// WaitUntilDistributionDeployed uses the CloudFront API operation
// GetDistribution to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *CloudFront) WaitUntilDistributionDeployed(input *GetDistributionInput) error {
	waiterCfg := waiter.Config{
		Operation:   "GetDistribution",
		Delay:       60,
		MaxAttempts: 25,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "path",
				Argument: "Distribution.Status",
				Expected: "Deployed",
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

// WaitUntilInvalidationCompleted uses the CloudFront API operation
// GetInvalidation to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *CloudFront) WaitUntilInvalidationCompleted(input *GetInvalidationInput) error {
	waiterCfg := waiter.Config{
		Operation:   "GetInvalidation",
		Delay:       20,
		MaxAttempts: 30,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "path",
				Argument: "Invalidation.Status",
				Expected: "Completed",
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

// WaitUntilStreamingDistributionDeployed uses the CloudFront API operation
// GetStreamingDistribution to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *CloudFront) WaitUntilStreamingDistributionDeployed(input *GetStreamingDistributionInput) error {
	waiterCfg := waiter.Config{
		Operation:   "GetStreamingDistribution",
		Delay:       60,
		MaxAttempts: 25,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "path",
				Argument: "StreamingDistribution.Status",
				Expected: "Deployed",
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
