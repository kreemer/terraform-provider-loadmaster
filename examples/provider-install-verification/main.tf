terraform {
  required_providers {
    loadmaster = {
      source = "registry.terraform.io/kreemer/loadmaster"
    }
  }
}

provider "loadmaster" {}

data "loadmaster_virtual_service" "example" {
  id = "1"
}

resource "loadmaster_virtual_service" "example" {
  address  = "10.0.0.4"
  port     = "8889"
  protocol = "tcp"

}

output "virtual_service" {
  value = data.loadmaster_virtual_service.example
}


output "virtual_service2" {
  value = loadmaster_virtual_service.example
}
