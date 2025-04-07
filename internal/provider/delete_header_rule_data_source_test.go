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

func TestDeleteHeaderRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testDeleteHeaderRuleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_delete_header_rule.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule_delete"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_delete_header_rule.test",
						tfjsonpath.New("header"),
						knownvalue.StringExact("TEST"),
					),
				},
			},
		},
	})
}

const testDeleteHeaderRuleDataSourceConfig = `
resource "loadmaster_delete_header_rule" "test_rule" {
  id = "test_rule_delete"
  header = "TEST"
}

data "loadmaster_delete_header_rule" "test" {
  id = loadmaster_delete_header_rule.test_rule.id
}
`
