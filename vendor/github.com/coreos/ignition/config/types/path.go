// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"errors"
	"path/filepath"

	"github.com/coreos/ignition/config/validate/report"
)

var (
	ErrPathRelative = errors.New("path not absolute")
)

type Path string

func (p Path) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(p) + `"`), nil
}

func (p Path) Validate() report.Report {
	if !filepath.IsAbs(string(p)) {
		return report.ReportFromError(ErrPathRelative, report.EntryError)
	}
	return report.Report{}
}
