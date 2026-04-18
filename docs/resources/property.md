# anytype_property (Resource)

Manages an [Anytype property](https://anytype.io) inside a space. Properties
are typed fields that can be attached to `anytype_type` definitions.

## Example Usage

```terraform
resource "anytype_space" "wiki" {
  name = "Engineering Wiki"
}

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
```

## Schema

### Required

- `space_id` (String) — The ID of the space. Changing this forces a new
  resource.
- `name` (String) — The human-readable name of the property.
- `format` (String) — One of `text`, `number`, `select`, `multi_select`,
  `date`, `files`, `checkbox`, `url`, `email`, `phone`, `objects`. Immutable
  once created.

### Optional

- `key` (String) — The snake_case key. If omitted, Anytype generates one.
- `tags` (List of Object) — Initial tags seeded on create for `select` /
  `multi_select` properties. Tags are immutable once created; changing them
  forces the property to be replaced. Each entry has:
  - `name` (String, required)
  - `color` (String, required) — One of `grey`, `yellow`, `orange`, `red`,
    `pink`, `purple`, `blue`, `ice`, `teal`, `lime`.
  - `key` (String, optional) — Custom key for the tag.

### Read-Only

- `id` (String) — The ID of the property.
- `object` (String) — The data model of the object (`property`).

## Import

Import using the composite ID `<space_id>/<property_id>`:

```shell
terraform import anytype_property.priority bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1/bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a
```
