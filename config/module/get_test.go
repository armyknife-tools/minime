package module

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/registry/response"
)

// Map of module names and location of test modules.
// Only one version for now, as we only lookup latest from the registry.
type testMod struct {
	location string
	version  string
}

// All the locationes from the mockRegistry start with a file:// scheme. If
// the the location string here doesn't have a scheme, the mockRegistry will
// find the absolute path and return a complete URL.
var testMods = map[string][]testMod{
	"registry/foo/bar": {{
		location: "file:///download/registry/foo/bar/0.2.3//*?archive=tar.gz",
		version:  "0.2.3",
	}},
	"registry/foo/baz": {{
		location: "file:///download/registry/foo/baz/1.10.0//*?archive=tar.gz",
		version:  "1.10.0",
	}},
	"registry/local/sub": {{
		location: "test-fixtures/registry-tar-subdir/foo.tgz//*?archive=tar.gz",
		version:  "0.1.2",
	}},
	"exists-in-registry/identifier/provider": {{
		location: "file:///registry/exists",
		version:  "0.2.0",
	}},
	"test-versions/name/provider": {
		{version: "2.2.0"},
		{version: "2.1.1"},
		{version: "1.2.2"},
		{version: "1.2.1"},
	},
}

func latestVersion(versions []string) string {
	var col version.Collection
	for _, v := range versions {
		ver, err := version.NewVersion(v)
		if err != nil {
			panic(err)
		}
		col = append(col, ver)
	}

	sort.Sort(col)
	return col[len(col)-1].String()
}

func mockRegHandler() http.Handler {
	mux := http.NewServeMux()

	download := func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimLeft(r.URL.Path, "/")
		// handle download request
		re := regexp.MustCompile(`^([-a-z]+/\w+/\w+)/download$`)
		// download lookup
		matches := re.FindStringSubmatch(p)
		if len(matches) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		versions, ok := testMods[matches[1]]
		if !ok {
			http.NotFound(w, r)
			return
		}
		mod := versions[0]

		location := mod.location
		if !strings.HasPrefix(location, "file:///") {
			// we can't use filepath.Abs because it will clean `//`
			wd, _ := os.Getwd()
			location = fmt.Sprintf("file://%s/%s", wd, location)
		}

		w.Header().Set("X-Terraform-Get", location)
		w.WriteHeader(http.StatusNoContent)
		// no body
		return
	}

	versions := func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimLeft(r.URL.Path, "/")
		re := regexp.MustCompile(`^([-a-z]+/\w+/\w+)/versions$`)
		matches := re.FindStringSubmatch(p)
		if len(matches) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := matches[1]
		versions, ok := testMods[name]
		if !ok {
			http.NotFound(w, r)
			return
		}

		// only adding the single requested module for now
		// this is the minimal that any regisry is epected to support
		mpvs := &response.ModuleProviderVersions{
			Source: name,
		}

		for _, v := range versions {
			mv := &response.ModuleVersion{
				Version: v.version,
			}
			mpvs.Versions = append(mpvs.Versions, mv)
		}

		resp := response.ModuleVersions{
			Modules: []*response.ModuleProviderVersions{mpvs},
		}

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	mux.Handle("/v1/modules/",
		http.StripPrefix("/v1/modules/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/download") {
				download(w, r)
				return
			}

			if strings.HasSuffix(r.URL.Path, "/versions") {
				versions(w, r)
				return
			}

			http.NotFound(w, r)
		})),
	)

	mux.HandleFunc("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"modules.v1":"/v1/modules/"}`)
	})
	return mux
}

// Just enough like a registry to exercise our code.
// Returns the location of the latest version
func mockRegistry() *httptest.Server {
	server := httptest.NewServer(mockRegHandler())
	return server
}

func mockTLSRegistry() *httptest.Server {
	server := httptest.NewTLSServer(mockRegHandler())
	return server
}

func setResetRegDetector(server *httptest.Server) func() {
	regDetector := &registryDetector{
		api:    server.URL + "/v1/modules",
		client: server.Client(),
	}

	origDetectors := detectors
	detectors = []getter.Detector{
		new(getter.GitHubDetector),
		new(getter.BitBucketDetector),
		new(getter.S3Detector),
		regDetector,
		new(getter.FileDetector),
	}

	return func() {
		detectors = origDetectors
	}
}

func TestDetectRegistry(t *testing.T) {
	server := mockRegistry()
	defer server.Close()

	detector := registryDetector{
		api:    server.URL + "/v1/modules",
		client: server.Client(),
	}

	for _, tc := range []struct {
		source   string
		location string
		found    bool
		err      bool
	}{
		{
			source:   "registry/foo/bar",
			location: testMods["registry/foo/bar"][0].location,
			found:    true,
		},
		{
			source:   "registry/foo/baz",
			location: testMods["registry/foo/baz"][0].location,
			found:    true,
		},
		// this should not be found, and is no longer valid as a local source
		{
			source: "registry/foo/notfound",
			err:    true,
		},

		// a full url should not be detected
		{
			source: "http://example.com/registry/foo/notfound",
			found:  false,
		},

		// paths should not be detected
		{
			source: "./local/foo/notfound",
			found:  false,
		},
		{
			source: "/local/foo/notfound",
			found:  false,
		},

		// wrong number of parts can't be regisry IDs
		{
			source: "something/registry/foo/notfound",
			found:  false,
		},
	} {

		t.Run(tc.source, func(t *testing.T) {
			loc, ok, err := detector.Detect(tc.source, "")
			if (err == nil) == tc.err {
				t.Fatalf("expected error? %t; got error: %v", tc.err, err)
			}

			if ok != tc.found {
				t.Fatalf("expected OK == %t", tc.found)
			}

			loc = strings.TrimPrefix(loc, server.URL+"/")
			if strings.TrimPrefix(loc, server.URL) != tc.location {
				t.Fatalf("expected location: %q, got %q", tc.location, loc)
			}
		})

	}
}

// check that the full set of detectors works as expected
func TestDetectors(t *testing.T) {
	server := mockRegistry()
	defer server.Close()
	defer setResetRegDetector(server)()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		source   string
		location string
		fixture  string
		err      bool
	}{
		{
			source:   "registry/foo/bar",
			location: "file:///download/registry/foo/bar/0.2.3//*?archive=tar.gz",
		},
		// this should not be found, and is no longer a valid local source
		{
			source: "registry/foo/notfound",
			err:    true,
		},
		// a full url should be unchanged
		{
			source: "http://example.com/registry/foo/notfound?" +
				"checksum=sha256:f19056b80a426d797ff9e470da069c171a6c6befa83e2da7f6c706207742acab",
			location: "http://example.com/registry/foo/notfound?" +
				"checksum=sha256:f19056b80a426d797ff9e470da069c171a6c6befa83e2da7f6c706207742acab",
		},

		// forced getters will return untouched
		{
			source:   "git::http://example.com/registry/foo/notfound?param=value",
			location: "git::http://example.com/registry/foo/notfound?param=value",
		},

		// local paths should be detected as such, even if they're match
		// registry modules.
		{
			source:   "./registry/foo/bar",
			location: "file://" + filepath.Join(wd, "registry/foo/bar"),
		},
		{
			source:   "/registry/foo/bar",
			location: "file:///registry/foo/bar",
		},

		// Wrong number of parts can't be registry IDs.
		// This is returned as a local path for now, but may return an error at
		// some point.
		{
			source:   "something/here/registry/foo/notfound",
			location: "file://" + filepath.Join(wd, "something/here/registry/foo/notfound"),
		},

		// make sure a local module that looks like a registry id can be found
		{
			source:  "namespace/identifier/provider",
			fixture: "discover-subdirs",
			err:     true,
		},

		// The registry takes precedence over local paths if they don't start
		// with a relative or absolute path
		{
			source:  "exists-in-registry/identifier/provider",
			fixture: "discover-registry-local",
			// registry should take precidence
			location: "file:///registry/exists",
		},
	} {

		t.Run(tc.source, func(t *testing.T) {
			dir := wd
			if tc.fixture != "" {
				dir = filepath.Join(wd, fixtureDir, tc.fixture)
				if err := os.Chdir(dir); err != nil {
					t.Fatal(err)
				}
				defer os.Chdir(wd)
			}

			loc, err := getter.Detect(tc.source, dir, detectors)
			if (err == nil) == tc.err {
				t.Fatalf("expected error? %t; got error :%v", tc.err, err)
			}

			loc = strings.TrimPrefix(loc, server.URL+"/")
			if strings.TrimPrefix(loc, server.URL) != tc.location {
				t.Fatalf("expected location: %q, got %q", tc.location, loc)
			}
		})

	}
}

// GitHub archives always contain the module source in a single subdirectory,
// so the registry will return a path with with a `//*` suffix. We need to make
// sure this doesn't intefere with our internal handling of `//` subdir.
func TestRegistryGitHubArchive(t *testing.T) {
	server := mockRegistry()
	defer server.Close()
	defer setResetRegDetector(server)()

	storage := testStorage(t)
	tree := NewTree("", testConfig(t, "registry-tar-subdir"))

	if err := tree.Load(storage, GetModeGet); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !tree.Loaded() {
		t.Fatal("should be loaded")
	}

	if err := tree.Load(storage, GetModeNone); err != nil {
		t.Fatalf("err: %s", err)
	}

	// stop the registry server, and make sure that we don't need to call out again
	server.Close()
	tree = NewTree("", testConfig(t, "registry-tar-subdir"))

	if err := tree.Load(storage, GetModeGet); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !tree.Loaded() {
		t.Fatal("should be loaded")
	}

	actual := strings.TrimSpace(tree.String())
	expected := strings.TrimSpace(treeLoadSubdirStr)
	if actual != expected {
		t.Fatalf("got: \n\n%s\nexpected: \n\n%s", actual, expected)
	}
}

// Test that the //subdir notation can be used with registry modules
func TestRegisryModuleSubdir(t *testing.T) {
	server := mockRegistry()
	defer server.Close()
	defer setResetRegDetector(server)()

	storage := testStorage(t)
	tree := NewTree("", testConfig(t, "registry-subdir"))

	if err := tree.Load(storage, GetModeGet); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !tree.Loaded() {
		t.Fatal("should be loaded")
	}

	if err := tree.Load(storage, GetModeNone); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := strings.TrimSpace(tree.String())
	expected := strings.TrimSpace(treeLoadRegistrySubdirStr)
	if actual != expected {
		t.Fatalf("got: \n\n%s\nexpected: \n\n%s", actual, expected)
	}
}

func TestAccRegistryDiscover(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping ACC test")
	}

	// simply check that we get a valid github URL for this from the registry
	loc, err := getter.Detect("hashicorp/consul/aws", "./", detectors)
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(loc)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(u.Host, "github.com") {
		t.Fatalf("expected host 'github.com', got: %q", u.Host)
	}

	if !strings.Contains(u.String(), "consul") {
		t.Fatalf("url doesn't contain 'consul': %s", u.String())
	}
}

func TestAccRegistryLoad(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping ACC test")
	}

	storage := testStorage(t)
	tree := NewTree("", testConfig(t, "registry-load"))

	if err := tree.Load(storage, GetModeGet); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !tree.Loaded() {
		t.Fatal("should be loaded")
	}

	if err := tree.Load(storage, GetModeNone); err != nil {
		t.Fatalf("err: %s", err)
	}

	// TODO expand this further by fetching some metadata from the registry
	actual := strings.TrimSpace(tree.String())
	if !strings.Contains(actual, "(path: vault)") {
		t.Fatal("missing vault module, got:\n", actual)
	}
}
