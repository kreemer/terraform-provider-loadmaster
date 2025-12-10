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

func TestOwaspCustomDataDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testOwaspCustomDataDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_owasp_custom_data.test",
						tfjsonpath.New("filename"),
						knownvalue.StringExact("test_rule_replace.txt"),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_owasp_custom_data.test",
						tfjsonpath.New("data"),
						knownvalue.StringExact("Data"),
					),
				},
			},
		},
	})
}

const testOwaspCustomDataDataSourceConfig = `
resource "loadmaster_owasp_custom_data" "test_data" {
  filename = "test_rule_replace.txt"
  data = "Data"
}

data "loadmaster_owasp_custom_data" "test" {
  filename = loadmaster_owasp_custom_data.test_data.filename
}
`
