package arguments

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/addrs"
)

func TestParseApply_basicValid(t *testing.T) {
	testCases := map[string]struct {
		args []string
		want *Apply
	}{
		"defaults": {
			nil,
			&Apply{
				AutoApprove:  false,
				InputEnabled: true,
				PlanPath:     "",
				ViewType:     ViewHuman,
			},
		},
		"auto-approve, disabled input, and plan path": {
			[]string{"-auto-approve", "-input=false", "saved.tfplan"},
			&Apply{
				AutoApprove:  true,
				InputEnabled: false,
				PlanPath:     "saved.tfplan",
				ViewType:     ViewHuman,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, diags := ParseApply(tc.args)
			if len(diags) > 0 {
				t.Fatalf("unexpected diags: %v", diags)
			}
			// Ignore the extended arguments for simplicity
			got.State = nil
			got.Operation = nil
			got.Vars = nil
			if *got != *tc.want {
				t.Fatalf("unexpected result\n got: %#v\nwant: %#v", got, tc.want)
			}
		})
	}
}

func TestParseApply_invalid(t *testing.T) {
	got, diags := ParseApply([]string{"-frob"})
	if len(diags) == 0 {
		t.Fatal("expected diags but got none")
	}
	if got, want := diags.Err().Error(), "flag provided but not defined"; !strings.Contains(got, want) {
		t.Fatalf("wrong diags\n got: %s\nwant: %s", got, want)
	}
	if got.ViewType != ViewHuman {
		t.Fatalf("wrong view type, got %#v, want %#v", got.ViewType, ViewHuman)
	}
}

func TestParseApply_tooManyArguments(t *testing.T) {
	got, diags := ParseApply([]string{"saved.tfplan", "please"})
	if len(diags) == 0 {
		t.Fatal("expected diags but got none")
	}
	if got, want := diags.Err().Error(), "Too many command line arguments"; !strings.Contains(got, want) {
		t.Fatalf("wrong diags\n got: %s\nwant: %s", got, want)
	}
	if got.ViewType != ViewHuman {
		t.Fatalf("wrong view type, got %#v, want %#v", got.ViewType, ViewHuman)
	}
}

func TestParseApply_targets(t *testing.T) {
	foobarbaz, _ := addrs.ParseTargetStr("foo_bar.baz")
	boop, _ := addrs.ParseTargetStr("module.boop")
	testCases := map[string]struct {
		args    []string
		want    []addrs.Targetable
		wantErr string
	}{
		"no targets by default": {
			args: nil,
			want: nil,
		},
		"one target": {
			args: []string{"-target=foo_bar.baz"},
			want: []addrs.Targetable{foobarbaz.Subject},
		},
		"two targets": {
			args: []string{"-target=foo_bar.baz", "-target", "module.boop"},
			want: []addrs.Targetable{foobarbaz.Subject, boop.Subject},
		},
		"invalid traversal": {
			args:    []string{"-target=foo."},
			want:    nil,
			wantErr: "Dot must be followed by attribute name",
		},
		"invalid target": {
			args:    []string{"-target=data[0].foo"},
			want:    nil,
			wantErr: "A data source name is required",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, diags := ParseApply(tc.args)
			if len(diags) > 0 {
				if tc.wantErr == "" {
					t.Fatalf("unexpected diags: %v", diags)
				} else if got := diags.Err().Error(); !strings.Contains(got, tc.wantErr) {
					t.Fatalf("wrong diags\n got: %s\nwant: %s", got, tc.wantErr)
				}
			}
			if !cmp.Equal(got.Operation.Targets, tc.want) {
				t.Fatalf("unexpected result\n%s", cmp.Diff(got.Operation.Targets, tc.want))
			}
		})
	}
}

func TestParseApply_vars(t *testing.T) {
	testCases := map[string]struct {
		args []string
		want []FlagNameValue
	}{
		"no var flags by default": {
			args: nil,
			want: nil,
		},
		"one var": {
			args: []string{"-var", "foo=bar"},
			want: []FlagNameValue{
				{Name: "-var", Value: "foo=bar"},
			},
		},
		"one var-file": {
			args: []string{"-var-file", "cool.tfvars"},
			want: []FlagNameValue{
				{Name: "-var-file", Value: "cool.tfvars"},
			},
		},
		"ordering preserved": {
			args: []string{
				"-var", "foo=bar",
				"-var-file", "cool.tfvars",
				"-var", "boop=beep",
			},
			want: []FlagNameValue{
				{Name: "-var", Value: "foo=bar"},
				{Name: "-var-file", Value: "cool.tfvars"},
				{Name: "-var", Value: "boop=beep"},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, diags := ParseApply(tc.args)
			if len(diags) > 0 {
				t.Fatalf("unexpected diags: %v", diags)
			}
			if vars := got.Vars.All(); !cmp.Equal(vars, tc.want) {
				t.Fatalf("unexpected result\n%s", cmp.Diff(vars, tc.want))
			}
			if got, want := got.Vars.Empty(), len(tc.want) == 0; got != want {
				t.Fatalf("expected Empty() to return %t, but was %t", want, got)
			}
		})
	}
}
