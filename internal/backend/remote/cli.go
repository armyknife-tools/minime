// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package remote

import (
	"github.com/opentofu/opentofu/internal/backend"
)

// CLIInit implements backend.CLI
func (b *Remote) CLIInit(opts *backend.CLIOpts) error {
	if cli, ok := b.local.(backend.CLI); ok {
		if err := cli.CLIInit(opts); err != nil {
			return err
		}
	}

	b.CLI = opts.CLI
	b.CLIColor = opts.CLIColor
	b.ContextOpts = opts.ContextOpts

	return nil
}
