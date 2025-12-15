resource "loadmaster_virtual_service" "this" {
  address  = "10.0.0.1"
  port     = 8080
  protocol = "tcp"

  lifecycle {
    action_trigger {
      events  = [after_create, after_update]
      actions = [action.loadmaster_virtual_service_restart.this]
    }
  }
}

action "loadmaster_virtual_service_restart" "this" {
  config {
    virtual_service_id = loadmaster_virtual_service.this.id
  }
}
