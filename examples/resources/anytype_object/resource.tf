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
# Referencing the type by `type_id` lets the provider resolve `type_key`
# automatically, so the object never has to reach into `anytype_type.note.key`.
resource "anytype_object" "welcome" {
  space_id = anytype_space.wiki.id
  type_id  = anytype_type.note.id
  name     = "Welcome"
  body     = <<-EOT
    # Welcome to the wiki
    This wiki is managed by Terraform.
  EOT
}
