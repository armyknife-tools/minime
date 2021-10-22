package command

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/internal/backend"
	"github.com/hashicorp/terraform/internal/cloud"
	"github.com/hashicorp/terraform/internal/command/arguments"
	"github.com/hashicorp/terraform/internal/command/clistate"
	"github.com/hashicorp/terraform/internal/command/views"
	"github.com/hashicorp/terraform/internal/states"
	"github.com/hashicorp/terraform/internal/states/statemgr"
	"github.com/hashicorp/terraform/internal/terraform"
)

type backendMigrateOpts struct {
	SourceType, DestinationType string
	Source, Destination         backend.Backend

	// Fields below are set internally when migrate is called

	sourceWorkspace      string
	destinationWorkspace string
	force                bool // if true, won't ask for confirmation
}

// backendMigrateState handles migrating (copying) state from one backend
// to another. This function handles asking the user for confirmation
// as well as the copy itself.
//
// This function can handle all scenarios of state migration regardless
// of the existence of state in either backend.
//
// After migrating the state, the existing state in the first backend
// remains untouched.
//
// This will attempt to lock both states for the migration.
func (m *Meta) backendMigrateState(opts *backendMigrateOpts) error {
	log.Printf("[INFO] backendMigrateState: need to migrate from %q to %q backend config", opts.SourceType, opts.DestinationType)
	// We need to check what the named state status is. If we're converting
	// from multi-state to single-state for example, we need to handle that.
	var sourceSingleState, destinationSingleState, sourceTFC, destinationTFC bool

	_, sourceTFC = opts.Source.(*cloud.Cloud)
	_, destinationTFC = opts.Destination.(*cloud.Cloud)

	sourceWorkspaces, sourceSingleState, err := retrieveWorkspaces(opts.Source, opts.SourceType)
	if err != nil {
		return err
	}
	destinationWorkspaces, destinationSingleState, err := retrieveWorkspaces(opts.Destination, opts.SourceType)
	if err != nil {
		return err
	}

	// Set up defaults
	opts.sourceWorkspace = backend.DefaultStateName
	opts.destinationWorkspace = backend.DefaultStateName
	opts.force = m.forceInitCopy

	// Disregard remote Terraform version for the state source backend. If it's a
	// Terraform Cloud remote backend, we don't care about the remote version,
	// as we are migrating away and will not break a remote workspace.
	m.ignoreRemoteVersionConflict(opts.Source)

	// Disregard remote Terraform version if instructed to do so via CLI flag.
	if m.ignoreRemoteVersion {
		m.ignoreRemoteVersionConflict(opts.Destination)
	} else {
		// Check the remote Terraform version for the state destination backend. If
		// it's a Terraform Cloud remote backend, we want to ensure that we don't
		// break the workspace by uploading an incompatible state file.
		for _, workspace := range destinationWorkspaces {
			diags := m.remoteVersionCheck(opts.Destination, workspace)
			if diags.HasErrors() {
				return diags.Err()
			}
		}
		// If there are no specified destination workspaces, perform a remote
		// backend version check with the default workspace.
		if len(destinationWorkspaces) == 0 {
			diags := m.remoteVersionCheck(opts.Destination, backend.DefaultStateName)
			if diags.HasErrors() {
				return diags.Err()
			}
		}
	}

	// Determine migration behavior based on whether the source/destination
	// supports multi-state.
	switch {
	case sourceTFC || destinationTFC:
		return m.backendMigrateTFC(opts)

	// Single-state to single-state. This is the easiest case: we just
	// copy the default state directly.
	case sourceSingleState && destinationSingleState:
		return m.backendMigrateState_s_s(opts)

	// Single-state to multi-state. This is easy since we just copy
	// the default state and ignore the rest in the destination.
	case sourceSingleState && !destinationSingleState:
		return m.backendMigrateState_s_s(opts)

	// Multi-state to single-state. If the source has more than the default
	// state this is complicated since we have to ask the user what to do.
	case !sourceSingleState && destinationSingleState:
		// If the source only has one state and it is the default,
		// treat it as if it doesn't support multi-state.
		if len(sourceWorkspaces) == 1 && sourceWorkspaces[0] == backend.DefaultStateName {
			return m.backendMigrateState_s_s(opts)
		}

		return m.backendMigrateState_S_s(opts)

	// Multi-state to multi-state. We merge the states together (migrating
	// each from the source to the destination one by one).
	case !sourceSingleState && !destinationSingleState:
		// If the source only has one state and it is the default,
		// treat it as if it doesn't support multi-state.
		if len(sourceWorkspaces) == 1 && sourceWorkspaces[0] == backend.DefaultStateName {
			return m.backendMigrateState_s_s(opts)
		}

		return m.backendMigrateState_S_S(opts)
	}

	return nil
}

//-------------------------------------------------------------------
// State Migration Scenarios
//
// The functions below cover handling all the various scenarios that
// can exist when migrating state. They are named in an immediately not
// obvious format but is simple:
//
// Format: backendMigrateState_s1_s2[_suffix]
//
// When s1 or s2 is lower case, it means that it is a single state backend.
// When either is uppercase, it means that state is a multi-state backend.
// The suffix is used to disambiguate multiple cases with the same type of
// states.
//
//-------------------------------------------------------------------

// Multi-state to multi-state.
func (m *Meta) backendMigrateState_S_S(opts *backendMigrateOpts) error {
	log.Print("[INFO] backendMigrateState: migrating all named workspaces")

	migrate := opts.force
	if !migrate {
		var err error
		// Ask the user if they want to migrate their existing remote state
		migrate, err = m.confirm(&terraform.InputOpts{
			Id: "backend-migrate-multistate-to-multistate",
			Query: fmt.Sprintf(
				"Do you want to migrate all workspaces to %q?",
				opts.DestinationType),
			Description: fmt.Sprintf(
				strings.TrimSpace(inputBackendMigrateMultiToMulti),
				opts.SourceType, opts.DestinationType),
		})
		if err != nil {
			return fmt.Errorf(
				"Error asking for state migration action: %s", err)
		}
	}
	if !migrate {
		return fmt.Errorf("Migration aborted by user.")
	}

	// Read all the states
	sourceWorkspaces, err := opts.Source.Workspaces()
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(
			errMigrateLoadStates), opts.SourceType, err)
	}

	// Sort the states so they're always copied alphabetically
	sort.Strings(sourceWorkspaces)

	// Go through each and migrate
	for _, name := range sourceWorkspaces {
		// Copy the same names
		opts.sourceWorkspace = name
		opts.destinationWorkspace = name

		// Force it, we confirmed above
		opts.force = true

		// Perform the migration
		if err := m.backendMigrateState_s_s(opts); err != nil {
			return fmt.Errorf(strings.TrimSpace(
				errMigrateMulti), name, opts.SourceType, opts.DestinationType, err)
		}
	}

	return nil
}

// Multi-state to single state.
func (m *Meta) backendMigrateState_S_s(opts *backendMigrateOpts) error {
	log.Printf("[INFO] backendMigrateState: destination backend type %q does not support named workspaces", opts.DestinationType)

	currentWorkspace, err := m.Workspace()
	if err != nil {
		return err
	}

	migrate := opts.force
	if !migrate {
		var err error
		// Ask the user if they want to migrate their existing remote state
		migrate, err = m.confirm(&terraform.InputOpts{
			Id: "backend-migrate-multistate-to-single",
			Query: fmt.Sprintf(
				"Destination state %q doesn't support workspaces.\n"+
					"Do you want to copy only your current workspace?",
				opts.DestinationType),
			Description: fmt.Sprintf(
				strings.TrimSpace(inputBackendMigrateMultiToSingle),
				opts.SourceType, opts.DestinationType, currentWorkspace),
		})
		if err != nil {
			return fmt.Errorf(
				"Error asking for state migration action: %s", err)
		}
	}

	if !migrate {
		return fmt.Errorf("Migration aborted by user.")
	}

	// Copy the default state
	opts.sourceWorkspace = currentWorkspace

	// now switch back to the default env so we can acccess the new backend
	m.SetWorkspace(backend.DefaultStateName)

	return m.backendMigrateState_s_s(opts)
}

// Single state to single state, assumed default state name.
func (m *Meta) backendMigrateState_s_s(opts *backendMigrateOpts) error {
	log.Printf("[INFO] backendMigrateState: single-to-single migrating %q workspace to %q workspace", opts.sourceWorkspace, opts.destinationWorkspace)

	sourceState, err := opts.Source.StateMgr(opts.sourceWorkspace)
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(
			errMigrateSingleLoadDefault), opts.SourceType, err)
	}
	if err := sourceState.RefreshState(); err != nil {
		return fmt.Errorf(strings.TrimSpace(
			errMigrateSingleLoadDefault), opts.SourceType, err)
	}

	// Do not migrate workspaces without state.
	if sourceState.State().Empty() {
		log.Print("[TRACE] backendMigrateState: source workspace has empty state, so nothing to migrate")
		return nil
	}

	destinationState, err := opts.Destination.StateMgr(opts.destinationWorkspace)
	if err == backend.ErrDefaultWorkspaceNotSupported {
		// If the backend doesn't support using the default state, we ask the user
		// for a new name and migrate the default state to the given named state.
		destinationState, err = func() (statemgr.Full, error) {
			log.Print("[TRACE] backendMigrateState: destination doesn't support a default workspace, so we must prompt for a new name")
			name, err := m.promptNewWorkspaceName(opts.DestinationType)
			if err != nil {
				return nil, err
			}

			// Update the name of the destination state.
			opts.destinationWorkspace = name

			destinationState, err := opts.Destination.StateMgr(opts.destinationWorkspace)
			if err != nil {
				return nil, err
			}

			// Ignore invalid workspace name as it is irrelevant in this context.
			workspace, _ := m.Workspace()

			// If the currently selected workspace is the default workspace, then set
			// the named workspace as the new selected workspace.
			if workspace == backend.DefaultStateName {
				if err := m.SetWorkspace(opts.destinationWorkspace); err != nil {
					return nil, fmt.Errorf("Failed to set new workspace: %s", err)
				}
			}

			return destinationState, nil
		}()
	}
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(
			errMigrateSingleLoadDefault), opts.DestinationType, err)
	}
	if err := destinationState.RefreshState(); err != nil {
		return fmt.Errorf(strings.TrimSpace(
			errMigrateSingleLoadDefault), opts.DestinationType, err)
	}

	// Check if we need migration at all.
	// This is before taking a lock, because they may also correspond to the same lock.
	source := sourceState.State()
	destination := destinationState.State()

	// no reason to migrate if the state is already there
	if source.Equal(destination) {
		// Equal isn't identical; it doesn't check lineage.
		sm1, _ := sourceState.(statemgr.PersistentMeta)
		sm2, _ := destinationState.(statemgr.PersistentMeta)
		if source != nil && destination != nil {
			if sm1 == nil || sm2 == nil {
				log.Print("[TRACE] backendMigrateState: both source and destination workspaces have no state, so no migration is needed")
				return nil
			}
			if sm1.StateSnapshotMeta().Lineage == sm2.StateSnapshotMeta().Lineage {
				log.Printf("[TRACE] backendMigrateState: both source and destination workspaces have equal state with lineage %q, so no migration is needed", sm1.StateSnapshotMeta().Lineage)
				return nil
			}
		}
	}

	if m.stateLock {
		lockCtx := context.Background()

		view := views.NewStateLocker(arguments.ViewHuman, m.View)
		locker := clistate.NewLocker(m.stateLockTimeout, view)

		lockerSource := locker.WithContext(lockCtx)
		if diags := lockerSource.Lock(sourceState, "migration source state"); diags.HasErrors() {
			return diags.Err()
		}
		defer lockerSource.Unlock()

		lockerDestination := locker.WithContext(lockCtx)
		if diags := lockerDestination.Lock(destinationState, "migration destination state"); diags.HasErrors() {
			return diags.Err()
		}
		defer lockerDestination.Unlock()

		// We now own a lock, so double check that we have the version
		// corresponding to the lock.
		log.Print("[TRACE] backendMigrateState: refreshing source workspace state")
		if err := sourceState.RefreshState(); err != nil {
			return fmt.Errorf(strings.TrimSpace(
				errMigrateSingleLoadDefault), opts.SourceType, err)
		}
		log.Print("[TRACE] backendMigrateState: refreshing destination workspace state")
		if err := destinationState.RefreshState(); err != nil {
			return fmt.Errorf(strings.TrimSpace(
				errMigrateSingleLoadDefault), opts.SourceType, err)
		}

		source = sourceState.State()
		destination = destinationState.State()
	}

	var confirmFunc func(statemgr.Full, statemgr.Full, *backendMigrateOpts) (bool, error)
	switch {
	// No migration necessary
	case source.Empty() && destination.Empty():
		log.Print("[TRACE] backendMigrateState: both source and destination workspaces have empty state, so no migration is required")
		return nil

	// No migration necessary if we're inheriting state.
	case source.Empty() && !destination.Empty():
		log.Print("[TRACE] backendMigrateState: source workspace has empty state, so no migration is required")
		return nil

	// We have existing state moving into no state. Ask the user if
	// they'd like to do this.
	case !source.Empty() && destination.Empty():
		log.Print("[TRACE] backendMigrateState: destination workspace has empty state, so might copy source workspace state")
		confirmFunc = m.backendMigrateEmptyConfirm

	// Both states are non-empty, meaning we need to determine which
	// state should be used and update accordingly.
	case !source.Empty() && !destination.Empty():
		log.Print("[TRACE] backendMigrateState: both source and destination workspaces have states, so might overwrite destination with source")
		confirmFunc = m.backendMigrateNonEmptyConfirm
	}

	if confirmFunc == nil {
		panic("confirmFunc must not be nil")
	}

	if !opts.force {
		// Abort if we can't ask for input.
		if !m.input {
			log.Print("[TRACE] backendMigrateState: can't prompt for input, so aborting migration")
			return errors.New("error asking for state migration action: input disabled")
		}

		// Confirm with the user whether we want to copy state over
		confirm, err := confirmFunc(sourceState, destinationState, opts)
		if err != nil {
			log.Print("[TRACE] backendMigrateState: error reading input, so aborting migration")
			return err
		}
		if !confirm {
			log.Print("[TRACE] backendMigrateState: user cancelled at confirmation prompt, so aborting migration")
			return nil
		}
	}

	// Confirmed! We'll have the statemgr package handle the migration, which
	// includes preserving any lineage/serial information where possible, if
	// both managers support such metadata.
	log.Print("[TRACE] backendMigrateState: migration confirmed, so migrating")
	if err := statemgr.Migrate(destinationState, sourceState); err != nil {
		return fmt.Errorf(strings.TrimSpace(errBackendStateCopy),
			opts.SourceType, opts.DestinationType, err)
	}
	if err := destinationState.PersistState(); err != nil {
		return fmt.Errorf(strings.TrimSpace(errBackendStateCopy),
			opts.SourceType, opts.DestinationType, err)
	}

	// And we're done.
	return nil
}

func (m *Meta) backendMigrateEmptyConfirm(source, destination statemgr.Full, opts *backendMigrateOpts) (bool, error) {
	var inputOpts *terraform.InputOpts
	if opts.DestinationType == "cloud" {
		inputOpts = &terraform.InputOpts{
			Id:          "backend-migrate-copy-to-empty-cloud",
			Query:       "Do you want to copy existing state to Terraform Cloud?",
			Description: fmt.Sprintf(strings.TrimSpace(inputBackendMigrateEmptyCloud), opts.SourceType),
		}
	} else {
		inputOpts = &terraform.InputOpts{
			Id:    "backend-migrate-copy-to-empty",
			Query: "Do you want to copy existing state to the new backend?",
			Description: fmt.Sprintf(
				strings.TrimSpace(inputBackendMigrateEmpty),
				opts.SourceType, opts.DestinationType),
		}
	}

	return m.confirm(inputOpts)
}

func (m *Meta) backendMigrateNonEmptyConfirm(
	sourceState, destinationState statemgr.Full, opts *backendMigrateOpts) (bool, error) {
	// We need to grab both states so we can write them to a file
	source := sourceState.State()
	destination := destinationState.State()

	// Save both to a temporary
	td, err := ioutil.TempDir("", "terraform")
	if err != nil {
		return false, fmt.Errorf("Error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(td)

	// Helper to write the state
	saveHelper := func(n, path string, s *states.State) error {
		mgr := statemgr.NewFilesystem(path)
		return mgr.WriteState(s)
	}

	// Write the states
	sourcePath := filepath.Join(td, fmt.Sprintf("1-%s.tfstate", opts.SourceType))
	destinationPath := filepath.Join(td, fmt.Sprintf("2-%s.tfstate", opts.DestinationType))
	if err := saveHelper(opts.SourceType, sourcePath, source); err != nil {
		return false, fmt.Errorf("Error saving temporary state: %s", err)
	}
	if err := saveHelper(opts.DestinationType, destinationPath, destination); err != nil {
		return false, fmt.Errorf("Error saving temporary state: %s", err)
	}

	// Ask for confirmation
	var inputOpts *terraform.InputOpts
	if opts.DestinationType == "cloud" {
		inputOpts = &terraform.InputOpts{
			Id:    "backend-migrate-to-tfc",
			Query: "Do you want to copy existing state to Terraform Cloud?",
			Description: fmt.Sprintf(
				strings.TrimSpace(inputBackendMigrateNonEmptyCloud),
				opts.SourceType, sourcePath, destinationPath),
		}
	} else {
		inputOpts = &terraform.InputOpts{
			Id:    "backend-migrate-to-backend",
			Query: "Do you want to copy existing state to the new backend?",
			Description: fmt.Sprintf(
				strings.TrimSpace(inputBackendMigrateNonEmpty),
				opts.SourceType, opts.DestinationType, sourcePath, destinationPath),
		}
	}

	// Confirm with the user that the copy should occur
	return m.confirm(inputOpts)
}

func retrieveWorkspaces(back backend.Backend, sourceType string) ([]string, bool, error) {
	var singleState bool
	var err error
	workspaces, err := back.Workspaces()
	if err == backend.ErrWorkspacesNotSupported {
		singleState = true
		err = nil
	}
	if err != nil {
		return nil, singleState, fmt.Errorf(strings.TrimSpace(
			errMigrateLoadStates), sourceType, err)
	}

	return workspaces, singleState, err
}

func (m *Meta) backendMigrateTFC(opts *backendMigrateOpts) error {
	_, sourceTFC := opts.Source.(*cloud.Cloud)
	cloudBackendDestination, destinationTFC := opts.Destination.(*cloud.Cloud)

	sourceWorkspaces, sourceSingleState, err := retrieveWorkspaces(opts.Source, opts.SourceType)
	if err != nil {
		return err
	}
	//to be used below, not yet implamented
	// destinationWorkspaces, destinationSingleState
	_, _, err = retrieveWorkspaces(opts.Destination, opts.SourceType)
	if err != nil {
		return err
	}

	// from TFC to non-TFC backend
	if sourceTFC && !destinationTFC {
		// From Terraform Cloud to another backend. This is not yet implemented, and
		// we recommend people to use the TFC API.
		return fmt.Errorf(strings.TrimSpace(errTFCMigrateNotYetImplemented))
	}

	// Everything below, by the above two conditionals, now assumes that the
	// destination is always Terraform Cloud (TFC).

	sourceSingle := sourceSingleState || (len(sourceWorkspaces) == 1)
	if sourceSingle {
		if cloudBackendDestination.WorkspaceMapping.Strategy() == cloud.WorkspaceNameStrategy {
			// If we know the name via WorkspaceNameStrategy, then set the
			// destinationWorkspace to the new Name and skip the user prompt. Here the
			// destinationWorkspace is not set to `default` thereby we will create it
			// in TFC if it does not exist.
			opts.destinationWorkspace = cloudBackendDestination.WorkspaceMapping.Name
		}

		currentWorkspace, err := m.Workspace()
		if err != nil {
			return err
		}
		opts.sourceWorkspace = currentWorkspace

		log.Printf("[INFO] backendMigrateTFC: single-to-single migration from source %s to destination %q", opts.sourceWorkspace, opts.destinationWorkspace)
		// Run normal single-to-single state migration
		// This will handle both situations where the new cloud backend
		// configuration is using a workspace.name strategy or workspace.tags
		// strategy.
		return m.backendMigrateState_s_s(opts)
	}

	destinationTagsStrategy := cloudBackendDestination.WorkspaceMapping.Strategy() == cloud.WorkspaceTagsStrategy
	destinationNameStrategy := cloudBackendDestination.WorkspaceMapping.Strategy() == cloud.WorkspaceNameStrategy

	multiSource := !sourceSingleState && len(sourceWorkspaces) > 1
	if multiSource && destinationNameStrategy {
		currentWorkspace, err := m.Workspace()
		if err != nil {
			return err
		}

		opts.sourceWorkspace = currentWorkspace
		opts.destinationWorkspace = cloudBackendDestination.WorkspaceMapping.Name
		if err := m.promptMultiToSingleCloudMigration(opts); err != nil {
			return err
		}

		log.Printf("[INFO] backendMigrateTFC: multi-to-single migration from source %s to destination %q", opts.sourceWorkspace, opts.destinationWorkspace)

		return m.backendMigrateState_s_s(opts)
	}

	// Multiple sources, and using tags strategy. So migrate every source
	// workspace over to new one, prompt for workspace name pattern (*),
	// and start migrating, and create tags for each workspace.
	if multiSource && destinationTagsStrategy {
		log.Printf("[INFO] backendMigrateTFC: multi-to-multi migration from source workspaces %q", sourceWorkspaces)
		return m.backendMigrateState_S_TFC(opts, sourceWorkspaces)
	}

	// TODO(omar): after the check for sourceSingle is done, everything following
	// it has to be multi. So rework the code to not need to check for multi, adn
	// return m.backendMigrateState_S_TFC here.
	return nil
}

// migrates a multi-state backend to Terraform Cloud
func (m *Meta) backendMigrateState_S_TFC(opts *backendMigrateOpts, sourceWorkspaces []string) error {
	log.Print("[TRACE] backendMigrateState: migrating all named workspaces")

	// This map is used later when doing the migration per source/destination.
	// If a source has 'default', then we ask what the new name should be.
	// And further down when we actually run state migration for each
	// sourc/destination workspce, we use this new name (where source is 'default')
	// and set as destinationWorkspace.
	defaultNewName := map[string]string{}
	for i := 0; i < len(sourceWorkspaces); i++ {
		if sourceWorkspaces[i] == backend.DefaultStateName {
			newName, err := m.promptNewWorkspaceName(opts.DestinationType)
			if err != nil {
				return err
			}
			defaultNewName[sourceWorkspaces[i]] = newName
		}
	}
	pattern, err := m.promptMultiStateMigrationPattern(opts.SourceType)
	if err != nil {
		return err
	}

	// Go through each and migrate
	for _, name := range sourceWorkspaces {

		// Copy the same names
		opts.sourceWorkspace = name
		if newName, ok := defaultNewName[name]; ok {
			// this has to be done before setting destinationWorkspace
			name = newName
		}
		opts.destinationWorkspace = strings.Replace(pattern, "*", name, -1)

		// Force it, we confirmed above
		opts.force = true

		// Perform the migration
		log.Printf("[INFO] backendMigrateTFC: multi-to-multi migration, source workspace %q to destination workspace %q", opts.sourceWorkspace, opts.destinationWorkspace)
		if err := m.backendMigrateState_s_s(opts); err != nil {
			return fmt.Errorf(strings.TrimSpace(
				errMigrateMulti), name, opts.SourceType, opts.DestinationType, err)
		}
	}

	// After migrating multiple workspaces, we want to ensure that a workspace is
	// set or we prompt the user to set a workspace.
	err = m.selectWorkspace(opts.Destination)
	if err != nil {
		return err
	}

	return nil
}

// Multi-state to single state.
func (m *Meta) promptMultiToSingleCloudMigration(opts *backendMigrateOpts) error {
	migrate := opts.force
	if !migrate {
		var err error
		// Ask the user if they want to migrate their existing remote state
		migrate, err = m.confirm(&terraform.InputOpts{
			Id:    "backend-migrate-multistate-to-single",
			Query: "Do you want to copy only your current workspace?",
			Description: fmt.Sprintf(
				strings.TrimSpace(tfcInputBackendMigrateMultiToSingle),
				opts.SourceType, opts.destinationWorkspace),
		})
		if err != nil {
			return fmt.Errorf("Error asking for state migration action: %s", err)
		}
	}

	if !migrate {
		return fmt.Errorf("Migration aborted by user.")
	}

	return nil
}

func (m *Meta) promptNewWorkspaceName(destinationType string) (string, error) {
	message := fmt.Sprintf("[reset][bold][yellow]The %q backend configuration only allows "+
		"named workspaces![reset]", destinationType)
	if destinationType == "cloud" {
		message = fmt.Sprintf("[reset][bold][yellow]The Terraform Cloud configuration only allows " +
			"named workspaces![reset]")
	}
	name, err := m.UIInput().Input(context.Background(), &terraform.InputOpts{
		Id:          "new-state-name",
		Query:       message,
		Description: strings.TrimSpace(inputBackendNewWorkspaceName),
	})
	if err != nil {
		return "", fmt.Errorf("Error asking for new state name: %s", err)
	}

	return name, nil
}

func (m *Meta) promptMultiStateMigrationPattern(sourceType string) (string, error) {
	renameWorkspaces, err := m.UIInput().Input(context.Background(), &terraform.InputOpts{
		Id:          "backend-migrate-multistate-to-tfc",
		Query:       fmt.Sprintf("[reset][bold][yellow]%s[reset]", "Would you like to rename your workspaces?"),
		Description: fmt.Sprintf(strings.TrimSpace(tfcInputBackendMigrateMultiToMulti), sourceType),
	})
	if err != nil {
		return "", fmt.Errorf("Error asking for state migration action: %s", err)
	}
	if renameWorkspaces != "2" && renameWorkspaces != "1" {
		return "", fmt.Errorf("Please select 1 or 2 as part of this option.")
	}
	if renameWorkspaces == "2" {
		// this means they did not want to rename their workspaces, and we are
		// returning a generic '*' that means use the same workspace name during
		// migration.
		return "*", nil
	}

	pattern, err := m.UIInput().Input(context.Background(), &terraform.InputOpts{
		Id:          "backend-migrate-multistate-to-tfc-pattern",
		Query:       fmt.Sprintf("[reset][bold][yellow]%s[reset]", "What pattern would you like to add to all your workspaces?"),
		Description: strings.TrimSpace(tfcInputBackendMigrateMultiToMultiPattern),
	})
	if err != nil {
		return "", fmt.Errorf("Error asking for state migration action: %s", err)
	}
	if !strings.Contains(pattern, "*") {
		return "", fmt.Errorf("The pattern must have an '*'")
	}

	if count := strings.Count(pattern, "*"); count > 1 {
		return "", fmt.Errorf("The pattern '*' cannot be used more than once.")
	}

	return pattern, nil
}

const errMigrateLoadStates = `
Error inspecting states in the %q backend:
    %s

Prior to changing backends, Terraform inspects the source and destination
states to determine what kind of migration steps need to be taken, if any.
Terraform failed to load the states. The data in both the source and the
destination remain unmodified. Please resolve the above error and try again.
`

const errMigrateSingleLoadDefault = `
Error loading state:
    %[2]s

Terraform failed to load the default state from the %[1]q backend.
State migration cannot occur unless the state can be loaded. Backend
modification and state migration has been aborted. The state in both the
source and the destination remain unmodified. Please resolve the
above error and try again.
`

const errMigrateMulti = `
Error migrating the workspace %q from the previous %q backend
to the newly configured %q backend:
    %s

Terraform copies workspaces in alphabetical order. Any workspaces
alphabetically earlier than this one have been copied. Any workspaces
later than this haven't been modified in the destination. No workspaces
in the source state have been modified.

Please resolve the error above and run the initialization command again.
This will attempt to copy (with permission) all workspaces again.
`

const errBackendStateCopy = `
Error copying state from the previous %q backend to the newly configured 
%q backend:
    %s

The state in the previous backend remains intact and unmodified. Please resolve
the error above and try again.
`

const errTFCMigrateNotYetImplemented = `
Migrating state from Terraform Cloud to another backend is not yet implemented.

Please use the API to do this: https://www.terraform.io/docs/cloud/api/state-versions.html
`

const tfcInputBackendMigrateMultiToMultiPattern = `
If you choose to NOT rename your workspaces, just input "*".

The asterisk "*" represents your workspace name. Here are a few examples
if a workspace was named 'prod':
* input: 'app-*'; output: 'app-prod'
* input: '*-app', output: 'prod-app'
* input: 'app-*-service', output: 'app-prod-service'
* input: '*'; output: 'prod'
`

const tfcInputBackendMigrateMultiToMulti = `
When migrating existing workspaces from the backend %[1]q to Terraform Cloud, would you like to
rename your workspaces?

Unlike typical Terraform workspaces representing an environment associated with a particular
configuration (e.g. production, staging, development), Terraform Cloud workspaces are named uniquely
across all configurations used within an organization. A typical strategy to start with is
<COMPONENT>-<ENVIRONMENT>-<REGION> (e.g. networking-prod-us-east, networking-staging-us-east).

For more information on workspace naming, see https://www.terraform.io/docs/cloud/workspaces/naming.html

1. Yes, rename workspaces according to a pattern.
2. No, I would not like to rename my workspaces. Migrate them as currently named.
`

const tfcInputBackendMigrateMultiToSingle = `
The previous backend %[1]q has multiple workspaces, but Terraform Cloud has been
configured to use a single workspace (%[2]q). By continuing, you will only
migrate your current workspace. If you wish to migrate all workspaces from the
previous backend, use the 'tags' strategy in your workspace configuration block
instead.
`

const inputBackendMigrateEmpty = `
Pre-existing state was found while migrating the previous %q backend to the
newly configured %q backend. No existing state was found in the newly
configured %[2]q backend. Do you want to copy this state to the new %[2]q
backend? Enter "yes" to copy and "no" to start with an empty state.
`

const inputBackendMigrateEmptyCloud = `
Pre-existing state was found while migrating the previous %q backend to Terraform Cloud.
No existing state was found in Terraform Cloud. Do you want to copy this state to Terraform Cloud?
Enter "yes" to copy and "no" to start with an empty state.
`

const inputBackendMigrateNonEmpty = `
Pre-existing state was found while migrating the previous %q backend to the
newly configured %q backend. An existing non-empty state already exists in
the new backend. The two states have been saved to temporary files that will be
removed after responding to this query.

Previous (type %[1]q): %[3]s
New      (type %[2]q): %[4]s

Do you want to overwrite the state in the new backend with the previous state?
Enter "yes" to copy and "no" to start with the existing state in the newly
configured %[2]q backend.
`

const inputBackendMigrateNonEmptyCloud = `
Pre-existing state was found while migrating the previous %q backend to
Terraform Cloud. An existing non-empty state already exists in Terraform Cloud.
The two states have been saved to temporary files that will be removed after
responding to this query.

Previous (type %[1]q): %[2]s
New      (Terraform Cloud): %[3]s

Do you want to overwrite the state in Terraform Cloud with the previous state?
Enter "yes" to copy and "no" to start with the existing state in Terraform Cloud.
`

const inputBackendMigrateMultiToSingle = `
The existing %[1]q backend supports workspaces and you currently are
using more than one. The newly configured %[2]q backend doesn't support
workspaces. If you continue, Terraform will copy your current workspace %[3]q
to the default workspace in the new backend. Your existing workspaces in the
source backend won't be modified. If you want to switch workspaces, back them
up, or cancel altogether, answer "no" and Terraform will abort.
`

const inputBackendMigrateMultiToMulti = `
Both the existing %[1]q backend and the newly configured %[2]q backend
support workspaces. When migrating between backends, Terraform will copy
all workspaces (with the same names). THIS WILL OVERWRITE any conflicting
states in the destination.

Terraform initialization doesn't currently migrate only select workspaces.
If you want to migrate a select number of workspaces, you must manually
pull and push those states.

If you answer "yes", Terraform will migrate all states. If you answer
"no", Terraform will abort.
`

const inputBackendNewWorkspaceName = `
Please provide a new workspace name (e.g. dev, test) that will be used
to migrate the existing default workspace. 
`

const inputBackendSelectWorkspace = `
This is expected behavior when the selected workspace did not have an
existing non-empty state. Please enter a number to select a workspace:

%s
`
