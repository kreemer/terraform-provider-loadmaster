resource "loadmaster_owasp_custom_rule" "name" {
  filename = "rule"

  data = <<-EOT
    # OWASP Custom Rule
    SecRule REQUEST_HEADERS:User-Agent "BadBot" "id:12345,phase:1,deny,status:403,msg:'Bad Bot Detected'"
    EOT
}
