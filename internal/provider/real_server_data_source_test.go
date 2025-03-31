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

func TestRealServerDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testRealServerDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_real_server.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_real_server.test",
						tfjsonpath.New("virtual_service_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const testRealServerDataSourceConfig = `
resource "loadmaster_virtual_service" "test" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"
}

resource "loadmaster_real_server" "test" {
  virtual_service_id = loadmaster_virtual_service.test.id
  address = "10.0.0.99"
  port = "80"
  protocol = "tcp"
}

data "loadmaster_real_server" "test" {
  id = loadmaster_real_server.test.id
  virtual_service_id = loadmaster_virtual_service.test.id
}
`
