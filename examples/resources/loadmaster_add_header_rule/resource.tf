resource "loadmaster_add_header_rule" "name" {
  id = "add-example-header"

  header      = "X-Example-Header"
  replacement = "ExampleValue"
}
