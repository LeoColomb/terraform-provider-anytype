# anytype_objects (Data Source)

List all objects in an Anytype space.

## Example Usage

```terraform
data "anytype_objects" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}
```

## Schema

### Required

- `space_id` (String).

### Read-Only

- `objects` (List of Object) — Each item has `id`, `name`, `layout`,
  `snippet`, `object`, `archived`.
