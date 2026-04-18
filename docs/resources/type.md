# anytype_type (Resource)

Manages an [Anytype type](https://anytype.io) inside a space. A type describes
the shape of objects and links to existing `anytype_property` resources by
`id`. The provider resolves the backend-required `key` / `name` / `format`
triplet automatically, so the consuming type never has to re-declare
attributes of the properties it links.

## Example Usage

```terraform
resource "anytype_space" "crm" {
  name = "CRM"
}

resource "anytype_property" "status" {
  space_id = anytype_space.crm.id
  name     = "Status"
  format   = "select"

  tags = [
    { name = "Lead", color = "yellow" },
    { name = "Qualified", color = "blue" },
    { name = "Closed", color = "lime" },
  ]
}

resource "anytype_type" "account" {
  space_id    = anytype_space.crm.id
  name        = "Account"
  plural_name = "Accounts"
  layout      = "basic"

  properties = [
    { id = anytype_property.status.id },
  ]
}
```

## Schema

### Required

- `space_id` (String) — The ID of the space the type belongs to. Changing this
  forces a new resource.
- `name` (String) — The singular name of the type (e.g. `Page`).
- `plural_name` (String) — The plural name of the type (e.g. `Pages`).
- `layout` (String) — One of `basic`, `profile`, `action`, `note`.

### Optional

- `key` (String) — The snake_case key of the type. If omitted, Anytype
  generates one.
- `properties` (List of Object) — Properties linked to this type. Each entry
  references an existing `anytype_property` by `id`:
  - `id` (String, required) — The ID of the `anytype_property` to link.
  - `key` (String, read-only) — Resolved from `id`.
  - `name` (String, read-only) — Resolved from `id`.
  - `format` (String, read-only) — Resolved from `id`.

### Read-Only

- `id` (String) — The ID of the type.
- `object` (String) — The data model of the object (`type`).
- `archived` (Bool) — Whether the type is archived.

## Import

Import using the composite ID `<space_id>/<type_id>`:

```shell
terraform import anytype_type.account bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1/bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu
```
