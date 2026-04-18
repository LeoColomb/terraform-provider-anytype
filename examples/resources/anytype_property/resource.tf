resource "anytype_space" "wiki" {
  name = "Engineering Wiki"
}

# A simple free-form text property.
resource "anytype_property" "summary" {
  space_id = anytype_space.wiki.id
  name     = "Summary"
  format   = "text"
}

# A select property seeded with a fixed set of tags.
resource "anytype_property" "priority" {
  space_id = anytype_space.wiki.id
  name     = "Priority"
  format   = "select"

  tags = [
    { name = "Low", color = "ice" },
    { name = "Medium", color = "yellow" },
    { name = "High", color = "red" },
  ]
}

# A date property.
resource "anytype_property" "due_date" {
  space_id = anytype_space.wiki.id
  name     = "Due date"
  format   = "date"
}
