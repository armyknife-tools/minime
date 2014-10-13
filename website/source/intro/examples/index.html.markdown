---
layout: "intro"
page_title: "Example Configurations"
sidebar_current: "examples"
---

# Example Configurations

These examples are designed to help you understand some
of the ways Terraform can be used.

All examples in this section are ready to run as-is. Terraform will
ask for input of things such as variables and API keys. If you want to
conitnue using the example, you should save those parameters in a
"terraform.tfvars" file or in a `provider` config bock.

<div class="alert alert-block alert-warning">
<div>
<strong>Note:</strong> The examples use real providers that launch
<em>real</em> resources. That means they can cost money to
experiment with. To avoid unexpected charges, be sure to understand the price
of resources before launching them, and verify any unneeded resources
are cleaned up afterwards.</div>
</div>

Experimenting in this way can help you learn how the Terraform lifecycle
works, as well as how to repeatedly create and destroy infrastructure.

If you're completely new to Terraform, we recommend reading the
[getting started guide](/intro/getting-started/install.html) before diving into
the examples. However, due to the intuitive configuration Terraform
uses it isn't required.

To use these examples, Terraform must first be installed on your machine.
You can install Terraform from the [downloads page](/downloads.html).
