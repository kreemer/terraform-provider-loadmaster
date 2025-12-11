// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestOwaspCustomRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testOwaspCustomRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_rule.test_rule",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("test_rule_replace_url.conf"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_rule.test_rule",
						tfjsonpath.New("data"),
						knownvalue.StringRegexp(regexp.MustCompile(`.*SecMarker BEGIN_ALLOWLIST_login.*`)),
					),
				},
			},
			{
				ResourceName:                         "loadmaster_owasp_custom_rule.test_rule",
				ImportStateVerifyIdentifierAttribute: "filename",
				ImportStateId:                        "test_rule_replace_url.conf",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				Config: testOwaspCustomRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_rule.test_rule",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("test_rule_replace_url.conf"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_rule.test_rule",
						tfjsonpath.New("data"),
						knownvalue.StringRegexp(regexp.MustCompile(`.*SecMarker BEGIN_ALLOWLIST_login.*`)),
					),
				},
			},
		},
	})
}

func TestOwaspCustomRuleResourceRealData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testRealOwaspCustomRuleWithData1(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_rule.rule",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("librechat_acl.conf"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.data",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("gpustack_acl.txt"),
					),
				},
			},
		},
	})
}

func testOwaspCustomRuleResource() string {
	return `
resource "loadmaster_owasp_custom_rule" "test_rule" {
  filename = "test_rule_replace_url.conf"
  data = <<EOT
SecMarker BEGIN_ALLOWLIST_login
 
# START allowlisting block for URI /login SecRule REQUEST_URI "!@beginsWith /login" \
    "id:11001,phase:1,pass,t:lowercase,nolog,skipAfter:END_ALLOWLIST_login"
SecRule REQUEST_URI "!@beginsWith /login" \
    "id:11002,phase:2,pass,t:lowercase,nolog,skipAfter:END_ALLOWLIST_login"
 
# Validate HTTP method
SecRule REQUEST_METHOD "!@pm GET HEAD POST OPTIONS" \
    "id:11100,phase:1,deny,status:405,log,tag:'Login Allowlist',\
    msg:'Method %%{MATCHED_VAR} not allowed'"
 
# Validate URIs
SecRule REQUEST_FILENAME "@beginsWith /login/static/css" \
    "id:11200,phase:1,pass,nolog,tag:'Login Allowlist',\
    skipAfter:END_ALLOWLIST_URIBLOCK_login"
SecRule REQUEST_FILENAME "@beginsWith /login/static/img" \
    "id:11201,phase:1,pass,nolog,tag:'Login Allowlist',\
    skipAfter:END_ALLOWLIST_URIBLOCK_login"
SecRule REQUEST_FILENAME "@beginsWith /login/static/js" \
    "id:11202,phase:1,pass,nolog,tag:'Login Allowlist',\
    skipAfter:END_ALLOWLIST_URIBLOCK_login"
SecRule REQUEST_FILENAME \
    "@rx ^/login/(displayLogin|login|logout).do$" \
    "id:11250,phase:1,pass,nolog,tag:'Login Allowlist',\
    skipAfter:END_ALLOWLIST_URIBLOCK_login"
 
# If we land here, we are facing an unknown URI, # which is why we will respond using the 404 status code SecAction "id:11299,phase:1,deny,status:404,log,tag:'Login Allowlist',\
    msg:'Unknown URI %%{REQUEST_URI}'"
 
SecMarker END_ALLOWLIST_URIBLOCK_login
 
# Validate parameter names
SecRule ARGS_NAMES "!@rx ^(username|password|sectoken)$" \
    "id:11300,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'Unknown parameter: %%{MATCHED_VAR_NAME}'"
 
# Validate each parameter's uniqueness
SecRule &ARGS:username  "@gt 1" \
    "id:11400,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'%%{MATCHED_VAR_NAME} occurring more than once'"
SecRule &ARGS:password  "@gt 1" \
    "id:11401,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'%%{MATCHED_VAR_NAME} occurring more than once'"
SecRule &ARGS:sectoken  "@gt 1" \
    "id:11402,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'%%{MATCHED_VAR_NAME} occurring more than once'"
 
# Check individual parameters
SecRule ARGS:username "!@rx ^[a-zA-Z0-9.@_-]{1,64}$" \
    "id:11500,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'Invalid parameter format: %%{MATCHED_VAR_NAME} (%%{MATCHED_VAR})'"
SecRule ARGS:sectoken "!@rx ^[a-zA-Z0-9]{32}$" \
    "id:11501,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'Invalid parameter format: %%{MATCHED_VAR_NAME} (%%{MATCHED_VAR})'"
SecRule ARGS:password "@gt 64" \
    "id:11502,phase:2,deny,log,t:length,tag:'Login Allowlist',\
    msg:'Invalid parameter format: %%{MATCHED_VAR_NAME} too long (%%{MATCHED_VAR} bytes)'"
SecRule ARGS:password "@validateByteRange 33-244" \
    "id:11503,phase:2,deny,log,tag:'Login Allowlist',\
    msg:'Invalid parameter format: %%{MATCHED_VAR_NAME} (%%{MATCHED_VAR})'"
 
SecMarker END_ALLOWLIST_login
EOT
}
`
}

func testRealOwaspCustomRuleWithData1() string {
	return `

resource "loadmaster_owasp_custom_data" "data" {
  filename = "gpustack_acl.txt"
  data = <<EOT
# ---------------------------------------------------------------
# List of allowed IP
# ---------------------------------------------------------------

# Uni
130.92.0.0/16
172.16.0.0/12
EOT
}

resource "loadmaster_owasp_custom_rule" "rule" {
  filename = "librechat_acl.conf"
  data = <<EOT
# ---------------------------------------------------------------
# ACL Check
# ---------------------------------------------------------------

SecMarker BEGIN_ALLOWLIST_URI
# gpustack
SecRule REQUEST_HEADERS:Host "@rx ^gpustack\.unibe\.ch$" "chain,id:100,phase:1,deny,t:lowercase,log,msg:'Access for IP %%{REMOTE_ADDR} is not allowed'" 
SecRule REMOTE_ADDR "!@ipMatchFromFile gpustack_acl.txt" \

SecMarker END_ALLOWLIST_URI
EOT

  depends_on = [loadmaster_owasp_custom_data.data]
}

`
}
