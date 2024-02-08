// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"

	svchost "github.com/hashicorp/terraform-svchost"
	svcauth "github.com/hashicorp/terraform-svchost/auth"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/opentofu/opentofu/internal/command/cliconfig"
)

func TestLogout(t *testing.T) {
	workDir := t.TempDir()

	ui := cli.NewMockUi()
	credsSrc := cliconfig.EmptyCredentialsSourceForTests(filepath.Join(workDir, "credentials.tfrc.json"))

	c := &LogoutCommand{
		Meta: Meta{
			Ui:       ui,
			Services: disco.NewWithCredentialsSource(credsSrc),
		},
	}

	t.Run("with no hostname", func(t *testing.T) {
		status := c.Run([]string{})

		if status != 1 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		if !strings.Contains(ui.ErrorWriter.String(), "The logout command expects exactly one argument") {
			t.Errorf("unexpected error message: %s", ui.ErrorWriter.String())
		}
	})

	testCases := []struct {
		// Hostname to associate a pre-stored token
		hostname string
		// Command-line arguments
		args []string
		// true iff the token at hostname should be removed by the command
		shouldRemove bool
	}{
		// Can remove token for a hostname
		{"tfe.example.com", []string{"tfe.example.com"}, true},

		// Logout does not remove tokens for other hostnames
		{"tfe.example.com", []string{"other-tfe.acme.com"}, false},
	}

	for _, tc := range testCases {
		host := svchost.Hostname(tc.hostname)
		token := svcauth.HostCredentialsToken("some-token")
		err := credsSrc.StoreForHost(host, token)
		if err != nil {
			t.Fatalf("unexpected error storing credentials: %s", err)
		}

		status := c.Run(tc.args)
		if status != 0 {
			t.Fatalf("unexpected error code %d\nstderr:\n%s", status, ui.ErrorWriter.String())
		}

		creds, err := credsSrc.ForHost(host)
		if err != nil {
			t.Errorf("failed to retrieve credentials: %s", err)
		}
		if tc.shouldRemove {
			if creds != nil {
				t.Errorf("wrong token %q; should have no token", creds.Token())
			}
		} else {
			if got, want := creds.Token(), "some-token"; got != want {
				t.Errorf("wrong token %q; want %q", got, want)
			}
		}
	}
}
