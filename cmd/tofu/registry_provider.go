package main

import (
	"strings"

	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/command"
	"github.com/posener/complete"
)

// RegistryProviderCommand is the command that just shows help for the subcommands.
type RegistryProviderCommand struct {
	Meta command.Meta
}

func (c *RegistryProviderCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *RegistryProviderCommand) Help() string {
	helpText := `
Usage: tofu registry provider <subcommand> [options] [args]

  This command has subcommands for working with OpenTofu providers.

  The main subcommand is 'install', which installs a provider from a registry
  and adds it to your configuration.

Subcommands:
    install    Install a provider from a registry
`
	return strings.TrimSpace(helpText)
}

func (c *RegistryProviderCommand) Synopsis() string {
	return "Registry provider management"
}

func (c *RegistryProviderCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RegistryProviderCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}
