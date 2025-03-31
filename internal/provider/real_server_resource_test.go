// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestRealServerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testRealServerResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_real_server.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_real_server.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.99"),
					),
				},
			},
			{
				ResourceName:      "loadmaster_real_server.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: generateRealServerImportId,
			},
			{
				Config: testRealServerResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_real_server.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_real_server.test",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.99"),
					),
				},
			},
		},
	})
}

func testRealServerResourceConfig() string {
	return `
resource "loadmaster_virtual_service" "test" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"
}

resource "loadmaster_real_server" "test" {
  virtual_service_id = loadmaster_virtual_service.test.id
  address = "10.0.0.99"
  port = "80"
}
`
}

func generateRealServerImportId(state *terraform.State) (string, error) {
	resourceName := "loadmaster_real_server.test"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	fmt.Printf("raw state %s", rawState)
	return fmt.Sprintf("%s/%s", rawState["virtual_service_id"], rawState["id"]), nil
}
