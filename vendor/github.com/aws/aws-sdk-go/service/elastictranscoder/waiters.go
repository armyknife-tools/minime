// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package elastictranscoder

import (
	"github.com/aws/aws-sdk-go/private/waiter"
)

// WaitUntilJobComplete uses the Amazon Elastic Transcoder API operation
// ReadJob to wait for a condition to be met before returning.
// If the condition is not meet within the max attempt window an error will
// be returned.
func (c *ElasticTranscoder) WaitUntilJobComplete(input *ReadJobInput) error {
	waiterCfg := waiter.Config{
		Operation:   "ReadJob",
		Delay:       30,
		MaxAttempts: 120,
		Acceptors: []waiter.WaitAcceptor{
			{
				State:    "success",
				Matcher:  "path",
				Argument: "Job.Status",
				Expected: "Complete",
			},
			{
				State:    "failure",
				Matcher:  "path",
				Argument: "Job.Status",
				Expected: "Canceled",
			},
			{
				State:    "failure",
				Matcher:  "path",
				Argument: "Job.Status",
				Expected: "Error",
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
