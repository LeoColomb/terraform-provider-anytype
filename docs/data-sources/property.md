# anytype_property (Data Source)

Look up a single Anytype property by ID in a given space.

## Example Usage

```terraform
data "anytype_property" "status" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  id       = "bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a"
}
```

## Schema

### Required

- `id` (String) — The ID of the property.
- `space_id` (String) — The ID of the space.

### Read-Only

- `key`, `name`, `format`, `object` — see the `anytype_property` resource.
