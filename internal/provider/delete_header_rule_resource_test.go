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

func TestDeleteHeaderRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testDeleteHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_delete"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST1"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_delete_header_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testDeleteHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_delete"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST1"),
					),
				},
			},
			{
				Config: testDeleteHeaderRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_delete"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_delete_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST2"),
					),
				},
			},
		},
	})
}

func testDeleteHeaderRuleResource() string {
	return `
resource "loadmaster_delete_header_rule" "test_rule" {
  id = "test_rule_delete"
  replacement = "test"
  header = "TEST1"
}
`
}

func testDeleteHeaderRuleResourceUpdate() string {
	return `
resource "loadmaster_delete_header_rule" "test_rule" {
  id = "test_rule_delete"
  header = "TEST2"
}
`
}
