# anytype_tags (Data Source)

List all tags defined on an Anytype property.

## Example Usage

```terraform
data "anytype_tags" "status" {
  space_id    = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  property_id = "bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a"
}
```

## Schema

### Required

- `space_id` (String), `property_id` (String).

### Read-Only

- `tags` (List of Object) — Each item has `id`, `key`, `name`, `color`,
  `object`.
