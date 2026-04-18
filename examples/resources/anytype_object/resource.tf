resource "anytype_space" "wiki" {
  name = "Engineering Wiki"
}

resource "anytype_type" "note" {
  space_id    = anytype_space.wiki.id
  name        = "Note"
  plural_name = "Notes"
  layout      = "note"
}

# A concrete object of the "note" type, created from a markdown body.
resource "anytype_object" "welcome" {
  space_id = anytype_space.wiki.id
  type_key = anytype_type.note.key
  name     = "Welcome"
  body     = <<-EOT
    # Welcome to the wiki
    This wiki is managed by Terraform.
  EOT
}
