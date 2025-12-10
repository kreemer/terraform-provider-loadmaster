resource "loadmaster_owasp_custom_data" "name" {
  filename = "rule.txt"

  data = <<-EOT
    Data
    EOT
}
