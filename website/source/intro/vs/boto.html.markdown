---
layout: "intro"
page_title: "Terraform vs. Boto, Fog, etc."
sidebar_current: "vs-other-boto"
---

# Terraform vs. Boto, Fog, etc.

Libraries like Boto, Fog, etc. are used to provide native access
clients to cloud providers and services by using their APIs. Some
libraries are focused on specific clouds, while others attempt
to bridge them all and mask the semantic differences. Using a client
library only provides low-level access to APIs, requiring application
developers to build their own tooling to build and manage their infrastructure.

Terraform is not intended to give low-level programmatic access to
providers, but instead provides a high level syntax for describing
how cloud resources and services should be created, provisioned, and
combined.  It acts as a feature rich and flexible tool, using a
a plugin-based model to support providers and provisioners, giving
it the ability to support almost any service that exposes APIs.

