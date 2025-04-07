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

func TestModifyUrlRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testModifyUrlRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_url"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_modify_url_rule.test_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testModifyUrlRuleResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_url"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
			{
				Config: testModifyUrlRuleResourceUpdate(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_url"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_modify_url_rule.test_rule",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("C"),
					),
				},
			},
		},
	})
}

func testModifyUrlRuleResource() string {
	return `
resource "loadmaster_modify_url_rule" "test_rule" {
  id = "test_rule_replace_url"
  pattern = "A"
  replacement = "B"
}
`
}

func testModifyUrlRuleResourceUpdate() string {
	return `
resource "loadmaster_modify_url_rule" "test_rule" {
  id = "test_rule_replace_url"
  pattern = "A"
  replacement = "C"
}
`
}
