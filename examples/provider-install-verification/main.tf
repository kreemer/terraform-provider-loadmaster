terraform {
  required_providers {
    loadmaster = {
      source = "registry.terraform.io/kreemer/loadmaster"
    }
  }
}

provider "loadmaster" {}

resource "loadmaster_virtual_service" "example" {
  address  = "10.0.0.4"
  port     = "9123"
  protocol = "tcp"

  nickname = "master"
  enabled  = true
  standard_options = {
    scheduling_method = "rr"
  }
}

resource "loadmaster_sub_virtual_service" "example1" {
  virtual_service_id = loadmaster_virtual_service.example.id
  type               = "http"
  nickname           = "subvs1"
  standard_options = {
    scheduling_method = "rr"
  }
}
resource "loadmaster_sub_virtual_service" "example2" {
  virtual_service_id = loadmaster_virtual_service.example.id
  type               = "gen"
  nickname           = "subvs2"
  standard_options = {
    scheduling_method = "wrr"
  }
}
