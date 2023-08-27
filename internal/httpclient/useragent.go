// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/placeholderplaceholderplaceholder/opentf/version"
)

const userAgentFormat = "OpenTF/%s"
const uaEnvVar = "TF_APPEND_USER_AGENT"

// Deprecated: Use TerraformUserAgent(version) instead
func UserAgentString() string {
	ua := fmt.Sprintf(userAgentFormat, version.Version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", rt.userAgent)
	}
	log.Printf("[TRACE] HTTP client %s request to %s", req.Method, req.URL.String())
	return rt.inner.RoundTrip(req)
}

func TerraformUserAgent(version string) string {
	ua := fmt.Sprintf("placeholderplaceholderplaceholder OpenTF/%s (+https://www.opentf.org)", version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
