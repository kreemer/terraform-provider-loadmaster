resource "loadmaster_match_content_rule" "name" {
  id         = "match-url-by-domainname"
  pattern    = "^(.*)www.google.com(.*)$"
  match_type = "regex"
  inc_host   = true
  no_case    = true
}
