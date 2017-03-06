package storage

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator 1.0.1.0
// Changes may cause incorrect behavior and will be lost if the code is
// regenerated.

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	major           = "8"
	minor           = "1"
	patch           = "0"
	tag             = "beta"
	userAgentFormat = "Azure-SDK-For-Go/%s arm-%s/%s"
)

// cached results of UserAgent and Version to prevent repeated operations.
var (
	userAgent string
	version   string
)

// UserAgent returns the UserAgent string to use when sending http.Requests.
func UserAgent() string {
	if userAgent == "" {
		userAgent = fmt.Sprintf(userAgentFormat, Version(), "storage", "2016-01-01")
	}
	return userAgent
}

// Version returns the semantic version (see http://semver.org) of the client.
func Version() string {
	if version == "" {
		versionBuilder := bytes.NewBufferString(fmt.Sprintf("%s.%s.%s", major, minor, patch))
		if tag != "" {
			versionBuilder.WriteRune('-')
			versionBuilder.WriteString(strings.TrimPrefix(tag, "-"))
		}
		version = string(versionBuilder.Bytes())
	}
	return version
}
