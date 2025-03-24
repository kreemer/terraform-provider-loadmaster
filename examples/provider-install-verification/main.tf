terraform {
  required_providers {
    loadmaster = {
      source = "registry.terraform.io/kreemer/loadmaster"
    }
  }
}

provider "loadmaster" {}
