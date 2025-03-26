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

func TestSubVirtualServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testSubVirtualServiceResourceConfig("hello"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_sub_virtual_service.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_sub_virtual_service.test",
						tfjsonpath.New("nickname"),
						knownvalue.StringExact("hello"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_sub_virtual_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testSubVirtualServiceResourceConfig("blupp"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_sub_virtual_service.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_sub_virtual_service.test",
						tfjsonpath.New("nickname"),
						knownvalue.StringExact("blupp"),
					),
				},
			},
		},
	})
}

func testSubVirtualServiceResourceConfig(nickname string) string {
	return fmt.Sprintf(`
resource "loadmaster_virtual_service" "test" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"
}

resource "loadmaster_sub_virtual_service" "test" {
  virtual_service_id = loadmaster_virtual_service.test.id
  nickname = "%s"
}
`, nickname)
}
