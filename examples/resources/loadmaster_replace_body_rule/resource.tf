resource "loadmaster_replace_body_rule" "name" {
  id = "add-example-replace-body"

  pattern     = "A"
  replacement = "B"
}
