package cloud

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	version "github.com/hashicorp/go-version"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/internal/backend"
	"github.com/hashicorp/terraform/internal/configs/configschema"
	"github.com/hashicorp/terraform/internal/states/remote"
	"github.com/hashicorp/terraform/internal/states/statemgr"
	"github.com/hashicorp/terraform/internal/terraform"
	"github.com/hashicorp/terraform/internal/tfdiags"
	tfversion "github.com/hashicorp/terraform/version"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	backendLocal "github.com/hashicorp/terraform/internal/backend/local"
)

const (
	defaultHostname    = "app.terraform.io"
	defaultParallelism = 10
	tfeServiceID       = "tfe.v2"
)

// Cloud is an implementation of EnhancedBackend in service of the Terraform Cloud/Enterprise
// integration for Terraform CLI. This backend is not intended to be surfaced at the user level and
// is instead an implementation detail of cloud.Cloud.
type Cloud struct {
	// CLI and Colorize control the CLI output. If CLI is nil then no CLI
	// output will be done. If CLIColor is nil then no coloring will be done.
	CLI      cli.Ui
	CLIColor *colorstring.Colorize

	// ContextOpts are the base context options to set when initializing a
	// new Terraform context. Many of these will be overridden or merged by
	// Operation. See Operation for more details.
	ContextOpts *terraform.ContextOpts

	// client is the Terraform Cloud/Enterprise API client.
	client *tfe.Client

	// lastRetry is set to the last time a request was retried.
	lastRetry time.Time

	// hostname of Terraform Cloud or Terraform Enterprise
	hostname string

	// organization is the organization that contains the target workspaces.
	organization string

	// WorkspaceMapping contains strategies for mapping CLI workspaces in the working directory
	// to remote Terraform Cloud workspaces.
	WorkspaceMapping WorkspaceMapping

	// services is used for service discovery
	services *disco.Disco

	// local allows local operations, where Terraform Cloud serves as a state storage backend.
	local backend.Enhanced

	// forceLocal, if true, will force the use of the local backend.
	forceLocal bool

	// opLock locks operations
	opLock sync.Mutex

	// ignoreVersionConflict, if true, will disable the requirement that the
	// local Terraform version matches the remote workspace's configured
	// version. This will also cause VerifyWorkspaceTerraformVersion to return
	// a warning diagnostic instead of an error.
	ignoreVersionConflict bool
}

var _ backend.Backend = (*Cloud)(nil)
var _ backend.Enhanced = (*Cloud)(nil)
var _ backend.Local = (*Cloud)(nil)

// New creates a new initialized cloud backend.
func New(services *disco.Disco) *Cloud {
	return &Cloud{
		services: services,
	}
}

// ConfigSchema implements backend.Enhanced.
func (b *Cloud) ConfigSchema() *configschema.Block {
	return &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"hostname": {
				Type:        cty.String,
				Optional:    true,
				Description: schemaDescriptionHostname,
			},
			"organization": {
				Type:        cty.String,
				Required:    true,
				Description: schemaDescriptionOrganization,
			},
			"token": {
				Type:        cty.String,
				Optional:    true,
				Description: schemaDescriptionToken,
			},
		},

		BlockTypes: map[string]*configschema.NestedBlock{
			"workspaces": {
				Block: configschema.Block{
					Attributes: map[string]*configschema.Attribute{
						"name": {
							Type:        cty.String,
							Optional:    true,
							Description: schemaDescriptionName,
						},
						"prefix": {
							Type:        cty.String,
							Optional:    true,
							Description: schemaDescriptionPrefix,
						},
						"tags": {
							Type:        cty.Set(cty.String),
							Optional:    true,
							Description: schemaDescriptionTags,
						},
					},
				},
				Nesting: configschema.NestingSingle,
			},
		},
	}
}

// PrepareConfig implements backend.Backend.
func (b *Cloud) PrepareConfig(obj cty.Value) (cty.Value, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	if obj.IsNull() {
		return obj, diags
	}

	if val := obj.GetAttr("organization"); val.IsNull() || val.AsString() == "" {
		diags = diags.Append(invalidOrganizationConfigMissingValue)
	}

	WorkspaceMapping := WorkspaceMapping{}
	if workspaces := obj.GetAttr("workspaces"); !workspaces.IsNull() {
		if val := workspaces.GetAttr("name"); !val.IsNull() {
			WorkspaceMapping.Name = val.AsString()
		}
		if val := workspaces.GetAttr("prefix"); !val.IsNull() {
			WorkspaceMapping.Prefix = val.AsString()
		}
		if val := workspaces.GetAttr("tags"); !val.IsNull() {
			err := gocty.FromCtyValue(val, &WorkspaceMapping.Tags)
			if err != nil {
				log.Panicf("An unxpected error occurred: %s", err)
			}
		}
	}

	switch WorkspaceMapping.Strategy() {
	// Make sure have a workspace mapping strategy present
	case WorkspaceNoneStrategy:
		diags = diags.Append(invalidWorkspaceConfigMissingValues)
	// Make sure that only one of workspace name or a prefix is configured.
	case WorkspaceInvalidStrategy:
		diags = diags.Append(invalidWorkspaceConfigMisconfiguration)
	}

	return obj, diags
}

// Configure implements backend.Enhanced.
func (b *Cloud) Configure(obj cty.Value) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics
	if obj.IsNull() {
		return diags
	}

	diagErr := b.setConfigurationFields(obj)
	if diagErr.HasErrors() {
		return diagErr
	}

	// Discover the service URL to confirm that it provides the Terraform Cloud/Enterprise API
	service, err := b.discover()

	// Check for errors before we continue.
	if err != nil {
		diags = diags.Append(tfdiags.AttributeValue(
			tfdiags.Error,
			strings.ToUpper(err.Error()[:1])+err.Error()[1:],
			"", // no description is needed here, the error is clear
			cty.Path{cty.GetAttrStep{Name: "hostname"}},
		))
		return diags
	}

	// Retrieve the token for this host as configured in the credentials
	// section of the CLI Config File.
	token, err := b.token()
	if err != nil {
		diags = diags.Append(tfdiags.AttributeValue(
			tfdiags.Error,
			strings.ToUpper(err.Error()[:1])+err.Error()[1:],
			"", // no description is needed here, the error is clear
			cty.Path{cty.GetAttrStep{Name: "hostname"}},
		))
		return diags
	}

	// Get the token from the config if no token was configured for this
	// host in credentials section of the CLI Config File.
	if token == "" {
		if val := obj.GetAttr("token"); !val.IsNull() {
			token = val.AsString()
		}
	}

	// Return an error if we still don't have a token at this point.
	if token == "" {
		loginCommand := "terraform login"
		if b.hostname != defaultHostname {
			loginCommand = loginCommand + " " + b.hostname
		}
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Required token could not be found",
			fmt.Sprintf(
				"Run the following command to generate a token for %s:\n    %s",
				b.hostname,
				loginCommand,
			),
		))
		return diags
	}

	cfg := &tfe.Config{
		Address:      service.String(),
		BasePath:     service.Path,
		Token:        token,
		Headers:      make(http.Header),
		RetryLogHook: b.retryLogHook,
	}

	// Set the version header to the current version.
	cfg.Headers.Set(tfversion.Header, tfversion.Version)

	// Create the TFC/E API client.
	b.client, err = tfe.NewClient(cfg)
	if err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Failed to create the Terraform Enterprise client",
			fmt.Sprintf(
				`Encountered an unexpected error while creating the `+
					`Terraform Enterprise client: %s.`, err,
			),
		))
		return diags
	}

	// Check if the organization exists by reading its entitlements.
	entitlements, err := b.client.Organizations.Entitlements(context.Background(), b.organization)
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			err = fmt.Errorf("organization %q at host %s not found.\n\n"+
				"Please ensure that the organization and hostname are correct "+
				"and that your API token for %s is valid.",
				b.organization, b.hostname, b.hostname)
		}
		diags = diags.Append(tfdiags.AttributeValue(
			tfdiags.Error,
			fmt.Sprintf("Failed to read organization %q at host %s", b.organization, b.hostname),
			fmt.Sprintf("Encountered an unexpected error while reading the "+
				"organization settings: %s", err),
			cty.Path{cty.GetAttrStep{Name: "organization"}},
		))
		return diags
	}

	// Configure a local backend for when we need to run operations locally.
	b.local = backendLocal.NewWithBackend(b)
	b.forceLocal = b.forceLocal || !entitlements.Operations

	// Enable retries for server errors as the backend is now fully configured.
	b.client.RetryServerErrors(true)

	return diags
}

func (b *Cloud) setConfigurationFields(obj cty.Value) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	// Get the hostname.
	if val := obj.GetAttr("hostname"); !val.IsNull() && val.AsString() != "" {
		b.hostname = val.AsString()
	} else {
		b.hostname = defaultHostname
	}

	// Get the organization.
	if val := obj.GetAttr("organization"); !val.IsNull() {
		b.organization = val.AsString()
	}

	// Get the workspaces configuration block and retrieve the
	// default workspace name and prefix.
	if workspaces := obj.GetAttr("workspaces"); !workspaces.IsNull() {

		// PrepareConfig checks that you cannot set both of these.
		if val := workspaces.GetAttr("name"); !val.IsNull() {
			b.WorkspaceMapping.Name = val.AsString()
		}
		if val := workspaces.GetAttr("prefix"); !val.IsNull() {
			b.WorkspaceMapping.Prefix = val.AsString()
		}
		if val := workspaces.GetAttr("tags"); !val.IsNull() {
			var tags []string
			err := gocty.FromCtyValue(val, &tags)
			if err != nil {
				log.Panicf("An unxpected error occurred: %s", err)
			}

			b.WorkspaceMapping.Tags = tags
		}
	}

	// Determine if we are forced to use the local backend.
	b.forceLocal = os.Getenv("TF_FORCE_LOCAL_BACKEND") != ""

	return diags
}

// discover the TFC/E API service URL and version constraints.
func (b *Cloud) discover() (*url.URL, error) {
	hostname, err := svchost.ForComparison(b.hostname)
	if err != nil {
		return nil, err
	}

	host, err := b.services.Discover(hostname)
	if err != nil {
		return nil, err
	}

	service, err := host.ServiceURL(tfeServiceID)
	// Return the error, unless its a disco.ErrVersionNotSupported error.
	if _, ok := err.(*disco.ErrVersionNotSupported); !ok && err != nil {
		return nil, err
	}

	return service, err
}

// token returns the token for this host as configured in the credentials
// section of the CLI Config File. If no token was configured, an empty
// string will be returned instead.
func (b *Cloud) token() (string, error) {
	hostname, err := svchost.ForComparison(b.hostname)
	if err != nil {
		return "", err
	}
	creds, err := b.services.CredentialsForHost(hostname)
	if err != nil {
		log.Printf("[WARN] Failed to get credentials for %s: %s (ignoring)", b.hostname, err)
		return "", nil
	}
	if creds != nil {
		return creds.Token(), nil
	}
	return "", nil
}

// retryLogHook is invoked each time a request is retried allowing the
// backend to log any connection issues to prevent data loss.
func (b *Cloud) retryLogHook(attemptNum int, resp *http.Response) {
	if b.CLI != nil {
		// Ignore the first retry to make sure any delayed output will
		// be written to the console before we start logging retries.
		//
		// The retry logic in the TFE client will retry both rate limited
		// requests and server errors, but in the cloud backend we only
		// care about server errors so we ignore rate limit (429) errors.
		if attemptNum == 0 || (resp != nil && resp.StatusCode == 429) {
			// Reset the last retry time.
			b.lastRetry = time.Now()
			return
		}

		if attemptNum == 1 {
			b.CLI.Output(b.Colorize().Color(strings.TrimSpace(initialRetryError)))
		} else {
			b.CLI.Output(b.Colorize().Color(strings.TrimSpace(
				fmt.Sprintf(repeatedRetryError, time.Since(b.lastRetry).Round(time.Second)))))
		}
	}
}

// Workspaces implements backend.Enhanced, returning a filtered list of workspace names according to
// the workspace mapping strategy configured.
func (b *Cloud) Workspaces() ([]string, error) {
	// Create a slice to contain all the names.
	var names []string

	// If configured for a single workspace, return that exact name only.  The StateMgr for this
	// backend will automatically create the remote workspace if it does not yet exist.
	if b.WorkspaceMapping.Strategy() == WorkspaceNameStrategy {
		names = append(names, b.WorkspaceMapping.Name)
		return names, nil
	}

	// Otherwise, multiple workspaces are being mapped. Query Terraform Cloud for all the remote
	// workspaces by the provided mapping strategy.
	options := tfe.WorkspaceListOptions{}
	switch b.WorkspaceMapping.Strategy() {
	case WorkspacePrefixStrategy:
		options.Search = tfe.String(b.WorkspaceMapping.Prefix)
	case WorkspaceTagsStrategy:
		taglist := strings.Join(b.WorkspaceMapping.Tags, ",")
		options.Tags = &taglist
	}

	for {
		wl, err := b.client.Workspaces.List(context.Background(), b.organization, options)
		if err != nil {
			return nil, err
		}

		for _, w := range wl.Items {
			switch b.WorkspaceMapping.Strategy() {
			case WorkspacePrefixStrategy:
				if strings.HasPrefix(w.Name, b.WorkspaceMapping.Prefix) {
					names = append(names, strings.TrimPrefix(w.Name, b.WorkspaceMapping.Prefix))
					continue
				}
			default:
				// Pass-through. The "prefix" strategy is naive and does
				// client-side filtering, but for tags and any other future
				// strategy this filtering should be left to the API.
				names = append(names, w.Name)
			}
		}

		// Exit the loop when we've seen all pages.
		if wl.CurrentPage >= wl.TotalPages {
			break
		}

		// Update the page number to get the next page.
		options.PageNumber = wl.NextPage
	}

	// Sort the result so we have consistent output.
	sort.StringSlice(names).Sort()

	return names, nil
}

// DeleteWorkspace implements backend.Enhanced.
func (b *Cloud) DeleteWorkspace(name string) error {
	if name == backend.DefaultStateName {
		return backend.ErrDefaultWorkspaceNotSupported
	}

	if b.WorkspaceMapping.Strategy() == WorkspaceNameStrategy {
		return backend.ErrWorkspacesNotSupported
	}

	// Configure the remote workspace name.
	switch {
	case b.WorkspaceMapping.Strategy() == WorkspacePrefixStrategy && !strings.HasPrefix(name, b.WorkspaceMapping.Prefix):
		name = b.WorkspaceMapping.Prefix + name
	}

	client := &remoteClient{
		client:       b.client,
		organization: b.organization,
		workspace: &tfe.Workspace{
			Name: name,
		},
	}

	return client.Delete()
}

// StateMgr implements backend.Enhanced.
func (b *Cloud) StateMgr(name string) (statemgr.Full, error) {
	if name == backend.DefaultStateName {
		return nil, backend.ErrDefaultWorkspaceNotSupported
	}

	if b.WorkspaceMapping.Strategy() == WorkspaceNameStrategy && name != b.WorkspaceMapping.Name {
		return nil, backend.ErrWorkspacesNotSupported
	}

	// If the prefix strategy is used, translate the local name to the TFC workspace name.
	if b.WorkspaceMapping.Strategy() == WorkspacePrefixStrategy {
		name = b.WorkspaceMapping.Prefix + name
	}

	workspace, err := b.client.Workspaces.Read(context.Background(), b.organization, name)
	if err != nil && err != tfe.ErrResourceNotFound {
		return nil, fmt.Errorf("Failed to retrieve workspace %s: %v", name, err)
	}

	if err == tfe.ErrResourceNotFound {
		options := tfe.WorkspaceCreateOptions{
			Name: tfe.String(name),
		}

		var tags []*tfe.Tag
		for _, tag := range b.WorkspaceMapping.Tags {
			t := tfe.Tag{Name: tag}
			tags = append(tags, &t)
		}
		options.Tags = tags

		options.TerraformVersion = tfe.String(tfversion.String())

		workspace, err = b.client.Workspaces.Create(context.Background(), b.organization, options)
		if err != nil {
			return nil, fmt.Errorf("Error creating workspace %s: %v", name, err)
		}
	}

	// This is a fallback error check. Most code paths should use other
	// mechanisms to check the version, then set the ignoreVersionConflict
	// field to true. This check is only in place to ensure that we don't
	// accidentally upgrade state with a new code path, and the version check
	// logic is coarser and simpler.
	if !b.ignoreVersionConflict {
		wsv := workspace.TerraformVersion
		// Explicitly ignore the pseudo-version "latest" here, as it will cause
		// plan and apply to always fail.
		if wsv != tfversion.String() && wsv != "latest" {
			return nil, fmt.Errorf("Remote workspace Terraform version %q does not match local Terraform version %q", workspace.TerraformVersion, tfversion.String())
		}
	}

	client := &remoteClient{
		client:       b.client,
		organization: b.organization,
		workspace:    workspace,

		// This is optionally set during Terraform Enterprise runs.
		runID: os.Getenv("TFE_RUN_ID"),
	}

	return &remote.State{Client: client}, nil
}

// Operation implements backend.Enhanced.
func (b *Cloud) Operation(ctx context.Context, op *backend.Operation) (*backend.RunningOperation, error) {
	name := op.Workspace

	// If the prefix strategy is used, translate the local name to the TFC workspace name.
	if b.WorkspaceMapping.Strategy() == WorkspacePrefixStrategy {
		name = b.WorkspaceMapping.Prefix + op.Workspace
	}

	// Retrieve the workspace for this operation.
	w, err := b.client.Workspaces.Read(ctx, b.organization, name)
	if err != nil {
		switch err {
		case context.Canceled:
			return nil, err
		case tfe.ErrResourceNotFound:
			return nil, fmt.Errorf(
				"workspace %s not found\n\n"+
					"For security, Terraform Cloud returns '404 Not Found' responses for resources\n"+
					"for resources that a user doesn't have access to, in addition to resources that\n"+
					"do not exist. If the resource does exist, please check the permissions of the provided token.",
				name,
			)
		default:
			return nil, fmt.Errorf(
				"Terraform Cloud returned an unexpected error:\n\n%s",
				err,
			)
		}
	}

	// Terraform remote version conflicts are not a concern for operations. We
	// are in one of three states:
	//
	// - Running remotely, in which case the local version is irrelevant;
	// - Workspace configured for local operations, in which case the remote
	//   version is meaningless;
	// - Forcing local operations, which should only happen in the Terraform Cloud worker, in
	//   which case the Terraform versions by definition match.
	b.IgnoreVersionConflict()

	// Check if we need to use the local backend to run the operation.
	if b.forceLocal || !w.Operations {
		// Record that we're forced to run operations locally to allow the
		// command package UI to operate correctly
		b.forceLocal = true
		return b.local.Operation(ctx, op)
	}

	// Set the remote workspace name.
	op.Workspace = w.Name

	// Determine the function to call for our operation
	var f func(context.Context, context.Context, *backend.Operation, *tfe.Workspace) (*tfe.Run, error)
	switch op.Type {
	case backend.OperationTypePlan:
		f = b.opPlan
	case backend.OperationTypeApply:
		f = b.opApply
	case backend.OperationTypeRefresh:
		return nil, fmt.Errorf(
			"\n\nThe \"refresh\" operation is not supported when using Terraform Cloud. " +
				"Use \"terraform apply -refresh-only\" instead.")
	default:
		return nil, fmt.Errorf(
			"\n\nTerraform Cloud does not support the %q operation.", op.Type)
	}

	// Lock
	b.opLock.Lock()

	// Build our running operation
	// the runninCtx is only used to block until the operation returns.
	runningCtx, done := context.WithCancel(context.Background())
	runningOp := &backend.RunningOperation{
		Context:   runningCtx,
		PlanEmpty: true,
	}

	// stopCtx wraps the context passed in, and is used to signal a graceful Stop.
	stopCtx, stop := context.WithCancel(ctx)
	runningOp.Stop = stop

	// cancelCtx is used to cancel the operation immediately, usually
	// indicating that the process is exiting.
	cancelCtx, cancel := context.WithCancel(context.Background())
	runningOp.Cancel = cancel

	// Do it.
	go func() {
		defer done()
		defer stop()
		defer cancel()

		defer b.opLock.Unlock()

		r, opErr := f(stopCtx, cancelCtx, op, w)
		if opErr != nil && opErr != context.Canceled {
			var diags tfdiags.Diagnostics
			diags = diags.Append(opErr)
			op.ReportResult(runningOp, diags)
			return
		}

		if r == nil && opErr == context.Canceled {
			runningOp.Result = backend.OperationFailure
			return
		}

		if r != nil {
			// Retrieve the run to get its current status.
			r, err := b.client.Runs.Read(cancelCtx, r.ID)
			if err != nil {
				var diags tfdiags.Diagnostics
				diags = diags.Append(generalError("Failed to retrieve run", err))
				op.ReportResult(runningOp, diags)
				return
			}

			// Record if there are any changes.
			runningOp.PlanEmpty = !r.HasChanges

			if opErr == context.Canceled {
				if err := b.cancel(cancelCtx, op, r); err != nil {
					var diags tfdiags.Diagnostics
					diags = diags.Append(generalError("Failed to retrieve run", err))
					op.ReportResult(runningOp, diags)
					return
				}
			}

			if r.Status == tfe.RunCanceled || r.Status == tfe.RunErrored {
				runningOp.Result = backend.OperationFailure
			}
		}
	}()

	// Return the running operation.
	return runningOp, nil
}

func (b *Cloud) cancel(cancelCtx context.Context, op *backend.Operation, r *tfe.Run) error {
	if r.Actions.IsCancelable {
		// Only ask if the remote operation should be canceled
		// if the auto approve flag is not set.
		if !op.AutoApprove {
			v, err := op.UIIn.Input(cancelCtx, &terraform.InputOpts{
				Id:          "cancel",
				Query:       "\nDo you want to cancel the remote operation?",
				Description: "Only 'yes' will be accepted to cancel.",
			})
			if err != nil {
				return generalError("Failed asking to cancel", err)
			}
			if v != "yes" {
				if b.CLI != nil {
					b.CLI.Output(b.Colorize().Color(strings.TrimSpace(operationNotCanceled)))
				}
				return nil
			}
		} else {
			if b.CLI != nil {
				// Insert a blank line to separate the ouputs.
				b.CLI.Output("")
			}
		}

		// Try to cancel the remote operation.
		err := b.client.Runs.Cancel(cancelCtx, r.ID, tfe.RunCancelOptions{})
		if err != nil {
			return generalError("Failed to cancel run", err)
		}
		if b.CLI != nil {
			b.CLI.Output(b.Colorize().Color(strings.TrimSpace(operationCanceled)))
		}
	}

	return nil
}

// IgnoreVersionConflict allows commands to disable the fall-back check that
// the local Terraform version matches the remote workspace's configured
// Terraform version. This should be called by commands where this check is
// unnecessary, such as those performing remote operations, or read-only
// operations. It will also be called if the user uses a command-line flag to
// override this check.
func (b *Cloud) IgnoreVersionConflict() {
	b.ignoreVersionConflict = true
}

// VerifyWorkspaceTerraformVersion compares the local Terraform version against
// the workspace's configured Terraform version. If they are equal, this means
// that there are no compatibility concerns, so it returns no diagnostics.
//
// If the versions differ,
func (b *Cloud) VerifyWorkspaceTerraformVersion(workspaceName string) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	workspace, err := b.getRemoteWorkspace(context.Background(), workspaceName)
	if err != nil {
		// If the workspace doesn't exist, there can be no compatibility
		// problem, so we can return. This is most likely to happen when
		// migrating state from a local backend to a new workspace.
		if err == tfe.ErrResourceNotFound {
			return nil
		}

		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Error looking up workspace",
			fmt.Sprintf("Workspace read failed: %s", err),
		))
		return diags
	}

	// If the workspace has the pseudo-version "latest", all bets are off. We
	// cannot reasonably determine what the intended Terraform version is, so
	// we'll skip version verification.
	if workspace.TerraformVersion == "latest" {
		return nil
	}

	// If the workspace has remote operations disabled, the remote Terraform
	// version is effectively meaningless, so we'll skip version verification.
	if !workspace.Operations {
		return nil
	}

	remoteVersion, err := version.NewSemver(workspace.TerraformVersion)
	if err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Error looking up workspace",
			fmt.Sprintf("Invalid Terraform version: %s", err),
		))
		return diags
	}

	v014 := version.Must(version.NewSemver("0.14.0"))
	if tfversion.SemVer.LessThan(v014) || remoteVersion.LessThan(v014) {
		// Versions of Terraform prior to 0.14.0 will refuse to load state files
		// written by a newer version of Terraform, even if it is only a patch
		// level difference. As a result we require an exact match.
		if tfversion.SemVer.Equal(remoteVersion) {
			return diags
		}
	}
	if tfversion.SemVer.GreaterThanOrEqual(v014) && remoteVersion.GreaterThanOrEqual(v014) {
		// Versions of Terraform after 0.14.0 should be compatible with each
		// other.  At the time this code was written, the only constraints we
		// are aware of are:
		//
		// - 0.14.0 is guaranteed to be compatible with versions up to but not
		//   including 1.2.0
		v120 := version.Must(version.NewSemver("1.2.0"))
		if tfversion.SemVer.LessThan(v120) && remoteVersion.LessThan(v120) {
			return diags
		}
		// - Any new Terraform state version will require at least minor patch
		//   increment, so x.y.* will always be compatible with each other
		tfvs := tfversion.SemVer.Segments64()
		rwvs := remoteVersion.Segments64()
		if len(tfvs) == 3 && len(rwvs) == 3 && tfvs[0] == rwvs[0] && tfvs[1] == rwvs[1] {
			return diags
		}
	}

	// Even if ignoring version conflicts, it may still be useful to call this
	// method and warn the user about a mismatch between the local and remote
	// Terraform versions.
	severity := tfdiags.Error
	if b.ignoreVersionConflict {
		severity = tfdiags.Warning
	}

	suggestion := " If you're sure you want to upgrade the state, you can force Terraform to continue using the -ignore-remote-version flag. This may result in an unusable workspace."
	if b.ignoreVersionConflict {
		suggestion = ""
	}
	diags = diags.Append(tfdiags.Sourceless(
		severity,
		"Terraform version mismatch",
		fmt.Sprintf(
			"The local Terraform version (%s) does not match the configured version for remote workspace %s/%s (%s).%s",
			tfversion.String(),
			b.organization,
			workspace.Name,
			workspace.TerraformVersion,
			suggestion,
		),
	))

	return diags
}

func (b *Cloud) IsLocalOperations() bool {
	return b.forceLocal
}

// Colorize returns the Colorize structure that can be used for colorizing
// output. This is guaranteed to always return a non-nil value and so useful
// as a helper to wrap any potentially colored strings.
//
// TODO SvH: Rename this back to Colorize as soon as we can pass -no-color.
//lint:ignore U1000 see above todo
func (b *Cloud) cliColorize() *colorstring.Colorize {
	if b.CLIColor != nil {
		return b.CLIColor
	}

	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: true,
	}
}

type WorkspaceMapping struct {
	Name   string
	Prefix string
	Tags   []string
}

type workspaceStrategy string

const (
	WorkspaceTagsStrategy    workspaceStrategy = "tags"
	WorkspaceNameStrategy    workspaceStrategy = "name"
	WorkspacePrefixStrategy  workspaceStrategy = "prefix"
	WorkspaceNoneStrategy    workspaceStrategy = "none"
	WorkspaceInvalidStrategy workspaceStrategy = "invalid"
)

func (wm WorkspaceMapping) Strategy() workspaceStrategy {
	switch {
	case len(wm.Tags) > 0 && wm.Name == "" && wm.Prefix == "":
		return WorkspaceTagsStrategy
	case len(wm.Tags) == 0 && wm.Name != "" && wm.Prefix == "":
		return WorkspaceNameStrategy
	case len(wm.Tags) == 0 && wm.Name == "" && wm.Prefix != "":
		return WorkspacePrefixStrategy
	case len(wm.Tags) == 0 && wm.Name == "" && wm.Prefix == "":
		return WorkspaceNoneStrategy
	default:
		// Any other combination is invalid as each strategy is mutually exclusive
		return WorkspaceInvalidStrategy
	}
}

func generalError(msg string, err error) error {
	var diags tfdiags.Diagnostics

	if urlErr, ok := err.(*url.Error); ok {
		err = urlErr.Err
	}

	switch err {
	case context.Canceled:
		return err
	case tfe.ErrResourceNotFound:
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			fmt.Sprintf("%s: %v", msg, err),
			"For security, Terraform Cloud returns '404 Not Found' responses for resources\n"+
				"for resources that a user doesn't have access to, in addition to resources that\n"+
				"do not exist. If the resource does exist, please check the permissions of the provided token.",
		))
		return diags.Err()
	default:
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			fmt.Sprintf("%s: %v", msg, err),
			`Terraform Cloud returned an unexpected error. Sometimes `+
				`this is caused by network connection problems, in which case you could retry `+
				`the command. If the issue persists please open a support ticket to get help `+
				`resolving the problem.`,
		))
		return diags.Err()
	}
}

// The newline in this error is to make it look good in the CLI!
const initialRetryError = `
[reset][yellow]There was an error connecting to Terraform Cloud. Please do not exit
Terraform to prevent data loss! Trying to restore the connection...
[reset]
`

const repeatedRetryError = `
[reset][yellow]Still trying to restore the connection... (%s elapsed)[reset]
`

const operationCanceled = `
[reset][red]The remote operation was successfully cancelled.[reset]
`

const operationNotCanceled = `
[reset][red]The remote operation was not cancelled.[reset]
`

var (
	workspaceConfigurationHelp = fmt.Sprintf(
		`The 'workspaces' block configures how Terraform CLI maps its workspaces for this single
configuration to workspaces within a Terraform Cloud organization. Three strategies are available:

[bold]tags[reset] - %s

[bold]name[reset] - %s

[bold]prefix[reset] - %s`, schemaDescriptionTags, schemaDescriptionName, schemaDescriptionPrefix)

	schemaDescriptionHostname = `The Terraform Enterprise hostname to connect to. This optional argument defaults to app.terraform.io
for use with Terraform Cloud.`

	schemaDescriptionOrganization = `The name of the organization containing the targeted workspace(s).`

	schemaDescriptionToken = `The token used to authenticate with Terraform Cloud/Enterprise. Typically this argument should not
be set, and 'terraform login' used instead; your credentials will then be fetched from your CLI
configuration file or configured credential helper.`

	schemaDescriptionTags = `A set of tags used to select remote Terraform Cloud workspaces to be used for this single
configuration.  New workspaces will automatically be tagged with these tag values.  Generally, this
is the primary and recommended strategy to use.  This option conflicts with "prefix" and "name".`

	schemaDescriptionName = `The name of a single Terraform Cloud workspace to be used with this configuration When configured
only the specified workspace can be used. This option conflicts with "tags" and "prefix".`

	schemaDescriptionPrefix = `DEPRECATED. A name prefix used to select remote Terraform Cloud to be used for this single configuration. New
workspaces will automatically be prefixed with this prefix. This option conflicts with "tags" and "name".`
)
