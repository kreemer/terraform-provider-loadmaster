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

func TestAddHeaderRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAddHeaderRuleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_add_header_rule.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_add_header_rule.test",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_add_header_rule.test",
						tfjsonpath.New("replacement"),
						knownvalue.StringExact("test"),
					),
				},
			},
		},
	})
}

const testAddHeaderRuleDataSourceConfig = `
resource "loadmaster_add_header_rule" "test_rule" {
  id = "test_rule"
  header = "TEST"
  replacement = "postfix"
}

data "loadmaster_add_header_rule" "test" {
  id = loadmaster_add_header_rule.test_rule.id
}
`
