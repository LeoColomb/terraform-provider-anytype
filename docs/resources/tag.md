# anytype_tag (Resource)

Manages a tag on a `select` or `multi_select` Anytype property. Tags carry the
set of allowed values for the property.

## Example Usage

```terraform
resource "anytype_property" "status" {
  space_id = anytype_space.workflow.id
  name     = "Status"
  format   = "select"
}

resource "anytype_tag" "done" {
  space_id    = anytype_space.workflow.id
  property_id = anytype_property.status.id
  name        = "Done"
  color       = "lime"
}
```

## Schema

### Required

- `space_id` (String) — The ID of the space. Changing this forces a new
  resource.
- `property_id` (String) — The ID of the property. Changing this forces a new
  resource.
- `name` (String) — The name of the tag.
- `color` (String) — One of `grey`, `yellow`, `orange`, `red`, `pink`,
  `purple`, `blue`, `ice`, `teal`, `lime`.

### Optional

- `key` (String) — Optional custom key. If omitted, Anytype generates one.

### Read-Only

- `id` (String) — The ID of the tag.
- `object` (String) — The data model of the object (`tag`).

## Import

Import using the composite ID `<space_id>/<property_id>/<tag_id>`:

```shell
terraform import anytype_tag.done bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1/bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a/bafyreiaixlnaefu3ci22zdenjhsdlyaeeoyjrsid5qhfeejzlccijbj7sq
```
