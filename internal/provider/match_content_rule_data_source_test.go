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

func TestMatchContentRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testMatchContentRuleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_match_content_rule.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("test_rule"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_match_content_rule.test",
						tfjsonpath.New("pattern"),
						knownvalue.StringExact("test"),
					),
				},
			},
		},
	})
}

const testMatchContentRuleDataSourceConfig = `
resource "loadmaster_match_content_rule" "test_rule" {
  id = "test_rule"
  pattern = "test"
  match_type = "postfix"
}

data "loadmaster_match_content_rule" "test" {
  id = loadmaster_match_content_rule.test_rule.id
}
`
