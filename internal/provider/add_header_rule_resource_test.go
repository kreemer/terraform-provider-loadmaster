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

func TestAddHeaderRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAddHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("test"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_add_header_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAddHeaderRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("test"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
				},
			},
			{
				Config: testAddHeaderRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("test2"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_add_header_rule.test_rule",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
				},
			},
		},
	})
}

func testAddHeaderRuleResource() string {
	return `
resource "loadmaster_add_header_rule" "test_rule" {
  id = "test_rule"
  replacement = "test"
  header = "TEST"
}
`
}

func testAddHeaderRuleResourceUpdate() string {
	return `
resource "loadmaster_add_header_rule" "test_rule" {
  id = "test_rule"
  replacement = "test2"
  match_type = "postfix"
}
`
}
