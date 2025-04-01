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

func TestMatchContentRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testMatchContentRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("match_type"),
						knownvalue.StringExact("Regex"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_match_content_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testMatchContentRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
				},
			},
			{
				Config: testMatchContentRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact(""),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_match_content_rule.test_rule",
						tfjsonpath.New("match_type"),
						knownvalue.StringExact("postfix"),
					),
				},
			},
		},
	})
}

func testMatchContentRuleResource() string {
	return `
resource "loadmaster_match_content_rule" "test_rule" {
  id = "test_rule"
  pattern = "test"
  header = "TEST"
}
`
}

func testMatchContentRuleResourceUpdate() string {
	return `
resource "loadmaster_match_content_rule" "test_rule" {
  id = "test_rule"
  pattern = "test"
  match_type = "postfix"
}
`
}
