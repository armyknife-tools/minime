---
layout: "docs"
page_title: "Command: add"
sidebar_current: "docs-commands-add"
description: |-
  The `terraform add` command generates resource configuration templates.
---

# Command: add

The `terraform add` command generates a starting point for the configuration
of a particular resource.

~> **Warning:** This command is currently experimental. Its exact behavior and
command line arguments are likely to change in future releases based on
feedback. We don't recommend building automation around the current design of
this command, but it's safe to use directly in a development environment
setting.

By default, Terraform will include only the subset of arguments that are marked
by the provider as being required, and will use `null` as a placeholder for
their values. You can then replace `null` with suitable expressions in order
to make the arguments valid.

If you use the `-optional` option then Terraform will also include arguments
that the provider declares as optional. You can then either write a suitable
expression for each argument or remove the arguments you wish to leave unset.

If you use the `-from-state` option then Terraform will instead generate a
configuration containing expressions which will produce the same values as
the corresponding resource instance object already tracked in the Terraform
state, if for example you've previously imported the object using
[`terraform import`](import.html).

-> **Note:** If you use `-from-state`, the result will not include expressions
for any values which are marked as sensitive in the state. If you want to
see those, you can inspect the state data directly using
`terraform state show ADDRESS`.

## Usage

Usage: `terraform add [options] ADDRESS`

This command requires an address that points to a resource which does not
already exist in the configuration. Addresses are in 
[resource addressing format](/docs/cli/state/resource-addressing.html).

This command accepts the following options:

* `-from-state` - Fill the template with values from an existing resource
  instance already tracked in the state. By default, Terraform will emit only
  placeholder values based on the resource type.

* `-optional` - Include optional arguments. By default, the result will
  include only required arguments.

* `-out=FILENAME` - Write the template to a file, instead of to standard
  output.

* `-provider=provider` - Override the provider configuration for the resource,
using the absolute provider configuration address syntax.

    Absolute provider configuration syntax uses the full source address of
    the provider, rather than a local name declared in the relevant module.
    For example, to select the aliased provider configuration "us-east-1"
    of the official AWS provider, use:

    ```
    -provider='provider["hashicorp/aws"].us-east-1'
    ```

    or, if you are using the Windows command prompt, use Windows-style escaping
    for the quotes inside the address:

    ```
    -provider=provider[\"hashicorp/aws\"].us-east-1
    ```

    This is incompatible with `-from-state`, because in that case Terraform
    will use the provider configuration already selected in the state, which
    is the provider configuration that most recently managed the object.
