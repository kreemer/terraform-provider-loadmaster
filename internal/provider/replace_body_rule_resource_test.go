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

func TestReplaceBodyRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testReplaceBodyRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_body"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_replace_body_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testReplaceBodyRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_body"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				Config: testReplaceBodyRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_body"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("C"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_replace_body_rule.test_rule",
						tfjsonpath.New("no_case"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func testReplaceBodyRuleResource() string {
	return `
resource "loadmaster_replace_body_rule" "test_rule" {
  id = "test_rule_replace_body"
  pattern = "A"
  replacement = "B"
}
`
}

func testReplaceBodyRuleResourceUpdate() string {
	return `
resource "loadmaster_replace_body_rule" "test_rule" {
  id = "test_rule_replace_body"
  pattern = "A"
  replacement = "C"
  no_case = true
}
`
}
