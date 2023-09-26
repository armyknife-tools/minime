// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package configs contains types that represent OpenTofu configurations and
// the different elements thereof.
//
// The functionality in this package can be used for some static analyses of
// OpenTofu configurations, but this package generally exposes representations
// of the configuration source code rather than the result of evaluating these
// objects. The sibling package "lang" deals with evaluation of structures
// and expressions in the configuration.
//
// Due to its close relationship with HCL, this package makes frequent use
// of types from the HCL API, including raw HCL diagnostic messages. Such
// diagnostics can be converted into OpenTofu-flavored diagnostics, if needed,
// using functions in the sibling package tfdiags.
//
// The Parser type is the main entry-point into this package. The LoadConfigDir
// method can be used to load a single module directory, and then a full
// configuration (including any descendent modules) can be produced using
// the top-level BuildConfig method.
package configs
