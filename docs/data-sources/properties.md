# anytype_properties (Data Source)

List all Anytype properties defined in a given space.

## Example Usage

```terraform
data "anytype_properties" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "property_names" {
  value = [for p in data.anytype_properties.all.properties : p.name]
}
```

## Schema

### Required

- `space_id` (String) — The ID of the space.

### Read-Only

- `properties` (List of Object) — The list of properties in the space. Each
  item has `id`, `key`, `name`, `format`, and `object`.
