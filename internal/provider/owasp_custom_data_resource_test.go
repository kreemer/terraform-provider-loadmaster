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

func TestOwaspCustomDataResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testOwaspCustomDataResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("test_rule_replace_url.txt"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("data"),
						knownvalue.StringExact("Data"),
					),
				},
			},
			{
				ResourceName:                         "loadmaster_owasp_custom_data.test_data",
				ImportStateVerifyIdentifierAttribute: "filename",
				ImportStateId:                        "test_rule_replace_url.txt",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				Config: testOwaspCustomDataResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("test_rule_replace_url.txt"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("data"),
						knownvalue.StringExact("Data"),
					),
				},
			},
		},
	})
}

func TestOwaspCustomDataResourceReal1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testOwaspCustomDataResourceReal1(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("real_rule_1.txt"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_owasp_custom_data.test_data",
						tfjsonpath.New("data"),
						knownvalue.StringRegexp(regexp.MustCompile(`130\.92\.0\.0/16`)),
					),
				},
			},
		},
	})
}

func testOwaspCustomDataResource() string {
	return `
resource "loadmaster_owasp_custom_data" "test_data" {
  filename = "test_rule_replace_url.txt"
  data = "Data"
}
`
}

func testOwaspCustomDataResourceReal1() string {
	return `
resource "loadmaster_owasp_custom_data" "test_data" {
  filename = "real_rule_1.txt"
  data = <<EOT
# ---------------------------------------------------------------
# List of allowed IP
# ---------------------------------------------------------------

# Universität Bern
130.92.0.0/16
172.16.0.0/12
EOT
}
`
}
