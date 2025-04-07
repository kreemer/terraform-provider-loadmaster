resource "loadmaster_delete_header_rule" "name" {
  id = "delete-example-header"

  header      = "X-Example-Header"
  pattern     = "A"
  replacement = "B"
}
