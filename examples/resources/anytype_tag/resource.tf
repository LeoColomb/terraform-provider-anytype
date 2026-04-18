resource "anytype_space" "workflow" {
  name = "Workflow"
}

resource "anytype_property" "status" {
  space_id = anytype_space.workflow.id
  name     = "Status"
  format   = "select"
}

resource "anytype_tag" "todo" {
  space_id    = anytype_space.workflow.id
  property_id = anytype_property.status.id
  name        = "Todo"
  color       = "grey"
}

resource "anytype_tag" "doing" {
  space_id    = anytype_space.workflow.id
  property_id = anytype_property.status.id
  name        = "Doing"
  color       = "yellow"
}

resource "anytype_tag" "done" {
  space_id    = anytype_space.workflow.id
  property_id = anytype_property.status.id
  name        = "Done"
  color       = "lime"
}
