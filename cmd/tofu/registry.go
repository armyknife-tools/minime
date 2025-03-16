// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"strings"

	"github.com/mitchellh/cli"
	"github.com/opentofu/opentofu/internal/command"
)

// RegistryCommand is a container for registry subcommands
type RegistryCommand struct {
	Meta command.Meta
}

func (c *RegistryCommand) Help() string {
	helpText := `
Usage: tofu registry <subcommand> [options] [args]

  This command has subcommands for registry operations.

Subcommands:
    provider    Provider registry operations
    refresh     Refresh local registry module and provider cache
    search      Search the registry for modules or providers
`
	return strings.TrimSpace(helpText)
}

func (c *RegistryCommand) Synopsis() string {
	return "Registry operations"
}

func (c *RegistryCommand) Run(args []string) int {
	if len(args) == 0 {
		c.Meta.Ui.Error("A subcommand is expected.\n" + c.Help())
		return 1
	}

	switch args[0] {
	case "provider":
		return cli.RunResultHelp
	case "refresh":
		return cli.RunResultHelp
	case "search":
		return cli.RunResultHelp
	default:
		c.Meta.Ui.Error(
			"The provided subcommand wasn't found.\n" +
				c.Help())
		return 1
	}
}
