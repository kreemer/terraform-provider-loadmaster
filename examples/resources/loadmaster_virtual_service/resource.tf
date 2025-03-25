resource "loadmaster_virtual_service" "example" {
  address  = "10.0.0.4"
  port     = "8889"
  protocol = "tcp"
}
