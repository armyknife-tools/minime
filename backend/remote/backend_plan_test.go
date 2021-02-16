package remote

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/internal/initwd"
	"github.com/hashicorp/terraform/plans/planfile"
	"github.com/hashicorp/terraform/states/statemgr"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/mitchellh/cli"
)

func testOperationPlan(t *testing.T, configDir string) (*backend.Operation, func()) {
	t.Helper()

	_, configLoader, configCleanup := initwd.MustLoadConfigForTests(t, configDir)

	return &backend.Operation{
		ConfigDir:       configDir,
		ConfigLoader:    configLoader,
		Parallelism:     defaultParallelism,
		PlanRefresh:     true,
		ShowDiagnostics: testLogDiagnostics(t),
		Type:            backend.OperationTypePlan,
	}, configCleanup
}

func testOperationPlanWithDiagnostics(t *testing.T, configDir string) (*backend.Operation, func(), func() tfdiags.Diagnostics) {
	t.Helper()

	op, cleanup := testOperationPlan(t, configDir)

	record, playback := testRecordDiagnostics(t)
	op.ShowDiagnostics = record

	return op, cleanup, playback
}

func TestRemote_planBasic(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatal("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}

	stateMgr, _ := b.StateMgr(backend.DefaultStateName)
	// An error suggests that the state was not unlocked after the operation finished
	if _, err := stateMgr.Lock(statemgr.NewLockInfo()); err != nil {
		t.Fatalf("unexpected error locking state after successful plan: %s", err.Error())
	}
}

func TestRemote_planCanceled(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	// Stop the run to simulate a Ctrl-C.
	run.Stop()

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}

	stateMgr, _ := b.StateMgr(backend.DefaultStateName)
	// An error suggests that the state was not unlocked after the operation finished
	if _, err := stateMgr.Lock(statemgr.NewLockInfo()); err != nil {
		t.Fatalf("unexpected error locking state after cancelled plan: %s", err.Error())
	}
}

func TestRemote_planLongLine(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan-long-line")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatal("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planWithoutPermissions(t *testing.T) {
	b, bCleanup := testBackendNoDefault(t)
	defer bCleanup()

	// Create a named workspace without permissions.
	w, err := b.client.Workspaces.Create(
		context.Background(),
		b.organization,
		tfe.WorkspaceCreateOptions{
			Name: tfe.String(b.prefix + "prod"),
		},
	)
	if err != nil {
		t.Fatalf("error creating named workspace: %v", err)
	}
	w.Permissions.CanQueueRun = false

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = "prod"

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "Insufficient rights to generate a plan") {
		t.Fatalf("expected a permissions error, got: %v", errOutput)
	}
}

func TestRemote_planWithParallelism(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	op.Parallelism = 3
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "parallelism values are currently not supported") {
		t.Fatalf("expected a parallelism error, got: %v", errOutput)
	}
}

func TestRemote_planWithPlan(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	op.PlanFile = &planfile.Reader{}
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "saved plan is currently not supported") {
		t.Fatalf("expected a saved plan error, got: %v", errOutput)
	}
}

func TestRemote_planWithPath(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	op.PlanOutPath = "./testdata/plan"
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "generated plan is currently not supported") {
		t.Fatalf("expected a generated plan error, got: %v", errOutput)
	}
}

func TestRemote_planWithoutRefresh(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	op.PlanRefresh = false
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "refresh is currently not supported") {
		t.Fatalf("expected a refresh error, got: %v", errOutput)
	}
}

func TestRemote_planWithTarget(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	// When the backend code creates a new run, we'll tweak it so that it
	// has a cost estimation object with the "skipped_due_to_targeting" status,
	// emulating how a real server is expected to behave in that case.
	b.client.Runs.(*mockRuns).modifyNewRun = func(client *mockClient, options tfe.RunCreateOptions, run *tfe.Run) {
		const fakeID = "fake"
		// This is the cost estimate object embedded in the run itself which
		// the backend will use to learn the ID to request from the cost
		// estimates endpoint. It's pending to simulate what a freshly-created
		// run is likely to look like.
		run.CostEstimate = &tfe.CostEstimate{
			ID:     fakeID,
			Status: "pending",
		}
		// The backend will then use the main cost estimation API to retrieve
		// the same ID indicated in the object above, where we'll then return
		// the status "skipped_due_to_targeting" to trigger the special skip
		// message in the backend output.
		client.CostEstimates.estimations[fakeID] = &tfe.CostEstimate{
			ID:     fakeID,
			Status: "skipped_due_to_targeting",
		}
	}

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	addr, _ := addrs.ParseAbsResourceStr("null_resource.foo")

	op.Targets = []addrs.Targetable{addr}
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatal("expected plan operation to succeed")
	}
	if run.PlanEmpty {
		t.Fatalf("expected plan to be non-empty")
	}

	// testBackendDefault above attached a "mock UI" to our backend, so we
	// can retrieve its non-error output via the OutputWriter in-memory buffer.
	gotOutput := b.CLI.(*cli.MockUi).OutputWriter.String()
	if wantOutput := "Not available for this plan, because it was created with the -target option."; !strings.Contains(gotOutput, wantOutput) {
		t.Errorf("missing message about skipped cost estimation\ngot:\n%s\nwant substring: %s", gotOutput, wantOutput)
	}

	// We should find a run inside the mock client that has the same
	// target address we requested above.
	runsAPI := b.client.Runs.(*mockRuns)
	if got, want := len(runsAPI.runs), 1; got != want {
		t.Fatalf("wrong number of runs in the mock client %d; want %d", got, want)
	}
	for _, run := range runsAPI.runs {
		if diff := cmp.Diff([]string{"null_resource.foo"}, run.TargetAddrs); diff != "" {
			t.Errorf("wrong TargetAddrs in the created run\n%s", diff)
		}

		if !strings.Contains(run.Message, "using -target") {
			t.Errorf("incorrect Message on the created run: %s", run.Message)
		}
	}
}

func TestRemote_planWithTargetIncompatibleAPIVersion(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan")
	defer configCleanup()

	// Set the tfe client's RemoteAPIVersion to an empty string, to mimic
	// API versions prior to 2.3.
	b.client.SetFakeRemoteAPIVersion("")

	addr, _ := addrs.ParseAbsResourceStr("null_resource.foo")

	op.Targets = []addrs.Targetable{addr}
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "Resource targeting is not supported") {
		t.Fatalf("expected a targeting error, got: %v", errOutput)
	}
}

func TestRemote_planWithVariables(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan-variables")
	defer configCleanup()

	op.Variables = testVariables(terraform.ValueFromCLIArg, "foo", "bar")
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "variables are currently not supported") {
		t.Fatalf("expected a variables error, got: %v", errOutput)
	}
}

func TestRemote_planNoConfig(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/empty")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "configuration files found") {
		t.Fatalf("expected configuration files error, got: %v", errOutput)
	}
}

func TestRemote_planNoChanges(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan-no-changes")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "No changes. Infrastructure is up-to-date.") {
		t.Fatalf("expected no changes in plan summary: %s", output)
	}
	if !strings.Contains(output, "Sentinel Result: true") {
		t.Fatalf("expected policy check result in output: %s", output)
	}
}

func TestRemote_planForceLocal(t *testing.T) {
	// Set TF_FORCE_LOCAL_BACKEND so the remote backend will use
	// the local backend with itself as embedded backend.
	if err := os.Setenv("TF_FORCE_LOCAL_BACKEND", "1"); err != nil {
		t.Fatalf("error setting environment variable TF_FORCE_LOCAL_BACKEND: %v", err)
	}
	defer os.Unsetenv("TF_FORCE_LOCAL_BACKEND")

	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("unexpected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planWithoutOperationsEntitlement(t *testing.T) {
	b, bCleanup := testBackendNoOperations(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("unexpected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planWorkspaceWithoutOperations(t *testing.T) {
	b, bCleanup := testBackendNoDefault(t)
	defer bCleanup()

	ctx := context.Background()

	// Create a named workspace that doesn't allow operations.
	_, err := b.client.Workspaces.Create(
		ctx,
		b.organization,
		tfe.WorkspaceCreateOptions{
			Name: tfe.String(b.prefix + "no-operations"),
		},
	)
	if err != nil {
		t.Fatalf("error creating named workspace: %v", err)
	}

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = "no-operations"

	run, err := b.Operation(ctx, op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("unexpected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planLockTimeout(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	ctx := context.Background()

	// Retrieve the workspace used to run this operation in.
	w, err := b.client.Workspaces.Read(ctx, b.organization, b.workspace)
	if err != nil {
		t.Fatalf("error retrieving workspace: %v", err)
	}

	// Create a new configuration version.
	c, err := b.client.ConfigurationVersions.Create(ctx, w.ID, tfe.ConfigurationVersionCreateOptions{})
	if err != nil {
		t.Fatalf("error creating configuration version: %v", err)
	}

	// Create a pending run to block this run.
	_, err = b.client.Runs.Create(ctx, tfe.RunCreateOptions{
		ConfigurationVersion: c,
		Workspace:            w,
	})
	if err != nil {
		t.Fatalf("error creating pending run: %v", err)
	}

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	input := testInput(t, map[string]string{
		"cancel":  "yes",
		"approve": "yes",
	})

	op.StateLockTimeout = 50 * time.Millisecond
	op.UIIn = input
	op.UIOut = b.CLI
	op.Workspace = backend.DefaultStateName

	_, err = b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	select {
	case <-sigint:
		// Stop redirecting SIGINT signals.
		signal.Stop(sigint)
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("expected lock timeout after 50 milliseconds, waited 200 milliseconds")
	}

	if len(input.answers) != 2 {
		t.Fatalf("expected unused answers, got: %v", input.answers)
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "Lock timeout exceeded") {
		t.Fatalf("expected lock timout error in output: %s", output)
	}
	if strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("unexpected plan summary in output: %s", output)
	}
}

func TestRemote_planDestroy(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Destroy = true
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}
}

func TestRemote_planDestroyNoConfig(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/empty")
	defer configCleanup()

	op.Destroy = true
	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}
}

func TestRemote_planWithWorkingDirectory(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	options := tfe.WorkspaceUpdateOptions{
		WorkingDirectory: tfe.String("terraform"),
	}

	// Configure the workspace to use a custom working directory.
	_, err := b.client.Workspaces.Update(context.Background(), b.organization, b.workspace, options)
	if err != nil {
		t.Fatalf("error configuring working directory: %v", err)
	}

	op, configCleanup := testOperationPlan(t, "./testdata/plan-with-working-directory/terraform")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "The remote workspace is configured to work with configuration") {
		t.Fatalf("expected working directory warning: %s", output)
	}
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planWithWorkingDirectoryFromCurrentPath(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	options := tfe.WorkspaceUpdateOptions{
		WorkingDirectory: tfe.String("terraform"),
	}

	// Configure the workspace to use a custom working directory.
	_, err := b.client.Workspaces.Update(context.Background(), b.organization, b.workspace, options)
	if err != nil {
		t.Fatalf("error configuring working directory: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current working directory: %v", err)
	}

	// We need to change into the configuration directory to make sure
	// the logic to upload the correct slug is working as expected.
	if err := os.Chdir("./testdata/plan-with-working-directory/terraform"); err != nil {
		t.Fatalf("error changing directory: %v", err)
	}
	defer os.Chdir(wd) // Make sure we change back again when were done.

	// For this test we need to give our current directory instead of the
	// full path to the configuration as we already changed directories.
	op, configCleanup := testOperationPlan(t, ".")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planCostEstimation(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan-cost-estimation")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "Resources: 1 of 1 estimated") {
		t.Fatalf("expected cost estimate result in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planPolicyPass(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan-policy-passed")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result != backend.OperationSuccess {
		t.Fatalf("operation failed: %s", b.CLI.(*cli.MockUi).ErrorWriter.String())
	}
	if run.PlanEmpty {
		t.Fatalf("expected a non-empty plan")
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "Sentinel Result: true") {
		t.Fatalf("expected policy check result in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planPolicyHardFail(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan-policy-hard-failed")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "hard failed") {
		t.Fatalf("expected a policy check error, got: %v", errOutput)
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "Sentinel Result: false") {
		t.Fatalf("expected policy check result in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planPolicySoftFail(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup, playback := testOperationPlanWithDiagnostics(t, "./testdata/plan-policy-soft-failed")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if !run.PlanEmpty {
		t.Fatalf("expected plan to be empty")
	}

	errOutput := playback().Err().Error()
	if !strings.Contains(errOutput, "soft failed") {
		t.Fatalf("expected a policy check error, got: %v", errOutput)
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "Sentinel Result: false") {
		t.Fatalf("expected policy check result in output: %s", output)
	}
	if !strings.Contains(output, "1 to add, 0 to change, 0 to destroy") {
		t.Fatalf("expected plan summary in output: %s", output)
	}
}

func TestRemote_planWithRemoteError(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan-with-error")
	defer configCleanup()

	op.Workspace = backend.DefaultStateName

	run, err := b.Operation(context.Background(), op)
	if err != nil {
		t.Fatalf("error starting operation: %v", err)
	}

	<-run.Done()
	if run.Result == backend.OperationSuccess {
		t.Fatal("expected plan operation to fail")
	}
	if run.Result.ExitStatus() != 1 {
		t.Fatalf("expected exit code 1, got %d", run.Result.ExitStatus())
	}

	output := b.CLI.(*cli.MockUi).OutputWriter.String()
	if !strings.Contains(output, "Running plan in the remote backend") {
		t.Fatalf("expected remote backend header in output: %s", output)
	}
	if !strings.Contains(output, "null_resource.foo: 1 error") {
		t.Fatalf("expected plan error in output: %s", output)
	}
}

func TestRemote_planOtherError(t *testing.T) {
	b, bCleanup := testBackendDefault(t)
	defer bCleanup()

	op, configCleanup := testOperationPlan(t, "./testdata/plan")
	defer configCleanup()

	op.Workspace = "network-error" // custom error response in backend_mock.go

	_, err := b.Operation(context.Background(), op)
	if err == nil {
		t.Errorf("expected error, got success")
	}

	if !strings.Contains(err.Error(),
		"The configured \"remote\" backend encountered an unexpected error:\n\nI'm a little teacup") {
		t.Fatalf("expected error message, got: %s", err.Error())
	}
}
