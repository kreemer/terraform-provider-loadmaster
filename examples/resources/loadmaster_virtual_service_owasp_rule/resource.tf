resource "loadmaster_virtual_service_owasp_rule" "example" {
  virtual_service_id = loadmaster_virtual_service.example.id
  rule               = <<-EOT
    [...]
    EOT
  run_first          = false
}
