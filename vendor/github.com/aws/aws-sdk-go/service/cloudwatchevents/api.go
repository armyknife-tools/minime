// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package cloudwatchevents provides a client for Amazon CloudWatch Events.
package cloudwatchevents

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/request"
)

const opDeleteRule = "DeleteRule"

// DeleteRuleRequest generates a request for the DeleteRule operation.
func (c *CloudWatchEvents) DeleteRuleRequest(input *DeleteRuleInput) (req *request.Request, output *DeleteRuleOutput) {
	op := &request.Operation{
		Name:       opDeleteRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DeleteRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &DeleteRuleOutput{}
	req.Data = output
	return
}

// Deletes a rule. You must remove all targets from a rule using RemoveTargets
// before you can delete the rule.
//
//  Note: When you make a change with this action, incoming events might still
// continue to match to the deleted rule. Please allow a short period of time
// for changes to take effect.
func (c *CloudWatchEvents) DeleteRule(input *DeleteRuleInput) (*DeleteRuleOutput, error) {
	req, out := c.DeleteRuleRequest(input)
	err := req.Send()
	return out, err
}

const opDescribeRule = "DescribeRule"

// DescribeRuleRequest generates a request for the DescribeRule operation.
func (c *CloudWatchEvents) DescribeRuleRequest(input *DescribeRuleInput) (req *request.Request, output *DescribeRuleOutput) {
	op := &request.Operation{
		Name:       opDescribeRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DescribeRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &DescribeRuleOutput{}
	req.Data = output
	return
}

// Describes the details of the specified rule.
func (c *CloudWatchEvents) DescribeRule(input *DescribeRuleInput) (*DescribeRuleOutput, error) {
	req, out := c.DescribeRuleRequest(input)
	err := req.Send()
	return out, err
}

const opDisableRule = "DisableRule"

// DisableRuleRequest generates a request for the DisableRule operation.
func (c *CloudWatchEvents) DisableRuleRequest(input *DisableRuleInput) (req *request.Request, output *DisableRuleOutput) {
	op := &request.Operation{
		Name:       opDisableRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DisableRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &DisableRuleOutput{}
	req.Data = output
	return
}

// Disables a rule. A disabled rule won't match any events, and won't self-trigger
// if it has a schedule expression.
//
//  Note: When you make a change with this action, incoming events might still
// continue to match to the disabled rule. Please allow a short period of time
// for changes to take effect.
func (c *CloudWatchEvents) DisableRule(input *DisableRuleInput) (*DisableRuleOutput, error) {
	req, out := c.DisableRuleRequest(input)
	err := req.Send()
	return out, err
}

const opEnableRule = "EnableRule"

// EnableRuleRequest generates a request for the EnableRule operation.
func (c *CloudWatchEvents) EnableRuleRequest(input *EnableRuleInput) (req *request.Request, output *EnableRuleOutput) {
	op := &request.Operation{
		Name:       opEnableRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &EnableRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &EnableRuleOutput{}
	req.Data = output
	return
}

// Enables a rule. If the rule does not exist, the operation fails.
//
//  Note: When you make a change with this action, incoming events might not
// immediately start matching to a newly enabled rule. Please allow a short
// period of time for changes to take effect.
func (c *CloudWatchEvents) EnableRule(input *EnableRuleInput) (*EnableRuleOutput, error) {
	req, out := c.EnableRuleRequest(input)
	err := req.Send()
	return out, err
}

const opListRuleNamesByTarget = "ListRuleNamesByTarget"

// ListRuleNamesByTargetRequest generates a request for the ListRuleNamesByTarget operation.
func (c *CloudWatchEvents) ListRuleNamesByTargetRequest(input *ListRuleNamesByTargetInput) (req *request.Request, output *ListRuleNamesByTargetOutput) {
	op := &request.Operation{
		Name:       opListRuleNamesByTarget,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ListRuleNamesByTargetInput{}
	}

	req = c.newRequest(op, input, output)
	output = &ListRuleNamesByTargetOutput{}
	req.Data = output
	return
}

// Lists the names of the rules that the given target is put to. Using this
// action, you can find out which of the rules in Amazon CloudWatch Events can
// invoke a specific target in your account. If you have more rules in your
// account than the given limit, the results will be paginated. In that case,
// use the next token returned in the response and repeat the ListRulesByTarget
// action until the NextToken in the response is returned as null.
func (c *CloudWatchEvents) ListRuleNamesByTarget(input *ListRuleNamesByTargetInput) (*ListRuleNamesByTargetOutput, error) {
	req, out := c.ListRuleNamesByTargetRequest(input)
	err := req.Send()
	return out, err
}

const opListRules = "ListRules"

// ListRulesRequest generates a request for the ListRules operation.
func (c *CloudWatchEvents) ListRulesRequest(input *ListRulesInput) (req *request.Request, output *ListRulesOutput) {
	op := &request.Operation{
		Name:       opListRules,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ListRulesInput{}
	}

	req = c.newRequest(op, input, output)
	output = &ListRulesOutput{}
	req.Data = output
	return
}

// Lists the Amazon CloudWatch Events rules in your account. You can either
// list all the rules or you can provide a prefix to match to the rule names.
// If you have more rules in your account than the given limit, the results
// will be paginated. In that case, use the next token returned in the response
// and repeat the ListRules action until the NextToken in the response is returned
// as null.
func (c *CloudWatchEvents) ListRules(input *ListRulesInput) (*ListRulesOutput, error) {
	req, out := c.ListRulesRequest(input)
	err := req.Send()
	return out, err
}

const opListTargetsByRule = "ListTargetsByRule"

// ListTargetsByRuleRequest generates a request for the ListTargetsByRule operation.
func (c *CloudWatchEvents) ListTargetsByRuleRequest(input *ListTargetsByRuleInput) (req *request.Request, output *ListTargetsByRuleOutput) {
	op := &request.Operation{
		Name:       opListTargetsByRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &ListTargetsByRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &ListTargetsByRuleOutput{}
	req.Data = output
	return
}

// Lists of targets assigned to the rule.
func (c *CloudWatchEvents) ListTargetsByRule(input *ListTargetsByRuleInput) (*ListTargetsByRuleOutput, error) {
	req, out := c.ListTargetsByRuleRequest(input)
	err := req.Send()
	return out, err
}

const opPutEvents = "PutEvents"

// PutEventsRequest generates a request for the PutEvents operation.
func (c *CloudWatchEvents) PutEventsRequest(input *PutEventsInput) (req *request.Request, output *PutEventsOutput) {
	op := &request.Operation{
		Name:       opPutEvents,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &PutEventsInput{}
	}

	req = c.newRequest(op, input, output)
	output = &PutEventsOutput{}
	req.Data = output
	return
}

// Sends custom events to Amazon CloudWatch Events so that they can be matched
// to rules.
func (c *CloudWatchEvents) PutEvents(input *PutEventsInput) (*PutEventsOutput, error) {
	req, out := c.PutEventsRequest(input)
	err := req.Send()
	return out, err
}

const opPutRule = "PutRule"

// PutRuleRequest generates a request for the PutRule operation.
func (c *CloudWatchEvents) PutRuleRequest(input *PutRuleInput) (req *request.Request, output *PutRuleOutput) {
	op := &request.Operation{
		Name:       opPutRule,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &PutRuleInput{}
	}

	req = c.newRequest(op, input, output)
	output = &PutRuleOutput{}
	req.Data = output
	return
}

// Creates or updates a rule. Rules are enabled by default, or based on value
// of the State parameter. You can disable a rule using DisableRule.
//
//  Note: When you make a change with this action, incoming events might not
// immediately start matching to new or updated rules. Please allow a short
// period of time for changes to take effect.
//
// A rule must contain at least an EventPattern or ScheduleExpression. Rules
// with EventPatterns are triggered when a matching event is observed. Rules
// with ScheduleExpressions self-trigger based on the given schedule. A rule
// can have both an EventPattern and a ScheduleExpression, in which case the
// rule will trigger on matching events as well as on a schedule.
//
//  Note: Most services in AWS treat : or / as the same character in Amazon
// Resource Names (ARNs). However, CloudWatch Events uses an exact match in
// event patterns and rules. Be sure to use the correct ARN characters when
// creating event patterns so that they match the ARN syntax in the event you
// want to match.
func (c *CloudWatchEvents) PutRule(input *PutRuleInput) (*PutRuleOutput, error) {
	req, out := c.PutRuleRequest(input)
	err := req.Send()
	return out, err
}

const opPutTargets = "PutTargets"

// PutTargetsRequest generates a request for the PutTargets operation.
func (c *CloudWatchEvents) PutTargetsRequest(input *PutTargetsInput) (req *request.Request, output *PutTargetsOutput) {
	op := &request.Operation{
		Name:       opPutTargets,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &PutTargetsInput{}
	}

	req = c.newRequest(op, input, output)
	output = &PutTargetsOutput{}
	req.Data = output
	return
}

// Adds target(s) to a rule. Updates the target(s) if they are already associated
// with the role. In other words, if there is already a target with the given
// target ID, then the target associated with that ID is updated.
//
//  Note: When you make a change with this action, when the associated rule
// triggers, new or updated targets might not be immediately invoked. Please
// allow a short period of time for changes to take effect.
func (c *CloudWatchEvents) PutTargets(input *PutTargetsInput) (*PutTargetsOutput, error) {
	req, out := c.PutTargetsRequest(input)
	err := req.Send()
	return out, err
}

const opRemoveTargets = "RemoveTargets"

// RemoveTargetsRequest generates a request for the RemoveTargets operation.
func (c *CloudWatchEvents) RemoveTargetsRequest(input *RemoveTargetsInput) (req *request.Request, output *RemoveTargetsOutput) {
	op := &request.Operation{
		Name:       opRemoveTargets,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &RemoveTargetsInput{}
	}

	req = c.newRequest(op, input, output)
	output = &RemoveTargetsOutput{}
	req.Data = output
	return
}

// Removes target(s) from a rule so that when the rule is triggered, those targets
// will no longer be invoked.
//
//  Note: When you make a change with this action, when the associated rule
// triggers, removed targets might still continue to be invoked. Please allow
// a short period of time for changes to take effect.
func (c *CloudWatchEvents) RemoveTargets(input *RemoveTargetsInput) (*RemoveTargetsOutput, error) {
	req, out := c.RemoveTargetsRequest(input)
	err := req.Send()
	return out, err
}

const opTestEventPattern = "TestEventPattern"

// TestEventPatternRequest generates a request for the TestEventPattern operation.
func (c *CloudWatchEvents) TestEventPatternRequest(input *TestEventPatternInput) (req *request.Request, output *TestEventPatternOutput) {
	op := &request.Operation{
		Name:       opTestEventPattern,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &TestEventPatternInput{}
	}

	req = c.newRequest(op, input, output)
	output = &TestEventPatternOutput{}
	req.Data = output
	return
}

// Tests whether an event pattern matches the provided event.
//
//  Note: Most services in AWS treat : or / as the same character in Amazon
// Resource Names (ARNs). However, CloudWatch Events uses an exact match in
// event patterns and rules. Be sure to use the correct ARN characters when
// creating event patterns so that they match the ARN syntax in the event you
// want to match.
func (c *CloudWatchEvents) TestEventPattern(input *TestEventPatternInput) (*TestEventPatternOutput, error) {
	req, out := c.TestEventPatternRequest(input)
	err := req.Send()
	return out, err
}

// Container for the parameters to the DeleteRule operation.
type DeleteRuleInput struct {
	_ struct{} `type:"structure"`

	// The name of the rule to be deleted.
	Name *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRuleInput) GoString() string {
	return s.String()
}

type DeleteRuleOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DeleteRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the DescribeRule operation.
type DescribeRuleInput struct {
	_ struct{} `type:"structure"`

	// The name of the rule you want to describe details for.
	Name *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s DescribeRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRuleInput) GoString() string {
	return s.String()
}

// The result of the DescribeRule operation.
type DescribeRuleOutput struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) associated with the rule.
	Arn *string `min:"1" type:"string"`

	// The rule's description.
	Description *string `type:"string"`

	// The event pattern.
	EventPattern *string `type:"string"`

	// The rule's name.
	Name *string `min:"1" type:"string"`

	// The Amazon Resource Name (ARN) of the IAM role associated with the rule.
	RoleArn *string `min:"1" type:"string"`

	// The scheduling expression. For example, "cron(0 20 * * ? *)", "rate(5 minutes)".
	ScheduleExpression *string `type:"string"`

	// Specifies whether the rule is enabled or disabled.
	State *string `type:"string" enum:"RuleState"`
}

// String returns the string representation
func (s DescribeRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DescribeRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the DisableRule operation.
type DisableRuleInput struct {
	_ struct{} `type:"structure"`

	// The name of the rule you want to disable.
	Name *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s DisableRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisableRuleInput) GoString() string {
	return s.String()
}

type DisableRuleOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DisableRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s DisableRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the EnableRule operation.
type EnableRuleInput struct {
	_ struct{} `type:"structure"`

	// The name of the rule that you want to enable.
	Name *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s EnableRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EnableRuleInput) GoString() string {
	return s.String()
}

type EnableRuleOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s EnableRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s EnableRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the ListRuleNamesByTarget operation.
type ListRuleNamesByTargetInput struct {
	_ struct{} `type:"structure"`

	// The maximum number of results to return.
	Limit *int64 `min:"1" type:"integer"`

	// The token returned by a previous call to indicate that there is more data
	// available.
	NextToken *string `min:"1" type:"string"`

	// The Amazon Resource Name (ARN) of the target resource that you want to list
	// the rules for.
	TargetArn *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s ListRuleNamesByTargetInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListRuleNamesByTargetInput) GoString() string {
	return s.String()
}

// The result of the ListRuleNamesByTarget operation.
type ListRuleNamesByTargetOutput struct {
	_ struct{} `type:"structure"`

	// Indicates that there are additional results to retrieve.
	NextToken *string `min:"1" type:"string"`

	// List of rules names that can invoke the given target.
	RuleNames []*string `type:"list"`
}

// String returns the string representation
func (s ListRuleNamesByTargetOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListRuleNamesByTargetOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the ListRules operation.
type ListRulesInput struct {
	_ struct{} `type:"structure"`

	// The maximum number of results to return.
	Limit *int64 `min:"1" type:"integer"`

	// The prefix matching the rule name.
	NamePrefix *string `min:"1" type:"string"`

	// The token returned by a previous call to indicate that there is more data
	// available.
	NextToken *string `min:"1" type:"string"`
}

// String returns the string representation
func (s ListRulesInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListRulesInput) GoString() string {
	return s.String()
}

// The result of the ListRules operation.
type ListRulesOutput struct {
	_ struct{} `type:"structure"`

	// Indicates that there are additional results to retrieve.
	NextToken *string `min:"1" type:"string"`

	// List of rules matching the specified criteria.
	Rules []*Rule `type:"list"`
}

// String returns the string representation
func (s ListRulesOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListRulesOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the ListTargetsByRule operation.
type ListTargetsByRuleInput struct {
	_ struct{} `type:"structure"`

	// The maximum number of results to return.
	Limit *int64 `min:"1" type:"integer"`

	// The token returned by a previous call to indicate that there is more data
	// available.
	NextToken *string `min:"1" type:"string"`

	// The name of the rule whose targets you want to list.
	Rule *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s ListTargetsByRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListTargetsByRuleInput) GoString() string {
	return s.String()
}

// The result of the ListTargetsByRule operation.
type ListTargetsByRuleOutput struct {
	_ struct{} `type:"structure"`

	// Indicates that there are additional results to retrieve.
	NextToken *string `min:"1" type:"string"`

	// Lists the targets assigned to the rule.
	Targets []*Target `type:"list"`
}

// String returns the string representation
func (s ListTargetsByRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s ListTargetsByRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the PutEvents operation.
type PutEventsInput struct {
	_ struct{} `type:"structure"`

	// The entry that defines an event in your system. You can specify several parameters
	// for the entry such as the source and type of the event, resources associated
	// with the event, and so on.
	Entries []*PutEventsRequestEntry `min:"1" type:"list" required:"true"`
}

// String returns the string representation
func (s PutEventsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutEventsInput) GoString() string {
	return s.String()
}

// The result of the PutEvents operation.
type PutEventsOutput struct {
	_ struct{} `type:"structure"`

	// A list of successfully and unsuccessfully ingested events results. If the
	// ingestion was successful, the entry will have the event ID in it. If not,
	// then the ErrorCode and ErrorMessage can be used to identify the problem with
	// the entry.
	Entries []*PutEventsResultEntry `type:"list"`

	// The number of failed entries.
	FailedEntryCount *int64 `type:"integer"`
}

// String returns the string representation
func (s PutEventsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutEventsOutput) GoString() string {
	return s.String()
}

// Contains information about the event to be used in the PutEvents action.
type PutEventsRequestEntry struct {
	_ struct{} `type:"structure"`

	// In the JSON sense, an object containing fields, which may also contain nested
	// sub-objects. No constraints are imposed on its contents.
	Detail *string `type:"string"`

	// Free-form string used to decide what fields to expect in the event detail.
	DetailType *string `type:"string"`

	// AWS resources, identified by Amazon Resource Name (ARN), which the event
	// primarily concerns. Any number, including zero, may be present.
	Resources []*string `type:"list"`

	// The source of the event.
	Source *string `type:"string"`

	// Timestamp of event, per RFC3339 (https://www.rfc-editor.org/rfc/rfc3339.txt).
	// If no timestamp is provided, the timestamp of the PutEvents call will be
	// used.
	Time *time.Time `type:"timestamp" timestampFormat:"unix"`
}

// String returns the string representation
func (s PutEventsRequestEntry) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutEventsRequestEntry) GoString() string {
	return s.String()
}

// A PutEventsResult contains a list of PutEventsResultEntry.
type PutEventsResultEntry struct {
	_ struct{} `type:"structure"`

	// The error code representing why the event submission failed on this entry.
	ErrorCode *string `type:"string"`

	// The error message explaining why the event submission failed on this entry.
	ErrorMessage *string `type:"string"`

	// The ID of the event submitted to Amazon CloudWatch Events.
	EventId *string `type:"string"`
}

// String returns the string representation
func (s PutEventsResultEntry) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutEventsResultEntry) GoString() string {
	return s.String()
}

// Container for the parameters to the PutRule operation.
type PutRuleInput struct {
	_ struct{} `type:"structure"`

	// A description of the rule.
	Description *string `type:"string"`

	// The event pattern.
	EventPattern *string `type:"string"`

	// The name of the rule that you are creating or updating.
	Name *string `min:"1" type:"string" required:"true"`

	// The Amazon Resource Name (ARN) of the IAM role associated with the rule.
	RoleArn *string `min:"1" type:"string"`

	// The scheduling expression. For example, "cron(0 20 * * ? *)", "rate(5 minutes)".
	ScheduleExpression *string `type:"string"`

	// Indicates whether the rule is enabled or disabled.
	State *string `type:"string" enum:"RuleState"`
}

// String returns the string representation
func (s PutRuleInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutRuleInput) GoString() string {
	return s.String()
}

// The result of the PutRule operation.
type PutRuleOutput struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) that identifies the rule.
	RuleArn *string `min:"1" type:"string"`
}

// String returns the string representation
func (s PutRuleOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutRuleOutput) GoString() string {
	return s.String()
}

// Container for the parameters to the PutTargets operation.
type PutTargetsInput struct {
	_ struct{} `type:"structure"`

	// The name of the rule you want to add targets to.
	Rule *string `min:"1" type:"string" required:"true"`

	// List of targets you want to update or add to the rule.
	Targets []*Target `type:"list" required:"true"`
}

// String returns the string representation
func (s PutTargetsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutTargetsInput) GoString() string {
	return s.String()
}

// The result of the PutTargets operation.
type PutTargetsOutput struct {
	_ struct{} `type:"structure"`

	// An array of failed target entries.
	FailedEntries []*PutTargetsResultEntry `type:"list"`

	// The number of failed entries.
	FailedEntryCount *int64 `type:"integer"`
}

// String returns the string representation
func (s PutTargetsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutTargetsOutput) GoString() string {
	return s.String()
}

// A PutTargetsResult contains a list of PutTargetsResultEntry.
type PutTargetsResultEntry struct {
	_ struct{} `type:"structure"`

	// The error code representing why the target submission failed on this entry.
	ErrorCode *string `type:"string"`

	// The error message explaining why the target submission failed on this entry.
	ErrorMessage *string `type:"string"`

	// The ID of the target submitted to Amazon CloudWatch Events.
	TargetId *string `min:"1" type:"string"`
}

// String returns the string representation
func (s PutTargetsResultEntry) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s PutTargetsResultEntry) GoString() string {
	return s.String()
}

// Container for the parameters to the RemoveTargets operation.
type RemoveTargetsInput struct {
	_ struct{} `type:"structure"`

	// The list of target IDs to remove from the rule.
	Ids []*string `min:"1" type:"list" required:"true"`

	// The name of the rule you want to remove targets from.
	Rule *string `min:"1" type:"string" required:"true"`
}

// String returns the string representation
func (s RemoveTargetsInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RemoveTargetsInput) GoString() string {
	return s.String()
}

// The result of the RemoveTargets operation.
type RemoveTargetsOutput struct {
	_ struct{} `type:"structure"`

	// An array of failed target entries.
	FailedEntries []*RemoveTargetsResultEntry `type:"list"`

	// The number of failed entries.
	FailedEntryCount *int64 `type:"integer"`
}

// String returns the string representation
func (s RemoveTargetsOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RemoveTargetsOutput) GoString() string {
	return s.String()
}

// The ID of the target requested to be removed from the rule by Amazon CloudWatch
// Events.
type RemoveTargetsResultEntry struct {
	_ struct{} `type:"structure"`

	// The error code representing why the target removal failed on this entry.
	ErrorCode *string `type:"string"`

	// The error message explaining why the target removal failed on this entry.
	ErrorMessage *string `type:"string"`

	// The ID of the target requested to be removed by Amazon CloudWatch Events.
	TargetId *string `min:"1" type:"string"`
}

// String returns the string representation
func (s RemoveTargetsResultEntry) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s RemoveTargetsResultEntry) GoString() string {
	return s.String()
}

// Contains information about a rule in Amazon CloudWatch Events. A ListRulesResult
// contains a list of Rules.
type Rule struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) of the rule.
	Arn *string `min:"1" type:"string"`

	// The description of the rule.
	Description *string `type:"string"`

	// The event pattern of the rule.
	EventPattern *string `type:"string"`

	// The rule's name.
	Name *string `min:"1" type:"string"`

	// The Amazon Resource Name (ARN) associated with the role that is used for
	// target invocation.
	RoleArn *string `min:"1" type:"string"`

	// The scheduling expression. For example, "cron(0 20 * * ? *)", "rate(5 minutes)".
	ScheduleExpression *string `type:"string"`

	// The rule's state.
	State *string `type:"string" enum:"RuleState"`
}

// String returns the string representation
func (s Rule) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Rule) GoString() string {
	return s.String()
}

// Targets are the resources that can be invoked when a rule is triggered. For
// example, AWS Lambda functions, Amazon Kinesis streams, and built-in targets.
//
// Input and InputPath are mutually-exclusive and optional parameters of a
// target. When a rule is triggered due to a matched event, if for a target:
//
//  Neither Input nor InputPath is specified, then the entire event is passed
// to the target in JSON form.  InputPath is specified in the form of JSONPath
// (e.g. $.detail), then only the part of the event specified in the path is
// passed to the target (e.g. only the detail part of the event is passed).
//   Input is specified in the form of a valid JSON, then the matched event
// is overridden with this constant.
type Target struct {
	_ struct{} `type:"structure"`

	// The Amazon Resource Name (ARN) associated of the target.
	Arn *string `min:"1" type:"string" required:"true"`

	// The unique target assignment ID.
	Id *string `min:"1" type:"string" required:"true"`

	// Valid JSON text passed to the target. For more information about JSON text,
	// see The JavaScript Object Notation (JSON) Data Interchange Format (http://www.rfc-editor.org/rfc/rfc7159.txt).
	Input *string `type:"string"`

	// The value of the JSONPath that is used for extracting part of the matched
	// event when passing it to the target. For more information about JSON paths,
	// see JSONPath (http://goessner.net/articles/JsonPath/).
	InputPath *string `type:"string"`
}

// String returns the string representation
func (s Target) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s Target) GoString() string {
	return s.String()
}

// Container for the parameters to the TestEventPattern operation.
type TestEventPatternInput struct {
	_ struct{} `type:"structure"`

	// The event in the JSON format to test against the event pattern.
	Event *string `type:"string" required:"true"`

	// The event pattern you want to test.
	EventPattern *string `type:"string" required:"true"`
}

// String returns the string representation
func (s TestEventPatternInput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s TestEventPatternInput) GoString() string {
	return s.String()
}

// The result of the TestEventPattern operation.
type TestEventPatternOutput struct {
	_ struct{} `type:"structure"`

	// Indicates whether the event matches the event pattern.
	Result *bool `type:"boolean"`
}

// String returns the string representation
func (s TestEventPatternOutput) String() string {
	return awsutil.Prettify(s)
}

// GoString returns the string representation
func (s TestEventPatternOutput) GoString() string {
	return s.String()
}

const (
	// @enum RuleState
	RuleStateEnabled = "ENABLED"
	// @enum RuleState
	RuleStateDisabled = "DISABLED"
)
