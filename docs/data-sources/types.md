# anytype_types (Data Source)

List all Anytype types defined in a given space.

## Example Usage

```terraform
data "anytype_types" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "type_names" {
  value = [for t in data.anytype_types.all.types : t.name]
}
```

## Schema

### Required

- `space_id` (String) — The ID of the space.

### Read-Only

- `types` (List of Object) — The list of types in the space. Each item has
  `id`, `key`, `name`, `plural_name`, `layout`, `object`, and `archived`.
