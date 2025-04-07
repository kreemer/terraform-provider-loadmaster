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

func TestReplaceBodyRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testReplaceBodyRuleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_replace_body_rule.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace_body"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_replace_body_rule.test",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_replace_body_rule.test",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
		},
	})
}

const testReplaceBodyRuleDataSourceConfig = `
resource "loadmaster_replace_body_rule" "test_rule" {
  id = "test_rule_replace_body"
  pattern = "A"
  replacement = "B"
}

data "loadmaster_replace_body_rule" "test" {
  id = loadmaster_replace_body_rule.test_rule.id
}
`
