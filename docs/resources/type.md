# anytype_type (Resource)

Manages an [Anytype type](https://anytype.io) inside a space. A type describes
the shape of objects and can declare a set of linked properties by
`key` / `name` / `format`.

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
    {
      key    = anytype_property.status.key
      name   = anytype_property.status.name
      format = anytype_property.status.format
    },
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
- `properties` (List of Object) — Properties linked to this type. Each object
  has:
  - `key` (String, required)
  - `name` (String, required)
  - `format` (String, required) — One of `text`, `number`, `select`,
    `multi_select`, `date`, `files`, `checkbox`, `url`, `email`, `phone`,
    `objects`.

### Read-Only

- `id` (String) — The ID of the type.
- `object` (String) — The data model of the object (`type`).
- `archived` (Bool) — Whether the type is archived.

## Import

Import using the composite ID `<space_id>/<type_id>`:

```shell
terraform import anytype_type.account bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1/bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu
```
