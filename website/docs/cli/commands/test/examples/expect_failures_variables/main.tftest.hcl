# Copyright (c) The OpenTofu Authors
# SPDX-License-Identifier: MPL-2.0
# Copyright (c) 2023 HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

run "main" {
  command = plan

  variables {
    instances = -1
  }

  expect_failures = [
    var.instances,
  ]
}