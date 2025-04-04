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

func TestSubVirtualServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testSubVirtualServiceDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.loadmaster_sub_virtual_service.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.loadmaster_sub_virtual_service.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("http2"),
					),
				},
			},
		},
	})
}

const testSubVirtualServiceDataSourceConfig = `
resource "loadmaster_virtual_service" "example" {
	address = "10.0.0.4"
	port = "9090"
	protocol = "tcp"
}

resource "loadmaster_sub_virtual_service" "example" {
  virtual_service_id = loadmaster_virtual_service.example.id

  nickname = "subvs"
  type = "http2"
}
  
data "loadmaster_sub_virtual_service" "test" {
  id = loadmaster_sub_virtual_service.example.id
}
`
