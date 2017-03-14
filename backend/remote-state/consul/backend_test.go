package consul

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/testutil"
	"github.com/hashicorp/terraform/backend"
)

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func newConsulTestServer(t *testing.T) *testutil.TestServer {
	skip := os.Getenv("TF_ACC") == "" && os.Getenv("TF_CONSUL_TEST") == ""
	if skip {
		t.Log("consul server tests require setting TF_ACC or TF_CONSUL_TEST")
		t.Skip()
	}

	srv := testutil.NewTestServerConfig(t, func(c *testutil.TestServerConfig) {
		c.LogLevel = "warn"

		if !testing.Verbose() {
			c.Stdout = ioutil.Discard
			c.Stderr = ioutil.Discard
		}
	})

	return srv
}

func TestBackend(t *testing.T) {
	srv := newConsulTestServer(t)
	defer srv.Stop()

	// Get the backend
	b := backend.TestBackendConfig(t, New(), map[string]interface{}{
		"address": srv.HTTPAddr,
		"path":    fmt.Sprintf("tf-unit/%s", time.Now().String()),
	})

	// Test
	backend.TestBackend(t, b)
}
