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

func TestVirtualServiceRestartAction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testVirtualServiceRestartActionResource(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test1",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.4"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test1",
						tfjsonpath.New("port"),
						knownvalue.StringExact("9090"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test1",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("tcp"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test1",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test1",
						tfjsonpath.New("enabled"),
						knownvalue.Bool(true),
					),
				},
			},
			{
				Config: testVirtualServiceRestartActionResourceDisabled(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test2",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.0.0.4"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test2",
						tfjsonpath.New("port"),
						knownvalue.StringExact("9091"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test2",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("tcp"),
					),
					statecheck.ExpectKnownValue(
						"loadmaster_virtual_service.test2",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
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

func testVirtualServiceRestartActionResource() string {
	return `
resource "loadmaster_virtual_service" "test1" {
  address = "10.0.0.4"
  port = "9090"
  protocol = "tcp"

  enabled = true

  lifecycle {
    action_trigger {
      events = [after_create, after_update]
      actions = [action.loadmaster_virtual_service_restart.this]
    }
  }
}

action "loadmaster_virtual_service_restart" "this" {
  config {
	virtual_service_id = loadmaster_virtual_service.test1.id
  }
}
`
}

func testVirtualServiceRestartActionResourceDisabled() string {
	return `
resource "loadmaster_virtual_service" "test2" {
  address = "10.0.0.4"
  port = "9091"
  protocol = "tcp"

  enabled = false

  lifecycle {
    action_trigger {
      events = [after_create, after_update]
      actions = [action.loadmaster_virtual_service_restart.this]
    }
  }
}

action "loadmaster_virtual_service_restart" "this" {
  config {
	virtual_service_id = loadmaster_virtual_service.test2.id
  }
}
`
}
