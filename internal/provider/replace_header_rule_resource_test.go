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

func TestReplaceHeaderRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testReplaceHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST1"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_replace_header_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testReplaceHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST1"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				Config: testReplaceHeaderRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST2"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("C"),
					),
				},
			},
		},
	})
}

func testReplaceHeaderRuleResource() string {
	return `
resource "loadmaster_replace_header_rule" "test_rule" {
  id = "test_rule_replace"
  header = "TEST1"
  pattern = "A"
  replacement = "B"
}
`
}

func testReplaceHeaderRuleResourceUpdate() string {
	return `
resource "loadmaster_replace_header_rule" "test_rule" {
  id = "test_rule_replace"
  header = "TEST2"
  pattern = "A"
  replacement = "C"
}
`
}
