terraform {
  required_providers {
    foo = {
      source = "opentofu/foo"
    }
  }
}

module "mod2" {
  source = "./mod1"
  providers = {
    foo = foo
  }
}
