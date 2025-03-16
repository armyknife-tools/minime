// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"strings"

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
    refresh    Refresh local registry module and provider cache
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

	c.Meta.Ui.Error(
		"The provided subcommand wasn't found.\n" +
			c.Help())
	return 1
}
