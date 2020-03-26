package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/posener/complete"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	backendInit "github.com/hashicorp/terraform/backend/init"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/internal/earlyconfig"
	"github.com/hashicorp/terraform/internal/getproviders"
	"github.com/hashicorp/terraform/internal/initwd"
	"github.com/hashicorp/terraform/internal/providercache"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/tfdiags"
)

// InitCommand is a Command implementation that takes a Terraform
// module and clones it to the working directory.
type InitCommand struct {
	Meta

	// getPlugins is for the -get-plugins flag
	getPlugins bool

	// providerInstaller is used to download and install providers that
	// aren't found locally. This uses a discovery.ProviderInstaller instance
	// by default, but it can be overridden here as a way to mock fetching
	// providers for tests.
	providerInstaller discovery.Installer
}

func (c *InitCommand) Run(args []string) int {
	var flagFromModule string
	var flagBackend, flagGet, flagUpgrade bool
	var flagPluginPath FlagStringSlice
	var flagVerifyPlugins bool
	flagConfigExtra := newRawFlags("-backend-config")

	args = c.Meta.process(args)
	cmdFlags := c.Meta.extendedFlagSet("init")
	cmdFlags.BoolVar(&flagBackend, "backend", true, "")
	cmdFlags.Var(flagConfigExtra, "backend-config", "")
	cmdFlags.StringVar(&flagFromModule, "from-module", "", "copy the source of the given module into the directory before init")
	cmdFlags.BoolVar(&flagGet, "get", true, "")
	cmdFlags.BoolVar(&c.getPlugins, "get-plugins", true, "")
	cmdFlags.BoolVar(&c.forceInitCopy, "force-copy", false, "suppress prompts about copying state data")
	cmdFlags.BoolVar(&c.Meta.stateLock, "lock", true, "lock state")
	cmdFlags.DurationVar(&c.Meta.stateLockTimeout, "lock-timeout", 0, "lock timeout")
	cmdFlags.BoolVar(&c.reconfigure, "reconfigure", false, "reconfigure")
	cmdFlags.BoolVar(&flagUpgrade, "upgrade", false, "")
	cmdFlags.Var(&flagPluginPath, "plugin-dir", "plugin directory")
	cmdFlags.BoolVar(&flagVerifyPlugins, "verify-plugins", true, "verify plugins")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	var diags tfdiags.Diagnostics

	if len(flagPluginPath) > 0 {
		c.pluginPath = flagPluginPath
		c.getPlugins = false
	}

	// set providerInstaller if we don't have a test version already
	if c.providerInstaller == nil {
		c.providerInstaller = &discovery.ProviderInstaller{
			Dir:                   c.pluginDir(),
			Cache:                 c.pluginCache(),
			PluginProtocolVersion: discovery.PluginInstallProtocolVersion,
			SkipVerify:            !flagVerifyPlugins,
			Ui:                    c.Ui,
			Services:              c.Services,
		}
	}

	// Validate the arg count
	args = cmdFlags.Args()
	if len(args) > 1 {
		c.Ui.Error("The init command expects at most one argument.\n")
		cmdFlags.Usage()
		return 1
	}

	if err := c.storePluginPath(c.pluginPath); err != nil {
		c.Ui.Error(fmt.Sprintf("Error saving -plugin-path values: %s", err))
		return 1
	}

	// Get our pwd. We don't always need it but always getting it is easier
	// than the logic to determine if it is or isn't needed.
	pwd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting pwd: %s", err))
		return 1
	}

	// If an argument is provided then it overrides our working directory.
	path := pwd
	if len(args) == 1 {
		path = args[0]
	}

	// This will track whether we outputted anything so that we know whether
	// to output a newline before the success message
	var header bool

	if flagFromModule != "" {
		src := flagFromModule

		empty, err := configs.IsEmptyDir(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error validating destination directory: %s", err))
			return 1
		}
		if !empty {
			c.Ui.Error(strings.TrimSpace(errInitCopyNotEmpty))
			return 1
		}

		c.Ui.Output(c.Colorize().Color(fmt.Sprintf(
			"[reset][bold]Copying configuration[reset] from %q...", src,
		)))
		header = true

		hooks := uiModuleInstallHooks{
			Ui:             c.Ui,
			ShowLocalPaths: false, // since they are in a weird location for init
		}

		initDiags := c.initDirFromModule(path, src, hooks)
		diags = diags.Append(initDiags)
		if initDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}

		c.Ui.Output("")
	}

	// If our directory is empty, then we're done. We can't get or setup
	// the backend with an empty directory.
	empty, err := configs.IsEmptyDir(path)
	if err != nil {
		diags = diags.Append(fmt.Errorf("Error checking configuration: %s", err))
		return 1
	}
	if empty {
		c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitEmpty)))
		return 0
	}

	// Before we do anything else, we'll try loading configuration with both
	// our "normal" and "early" configuration codepaths. If early succeeds
	// while normal fails, that strongly suggests that the configuration is
	// using syntax that worked in 0.11 but no longer in v0.12.
	rootMod, confDiags := c.loadSingleModule(path)
	rootModEarly, earlyConfDiags := c.loadSingleModuleEarly(path)
	if confDiags.HasErrors() {
		if earlyConfDiags.HasErrors() {
			// If both parsers produced errors then we'll assume the config
			// is _truly_ invalid and produce error messages as normal.
			// Since this may be the user's first ever interaction with Terraform,
			// we'll provide some additional context in this case.
			c.Ui.Error(strings.TrimSpace(errInitConfigError))
			diags = diags.Append(confDiags)
			c.showDiagnostics(diags)
			return 1
		}
		// If _only_ the main loader produced errors then that suggests the
		// configuration is written in 0.11-style syntax. We will return an
		// error suggesting the user upgrade their config manually or with
		// Terraform v0.12
		c.Ui.Error(strings.TrimSpace(errInitConfigErrorMaybeLegacySyntax))
		c.showDiagnostics(earlyConfDiags)
		return 1
	}

	// If _only_ the early loader encountered errors then that's unusual
	// (it should generally be a superset of the normal loader) but we'll
	// return those errors anyway since otherwise we'll probably get
	// some weird behavior downstream. Errors from the early loader are
	// generally not as high-quality since it has less context to work with.
	if earlyConfDiags.HasErrors() {
		c.Ui.Error(strings.TrimSpace(errInitConfigError))
		// Errors from the early loader are generally not as high-quality since
		// it has less context to work with.
		diags = diags.Append(confDiags)
		c.showDiagnostics(diags)
		return 1
	}

	if flagGet {
		modsOutput, modsDiags := c.getModules(path, rootModEarly, flagUpgrade)
		diags = diags.Append(modsDiags)
		if modsDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		if modsOutput {
			header = true
		}
	}

	// With all of the modules (hopefully) installed, we can now try to load the
	// whole configuration tree.
	//
	// Just as above, we'll try loading both with the early and normal config
	// loaders here. Subsequent work will only use the early config, but loading
	// both gives us an opportunity to prefer the better error messages from the
	// normal loader if both fail.

	earlyConfig, earlyConfDiags := c.loadConfigEarly(path)
	if earlyConfDiags.HasErrors() {
		c.Ui.Error(strings.TrimSpace(errInitConfigError))
		diags = diags.Append(earlyConfDiags)
		c.showDiagnostics(diags)
		return 1
	}

	_, confDiags = c.loadConfig(path)
	if confDiags.HasErrors() {
		c.Ui.Error(strings.TrimSpace(errInitConfigError))
		diags = diags.Append(confDiags)
		c.showDiagnostics(diags)
		return 1
	}

	// Before we go further, we'll check to make sure none of the modules in the
	// configuration declare that they don't support this Terraform version, so
	// we can produce a version-related error message rather than
	// potentially-confusing downstream errors.
	versionDiags := initwd.CheckCoreVersionRequirements(earlyConfig)
	diags = diags.Append(versionDiags)
	if versionDiags.HasErrors() {
		c.showDiagnostics(diags)
		return 1
	}

	var back backend.Backend
	if flagBackend {

		be, backendOutput, backendDiags := c.initBackend(rootMod, flagConfigExtra)
		diags = diags.Append(backendDiags)
		if backendDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		if backendOutput {
			header = true
		}
		back = be
	} else {
		// load the previously-stored backend config
		be, backendDiags := c.Meta.backendFromState()
		diags = diags.Append(backendDiags)
		if backendDiags.HasErrors() {
			c.showDiagnostics(diags)
			return 1
		}
		back = be
	}

	if back == nil {
		// If we didn't initialize a backend then we'll try to at least
		// instantiate one. This might fail if it wasn't already initialized
		// by a previous run, so we must still expect that "back" may be nil
		// in code that follows.
		var backDiags tfdiags.Diagnostics
		back, backDiags = c.Backend(nil)
		if backDiags.HasErrors() {
			// This is fine. We'll proceed with no backend, then.
			back = nil
		}
	}

	var state *states.State

	// If we have a functional backend (either just initialized or initialized
	// on a previous run) we'll use the current state as a potential source
	// of provider dependencies.
	if back != nil {
		sMgr, err := back.StateMgr(c.Workspace())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error loading state: %s", err))
			return 1
		}

		if err := sMgr.RefreshState(); err != nil {
			c.Ui.Error(fmt.Sprintf("Error refreshing state: %s", err))
			return 1
		}

		state = sMgr.State()
	}

	if v := os.Getenv(ProviderSkipVerifyEnvVar); v != "" {
		c.ignorePluginChecksum = true
	}

	// Now that we have loaded all modules, check the module tree for missing providers.
	providersOutput, providerDiags := c.getProviders(earlyConfig, state, flagUpgrade)
	diags = diags.Append(providerDiags)
	if providerDiags.HasErrors() {
		c.showDiagnostics(diags)
		return 1
	}
	if providersOutput {
		header = true
	}

	// If we outputted information, then we need to output a newline
	// so that our success message is nicely spaced out from prior text.
	if header {
		c.Ui.Output("")
	}

	// If we accumulated any warnings along the way that weren't accompanied
	// by errors then we'll output them here so that the success message is
	// still the final thing shown.
	c.showDiagnostics(diags)
	c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitSuccess)))
	if !c.RunningInAutomation {
		// If we're not running in an automation wrapper, give the user
		// some more detailed next steps that are appropriate for interactive
		// shell usage.
		c.Ui.Output(c.Colorize().Color(strings.TrimSpace(outputInitSuccessCLI)))
	}
	return 0
}

func (c *InitCommand) getModules(path string, earlyRoot *tfconfig.Module, upgrade bool) (output bool, diags tfdiags.Diagnostics) {
	if len(earlyRoot.ModuleCalls) == 0 {
		// Nothing to do
		return false, nil
	}

	if upgrade {
		c.Ui.Output(c.Colorize().Color(fmt.Sprintf("[reset][bold]Upgrading modules...")))
	} else {
		c.Ui.Output(c.Colorize().Color(fmt.Sprintf("[reset][bold]Initializing modules...")))
	}

	hooks := uiModuleInstallHooks{
		Ui:             c.Ui,
		ShowLocalPaths: true,
	}
	instDiags := c.installModules(path, upgrade, hooks)
	diags = diags.Append(instDiags)

	// Since module installer has modified the module manifest on disk, we need
	// to refresh the cache of it in the loader.
	if c.configLoader != nil {
		if err := c.configLoader.RefreshModules(); err != nil {
			// Should never happen
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to read module manifest",
				fmt.Sprintf("After installing modules, Terraform could not re-read the manifest of installed modules. This is a bug in Terraform. %s.", err),
			))
		}
	}

	return true, diags
}

func (c *InitCommand) initBackend(root *configs.Module, extraConfig rawFlags) (be backend.Backend, output bool, diags tfdiags.Diagnostics) {
	c.Ui.Output(c.Colorize().Color(fmt.Sprintf("\n[reset][bold]Initializing the backend...")))

	var backendConfig *configs.Backend
	var backendConfigOverride hcl.Body
	if root.Backend != nil {
		backendType := root.Backend.Type
		bf := backendInit.Backend(backendType)
		if bf == nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported backend type",
				Detail:   fmt.Sprintf("There is no backend type named %q.", backendType),
				Subject:  &root.Backend.TypeRange,
			})
			return nil, true, diags
		}

		b := bf()
		backendSchema := b.ConfigSchema()
		backendConfig = root.Backend

		var overrideDiags tfdiags.Diagnostics
		backendConfigOverride, overrideDiags = c.backendConfigOverrideBody(extraConfig, backendSchema)
		diags = diags.Append(overrideDiags)
		if overrideDiags.HasErrors() {
			return nil, true, diags
		}
	} else {
		// If the user supplied a -backend-config on the CLI but no backend
		// block was found in the configuration, it's likely - but not
		// necessarily - a mistake. Return a warning.
		if !extraConfig.Empty() {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Warning,
				"Missing backend configuration",
				`-backend-config was used without a "backend" block in the configuration.

If you intended to override the default local backend configuration,
no action is required, but you may add an explicit backend block to your
configuration to clear this warning:

terraform {
  backend "local" {}
}

However, if you intended to override a defined backend, please verify that
the backend configuration is present and valid.
`,
			))
		}
	}

	opts := &BackendOpts{
		Config:         backendConfig,
		ConfigOverride: backendConfigOverride,
		Init:           true,
	}

	back, backDiags := c.Backend(opts)
	diags = diags.Append(backDiags)
	return back, true, diags
}

// Load the complete module tree, and fetch any missing providers.
// This method outputs its own Ui.
func (c *InitCommand) getProviders(earlyConfig *earlyconfig.Config, state *states.State, upgrade bool) (output bool, diags tfdiags.Diagnostics) {
	// First we'll collect all the provider dependencies we can see in the
	// configuration and the state.
	reqs, moreDiags := earlyConfig.ProviderRequirements()
	diags = diags.Append(moreDiags)
	if moreDiags.HasErrors() {
		return false, diags
	}
	if state != nil {
		stateReqs := state.ProviderRequirements()
		reqs = reqs.Merge(stateReqs)
	}

	// TODO: If the user gave at least one -plugin-dir option on the command
	// line, we should construct a one-off getproviders.Source that consults
	// only those directories and use that instead of c.providerInstallSource()
	// here.
	targetDir := c.providerLocalCacheDir()
	globalCacheDir := c.providerGlobalCacheDir()
	source := c.providerInstallSource()
	inst := providercache.NewInstaller(targetDir, source)
	if globalCacheDir != nil {
		inst.SetGlobalCacheDir(globalCacheDir)
	}

	// Because we're currently just streaming a series of events sequentially
	// into the terminal, we're showing only a subset of the events to keep
	// things relatively concise. Later it'd be nice to have a progress UI
	// where statuses update in-place, but we can't do that as long as we
	// are shimming our vt100 output to the legacy console API on Windows.
	evts := &providercache.InstallerEvents{
		PendingProviders: func(reqs map[addrs.Provider]getproviders.VersionConstraints) {
			c.Ui.Output(c.Colorize().Color(
				"\n[reset][bold]Initializing provider plugins...",
			))
		},
		ProviderAlreadyInstalled: func(provider addrs.Provider, selectedVersion getproviders.Version) {
			c.Ui.Info(fmt.Sprintf("- Using previously-installed %s v%s", provider, selectedVersion))
		},
		QueryPackagesBegin: func(provider addrs.Provider, versionConstraints getproviders.VersionConstraints) {
			if len(versionConstraints) > 0 {
				c.Ui.Info(fmt.Sprintf("- Finding %s versions matching %q...", provider, getproviders.VersionConstraintsString(versionConstraints)))
			} else {
				c.Ui.Info(fmt.Sprintf("- Finding latest version of %s...", provider))
			}
		},
		LinkFromCacheBegin: func(provider addrs.Provider, version getproviders.Version, cacheRoot string) {
			c.Ui.Info(fmt.Sprintf("- Using %s v%s from the shared cache directory", provider, version))
		},
		FetchPackageBegin: func(provider addrs.Provider, version getproviders.Version, location getproviders.PackageLocation) {
			c.Ui.Info(fmt.Sprintf("- Installing %s v%s...", provider, version))
		},
		QueryPackagesFailure: func(provider addrs.Provider, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to query available provider packages",
				fmt.Sprintf("Could not retrieve the list of available versions for provider %s: %s.", provider, err),
			))
		},
		LinkFromCacheFailure: func(provider addrs.Provider, version getproviders.Version, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to install provider from shared cache",
				fmt.Sprintf("Error while importing %s v%s from the shared cache directory: %s.", provider, version, err),
			))
		},
		FetchPackageFailure: func(provider addrs.Provider, version getproviders.Version, err error) {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Failed to install provider",
				fmt.Sprintf("Error while installing %s v%s: %s.", provider, version, err),
			))
		},
	}

	mode := providercache.InstallNewProvidersOnly
	if upgrade {
		mode = providercache.InstallUpgrades
	}
	// TODO: Use a context that will be cancelled when the Terraform
	// process receives SIGINT.
	ctx := evts.OnContext(context.TODO())
	_, err := inst.EnsureProviderVersions(ctx, reqs, mode)
	if err != nil {
		// The errors captured in "err" should be redundant with what we
		// received via the InstallerEvents callbacks above, so we'll
		// just return those as long as we have some.
		if !diags.HasErrors() {
			diags = diags.Append(err)
		}
		return true, diags
	}

	return true, diags

	// TODO: Write the selections into the plugins lock file so we can be
	// sure that future commands will use exactly those provider packages.
	// TODO: Emit constraint suggestions for unconstrained providers.
	/*
		// With all the providers downloaded, we'll generate our lock file
		// that ensures the provider binaries remain unchanged until we init
		// again. If anything changes, other commands that use providers will
		// fail with an error instructing the user to re-run this command.
		available = c.providerPluginSet() // re-discover to see newly-installed plugins

		// internal providers were already filtered out, since we don't need to get them.
		chosen := chooseProviders(available, nil, requirements)

		digests := map[string][]byte{}
		for name, meta := range chosen {
			digest, err := meta.SHA256()
			if err != nil {
				diags = diags.Append(fmt.Errorf("Failed to read provider plugin %s: %s", meta.Path, err))
				return true, diags
			}
			digests[name] = digest
			if c.ignorePluginChecksum {
				digests[name] = nil
			}
		}
		err := c.providerPluginsLock().Write(digests)
		if err != nil {
			diags = diags.Append(fmt.Errorf("failed to save provider manifest: %s", err))
			return true, diags
		}

		{
			// Purge any auto-installed plugins that aren't being used.
			purged, err := c.providerInstaller.PurgeUnused(chosen)
			if err != nil {
				// Failure to purge old plugins is not a fatal error
				c.Ui.Warn(fmt.Sprintf("failed to purge unused plugins: %s", err))
			}
			if purged != nil {
				for meta := range purged {
					log.Printf("[DEBUG] Purged unused %s plugin %s", meta.Name, meta.Path)
				}
			}
		}

		// If any providers have "floating" versions (completely unconstrained)
		// we'll suggest the user constrain with a pessimistic constraint to
		// avoid implicitly adopting a later major release.
		constraintSuggestions := make(map[string]discovery.ConstraintStr)
		for name, meta := range chosen {
			req := requirements[name]
			if req == nil {
				// should never happen, but we don't want to crash here, so we'll
				// be cautious.
				continue
			}

			if req.Versions.Unconstrained() && meta.Version != discovery.VersionZero {
				// meta.Version.MustParse is safe here because our "chosen" metas
				// were already filtered for validity of versions.
				constraintSuggestions[name] = meta.Version.MustParse().MinorUpgradeConstraintStr()
			}
		}
		if len(constraintSuggestions) != 0 {
			names := make([]string, 0, len(constraintSuggestions))
			for name := range constraintSuggestions {
				names = append(names, name)
			}
			sort.Strings(names)

			c.Ui.Output(outputInitProvidersUnconstrained)
			for _, name := range names {
				c.Ui.Output(fmt.Sprintf("* provider.%s: version = %q", name, constraintSuggestions[name]))
			}
		}

		return true, diags
	*/
}

// backendConfigOverrideBody interprets the raw values of -backend-config
// arguments into a hcl Body that should override the backend settings given
// in the configuration.
//
// If the result is nil then no override needs to be provided.
//
// If the returned diagnostics contains errors then the returned body may be
// incomplete or invalid.
func (c *InitCommand) backendConfigOverrideBody(flags rawFlags, schema *configschema.Block) (hcl.Body, tfdiags.Diagnostics) {
	items := flags.AllItems()
	if len(items) == 0 {
		return nil, nil
	}

	var ret hcl.Body
	var diags tfdiags.Diagnostics
	synthVals := make(map[string]cty.Value)

	mergeBody := func(newBody hcl.Body) {
		if ret == nil {
			ret = newBody
		} else {
			ret = configs.MergeBodies(ret, newBody)
		}
	}
	flushVals := func() {
		if len(synthVals) == 0 {
			return
		}
		newBody := configs.SynthBody("-backend-config=...", synthVals)
		mergeBody(newBody)
		synthVals = make(map[string]cty.Value)
	}

	if len(items) == 1 && items[0].Value == "" {
		// Explicitly remove all -backend-config options.
		// We do this by setting an empty but non-nil ConfigOverrides.
		return configs.SynthBody("-backend-config=''", synthVals), diags
	}

	for _, item := range items {
		eq := strings.Index(item.Value, "=")

		if eq == -1 {
			// The value is interpreted as a filename.
			newBody, fileDiags := c.loadHCLFile(item.Value)
			diags = diags.Append(fileDiags)
			flushVals() // deal with any accumulated individual values first
			mergeBody(newBody)
		} else {
			name := item.Value[:eq]
			rawValue := item.Value[eq+1:]
			attrS := schema.Attributes[name]
			if attrS == nil {
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid backend configuration argument",
					fmt.Sprintf("The backend configuration argument %q given on the command line is not expected for the selected backend type.", name),
				))
				continue
			}
			value, valueDiags := configValueFromCLI(item.String(), rawValue, attrS.Type)
			diags = diags.Append(valueDiags)
			if valueDiags.HasErrors() {
				continue
			}
			synthVals[name] = value
		}
	}

	flushVals()

	return ret, diags
}

func (c *InitCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictDirs("")
}

func (c *InitCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-backend":        completePredictBoolean,
		"-backend-config": complete.PredictFiles("*.tfvars"), // can also be key=value, but we can't "predict" that
		"-force-copy":     complete.PredictNothing,
		"-from-module":    completePredictModuleSource,
		"-get":            completePredictBoolean,
		"-get-plugins":    completePredictBoolean,
		"-input":          completePredictBoolean,
		"-lock":           completePredictBoolean,
		"-lock-timeout":   complete.PredictAnything,
		"-no-color":       complete.PredictNothing,
		"-plugin-dir":     complete.PredictDirs(""),
		"-reconfigure":    complete.PredictNothing,
		"-upgrade":        completePredictBoolean,
		"-verify-plugins": completePredictBoolean,
	}
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: terraform init [options] [DIR]

  Initialize a new or existing Terraform working directory by creating
  initial files, loading any remote state, downloading modules, etc.

  This is the first command that should be run for any new or existing
  Terraform configuration per machine. This sets up all the local data
  necessary to run Terraform that is typically not committed to version
  control.

  This command is always safe to run multiple times. Though subsequent runs
  may give errors, this command will never delete your configuration or
  state. Even so, if you have important information, please back it up prior
  to running this command, just in case.

  If no arguments are given, the configuration in this working directory
  is initialized.

Options:

  -backend=true        Configure the backend for this configuration.

  -backend-config=path This can be either a path to an HCL file with key/value
                       assignments (same format as terraform.tfvars) or a
                       'key=value' format. This is merged with what is in the
                       configuration file. This can be specified multiple
                       times. The backend type must be in the configuration
                       itself.

  -force-copy          Suppress prompts about copying state data. This is
                       equivalent to providing a "yes" to all confirmation
                       prompts.

  -from-module=SOURCE  Copy the contents of the given module into the target
                       directory before initialization.

  -get=true            Download any modules for this configuration.

  -get-plugins=true    Download any missing plugins for this configuration.

  -input=true          Ask for input if necessary. If false, will error if
                       input was required.

  -lock=true           Lock the state file when locking is supported.

  -lock-timeout=0s     Duration to retry a state lock.

  -no-color            If specified, output won't contain any color.

  -plugin-dir          Directory containing plugin binaries. This overrides all
                       default search paths for plugins, and prevents the 
                       automatic installation of plugins. This flag can be used
                       multiple times.

  -reconfigure         Reconfigure the backend, ignoring any saved
                       configuration.

  -upgrade=false       If installing modules (-get) or plugins (-get-plugins),
                       ignore previously-downloaded objects and install the
                       latest version allowed within configured constraints.

  -verify-plugins=true Verify the authenticity and integrity of automatically
                       downloaded plugins.
`
	return strings.TrimSpace(helpText)
}

func (c *InitCommand) Synopsis() string {
	return "Initialize a Terraform working directory"
}

const errInitConfigError = `
There are some problems with the configuration, described below.

The Terraform configuration must be valid before initialization so that
Terraform can determine which modules and providers need to be installed.
`

const errInitConfigErrorMaybeLegacySyntax = `
There are some problems with the configuration, described below.

Terraform found syntax errors in the configuration that prevented full
initialization. If you've recently upgraded to Terraform v0.13 from Terraform
v0.11, this may be because your configuration uses syntax constructs that are no
longer valid, and so must be updated before full initialization is possible.

Manually update your configuration syntax, or install Terraform v0.12 and run
terraform init for this configuration at a shell prompt for more information
on how to update it for Terraform v0.12+ compatibility.
`

const errInitCopyNotEmpty = `
The working directory already contains files. The -from-module option requires
an empty directory into which a copy of the referenced module will be placed.

To initialize the configuration already in this working directory, omit the
-from-module option.
`

const outputInitEmpty = `
[reset][bold]Terraform initialized in an empty directory![reset]

The directory has no Terraform configuration files. You may begin working
with Terraform immediately by creating Terraform configuration files.
`

const outputInitSuccess = `
[reset][bold][green]Terraform has been successfully initialized![reset][green]
`

const outputInitSuccessCLI = `[reset][green]
You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
`

const outputInitProvidersUnconstrained = `
The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.
`

const errDiscoveryServiceUnreachable = `
[reset][bold][red]Registry service unreachable.[reset][red]

This may indicate a network issue, or an issue with the requested Terraform Registry.
`

const errProviderNotFound = `
[reset][bold][red]Provider %[1]q not available for installation.[reset][red]

A provider named %[1]q could not be found in the Terraform Registry.

This may result from mistyping the provider name, or the given provider may
be a third-party provider that cannot be installed automatically.

In the latter case, the plugin must be installed manually by locating and
downloading a suitable distribution package and placing the plugin's executable
file in the following directory:
    %[2]s

Terraform detects necessary plugins by inspecting the configuration and state.
To view the provider versions requested by each module, run
"terraform providers".
`

const errProviderVersionsUnsuitable = `
[reset][bold][red]No provider %[1]q plugins meet the constraint %[2]q.[reset][red]

The version constraint is derived from the "version" argument within the
provider %[1]q block in configuration. Child modules may also apply
provider version constraints. To view the provider versions requested by each
module in the current configuration, run "terraform providers".

To proceed, the version constraints for this provider must be relaxed by
either adjusting or removing the "version" argument in the provider blocks
throughout the configuration.
`

const errProviderIncompatible = `
[reset][bold][red]No available provider %[1]q plugins are compatible with this Terraform version.[reset][red]

From time to time, new Terraform major releases can change the requirements for
plugins such that older plugins become incompatible.

Terraform checked all of the plugin versions matching the given constraint:
    %[2]s

Unfortunately, none of the suitable versions are compatible with this version
of Terraform. If you have recently upgraded Terraform, it may be necessary to
move to a newer major release of this provider. Alternatively, if you are
attempting to upgrade the provider to a new major version you may need to
also upgrade Terraform to support the new version.

Consult the documentation for this provider for more information on
compatibility between provider versions and Terraform versions.
`

const errProviderInstallError = `
[reset][bold][red]Error installing provider %[1]q: %[2]s.[reset][red]

Terraform analyses the configuration and state and automatically downloads
plugins for the providers used. However, when attempting to download this
plugin an unexpected error occurred.

This may be caused if for some reason Terraform is unable to reach the
plugin repository. The repository may be unreachable if access is blocked
by a firewall.

If automatic installation is not possible or desirable in your environment,
you may alternatively manually install plugins by downloading a suitable
distribution package and placing the plugin's executable file in the
following directory:
    %[3]s
`

const errMissingProvidersNoInstall = `
[reset][bold][red]Missing required providers.[reset][red]

The following provider constraints are not met by the currently-installed
provider plugins:

%[1]s
Terraform can automatically download and install plugins to meet the given
constraints, but this step was skipped due to the use of -get-plugins=false
and/or -plugin-dir on the command line.

If automatic installation is not possible or desirable in your environment,
you may manually install plugins by downloading a suitable distribution package
and placing the plugin's executable file in one of the directories given in
by -plugin-dir on the command line, or in the following directory if custom
plugin directories are not set:
    %[2]s
`

const errChecksumVerification = `
[reset][bold][red]Error verifying checksum for provider %[1]q[reset][red]
The checksum for provider distribution from the Terraform Registry
did not match the source. This may mean that the distributed files
were changed after this version was released to the Registry.
`

const errSignatureVerification = `
[reset][bold][red]Error:[reset][bold] Untrusted signing key for provider %[1]q[reset]

This provider package is not signed with the HashiCorp signing key, and is
therefore incompatible with Terraform v%[2]s.

A later version of Terraform may have introduced other signing keys that would
accept this provider. Alternatively, an earlier version of this provider may
be compatible with Terraform v%[2]s.
`
