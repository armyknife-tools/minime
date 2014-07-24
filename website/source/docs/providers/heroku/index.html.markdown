---
layout: "heroku"
page_title: "Provider: Heroku"
sidebar_current: "docs-heroku-index"
---

# Heroku Provider

The Heroku provider is used to interact with the
resources supported by Heroku. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```
# Configure the Heroku provider
provider "heroku" {
    email = "ops@company.com"
	api_key = "${var.heroku_api_key}"
}

# Create a new applicaiton
resource "heroku_app" "default" {
    ...
}
```

## Argument Reference

The following arguments are supported:

* `api_key` - (Required) Heroku API token
* `email` - (Required) Email to be notified by Heroku

