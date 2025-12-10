// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
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

func testOwaspCustomDataResource() string {
	return `
resource "loadmaster_owasp_custom_data" "test_data" {
  filename = "test_rule_replace_url.txt"
  data = "Data"
}
`
}
