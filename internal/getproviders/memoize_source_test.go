package getproviders

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/addrs"
)

func TestMemoizeSource(t *testing.T) {
	provider := addrs.NewDefaultProvider("foo")
	version := MustParseVersion("1.0.0")
	protocols := VersionList{MustParseVersion("5.0")}
	platform := Platform{OS: "gameboy", Arch: "lr35902"}
	meta := FakePackageMeta(provider, version, protocols, platform)
	nonexistProvider := addrs.NewDefaultProvider("nonexist")
	nonexistPlatform := Platform{OS: "gamegear", Arch: "z80"}

	t.Run("AvailableVersions for existing provider", func(t *testing.T) {
		mock := NewMockSource([]PackageMeta{meta})
		source := NewMemoizeSource(mock)

		got, err := source.AvailableVersions(provider)
		want := VersionList{version}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from first call to AvailableVersions\n%s", diff)
		}

		got, err = source.AvailableVersions(provider)
		want = VersionList{version}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from second call to AvailableVersions\n%s", diff)
		}

		_, err = source.AvailableVersions(nonexistProvider)
		if want, ok := err.(ErrProviderNotKnown); !ok {
			t.Fatalf("wrong error type from nonexist call:\ngot:  %T\nwant: %T", err, want)
		}

		got, err = source.AvailableVersions(provider)
		want = VersionList{version}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from third call to AvailableVersions\n%s", diff)
		}

		gotLog := mock.CallLog()
		wantLog := [][]interface{}{
			// Only one call for the main provider, because the others were returned from the cache.
			{"AvailableVersions", provider},

			// The call for nonexist also shows through, because it didn't match the cache.
			{"AvailableVersions", nonexistProvider},
		}
		if diff := cmp.Diff(wantLog, gotLog); diff != "" {
			t.Fatalf("unexpected call log\n%s", diff)
		}
	})
	t.Run("PackageMeta for existing provider", func(t *testing.T) {
		mock := NewMockSource([]PackageMeta{meta})
		source := NewMemoizeSource(mock)

		got, err := source.PackageMeta(provider, version, platform)
		want := meta
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from first call to PackageMeta\n%s", diff)
		}

		got, err = source.PackageMeta(provider, version, platform)
		want = meta
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from second call to PackageMeta\n%s", diff)
		}

		_, err = source.PackageMeta(nonexistProvider, version, platform)
		if want, ok := err.(ErrPlatformNotSupported); !ok {
			t.Fatalf("wrong error type from nonexist provider call:\ngot:  %T\nwant: %T", err, want)
		}
		_, err = source.PackageMeta(provider, version, nonexistPlatform)
		if want, ok := err.(ErrPlatformNotSupported); !ok {
			t.Fatalf("wrong error type from nonexist platform call:\ngot:  %T\nwant: %T", err, want)
		}

		got, err = source.PackageMeta(provider, version, platform)
		want = meta
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("wrong result from third call to PackageMeta\n%s", diff)
		}

		gotLog := mock.CallLog()
		wantLog := [][]interface{}{
			// Only one call for the main provider, because the others were returned from the cache.
			{"PackageMeta", provider, version, platform},

			// The other calls for non-exist things also show through, because they missed the cache.
			{"PackageMeta", nonexistProvider, version, platform},
			{"PackageMeta", provider, version, nonexistPlatform},
		}
		if diff := cmp.Diff(wantLog, gotLog); diff != "" {
			t.Fatalf("unexpected call log\n%s", diff)
		}
	})
	t.Run("AvailableVersions for non-existing provider", func(t *testing.T) {
		mock := NewMockSource([]PackageMeta{meta})
		source := NewMemoizeSource(mock)

		_, err := source.AvailableVersions(nonexistProvider)
		if want, ok := err.(ErrProviderNotKnown); !ok {
			t.Fatalf("wrong error type from first call:\ngot:  %T\nwant: %T", err, want)
		}
		_, err = source.AvailableVersions(nonexistProvider)
		if want, ok := err.(ErrProviderNotKnown); !ok {
			t.Fatalf("wrong error type from second call:\ngot:  %T\nwant: %T", err, want)
		}

		gotLog := mock.CallLog()
		wantLog := [][]interface{}{
			// Only one call, because the other was returned from the cache.
			{"AvailableVersions", nonexistProvider},
		}
		if diff := cmp.Diff(wantLog, gotLog); diff != "" {
			t.Fatalf("unexpected call log\n%s", diff)
		}
	})
	t.Run("PackageMeta for non-existing provider", func(t *testing.T) {
		mock := NewMockSource([]PackageMeta{meta})
		source := NewMemoizeSource(mock)

		_, err := source.PackageMeta(nonexistProvider, version, platform)
		if want, ok := err.(ErrPlatformNotSupported); !ok {
			t.Fatalf("wrong error type from first call:\ngot:  %T\nwant: %T", err, want)
		}
		_, err = source.PackageMeta(nonexistProvider, version, platform)
		if want, ok := err.(ErrPlatformNotSupported); !ok {
			t.Fatalf("wrong error type from second call:\ngot:  %T\nwant: %T", err, want)
		}

		gotLog := mock.CallLog()
		wantLog := [][]interface{}{
			// Only one call, because the other was returned from the cache.
			{"PackageMeta", nonexistProvider, version, platform},
		}
		if diff := cmp.Diff(wantLog, gotLog); diff != "" {
			t.Fatalf("unexpected call log\n%s", diff)
		}
	})
}
