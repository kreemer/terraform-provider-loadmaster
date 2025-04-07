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

func TestModifyUrlRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testModifyUrlRuleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_modify_url_rule.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_replace"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_modify_url_rule.test",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("A"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_modify_url_rule.test",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("B"),
					),
				},
			},
		},
	})
}

const testModifyUrlRuleDataSourceConfig = `
resource "loadmaster_modify_url_rule" "test_rule" {
  id = "test_rule_replace"
  pattern = "A"
  replacement = "B"
}

data "loadmaster_modify_url_rule" "test" {
  id = loadmaster_modify_url_rule.test_rule.id
}
`
