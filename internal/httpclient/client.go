// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpclient

import (
	"github.com/placeholderplaceholderplaceholder/opentf/version"
	"net/http"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// New returns the DefaultPooledClient from the cleanhttp
// package that will also send a Terraform User-Agent string.
func New() *http.Client {
	cli := cleanhttp.DefaultPooledClient()
	cli.Transport = &userAgentRoundTripper{
		userAgent: TerraformUserAgent(version.Version),
		inner:     cli.Transport,
	}
	return cli
}
