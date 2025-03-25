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

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.4"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("port"),
						knownvalue.StringExact("9090"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("tcp"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ResourceName:      "loadmaster_virtual_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccExampleResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.4"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("port"),
						knownvalue.StringExact("9090"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("tcp"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig() string {
	return `
resource "loadmaster_virtual_service" "test" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"
}
`
}
