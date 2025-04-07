resource "loadmaster_modify_url_rule" "name" {
  id = "modify-url-example"

  pattern     = "example-pattern"
  replacement = "ExampleValue"
}
