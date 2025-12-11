resource "loadmaster_owasp_custom_rule" "name" {
  filename = "rule.conf"
  data     = file("${path.module}/rule.conf")
}
