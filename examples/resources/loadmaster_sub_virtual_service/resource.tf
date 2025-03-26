resource "loadmaster_virtual_service" "example" {
  address  = "10.0.0.4"
  port     = "8889"
  protocol = "tcp"
}

resource "loadmaster_sub_virtual_service" "example" {
  virtual_service_id = loadmaster_virtual_service.example.id

  nickname = "example"
}
