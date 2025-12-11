resource "loadmaster_owasp_custom_data" "name" {
  filename = "data.txt"
  data     = file("${path.module}/data.txt")
}
