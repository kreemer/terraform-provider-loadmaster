// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
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
				Config: testAccExampleResourceConfig("hello"),
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
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("nickname"),
						knownvalue.StringExact("hello"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				ResourceName:      "loadmaster_virtual_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccExampleResourceConfig("blupp"),
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
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("nickname"),
						knownvalue.StringExact("blupp"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: testAccExampleResourceConfigDisabled(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test2",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

func testAccExampleResourceConfig(nickname string) string {
	return fmt.Sprintf(`
resource "loadmaster_virtual_service" "test" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"

  nickname = "%s"
  enabled = true
}
`, nickname)
}

func testAccExampleResourceConfigDisabled() string {
	return `
resource "loadmaster_virtual_service" "test2" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"

  enabled = false
}
`
}
