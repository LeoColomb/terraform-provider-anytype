# anytype_type (Data Source)

Look up a single Anytype type by ID in a given space.

## Example Usage

```terraform
data "anytype_type" "page" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  id       = "bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu"
}

output "page_property_keys" {
  value = [for p in data.anytype_type.page.properties : p.key]
}
```

## Schema

### Required

- `id` (String) — The ID of the type.
- `space_id` (String) — The ID of the space.

### Read-Only

- `key`, `name`, `plural_name`, `layout`, `object`, `archived` — see the
  `anytype_type` resource.
- `properties` (List of Object) — The properties currently linked to this
  type. Each object has `id`, `key`, `name`, and `format`.
