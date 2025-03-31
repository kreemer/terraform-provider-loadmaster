resource "loadmaster_virtual_service" "test" {
  address  = "10.0.0.4"
  port     = "9090"
  protocol = "tcp"
}

resource "loadmaster_real_server" "test" {
  virtual_service_id = loadmaster_virtual_service.test.id
  address            = "10.0.0.99"
  port               = "80"
  protocol           = "tcp"
}
