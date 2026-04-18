# anytype_tag (Data Source)

Look up a single Anytype tag by ID in a given property.

## Example Usage

```terraform
data "anytype_tag" "done" {
  space_id    = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  property_id = "bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a"
  id          = "bafyreiaixlnaefu3ci22zdenjhsdlyaeeoyjrsid5qhfeejzlccijbj7sq"
}
```

## Schema

### Required

- `space_id` (String), `property_id` (String), `id` (String).

### Read-Only

- `key`, `name`, `color`, `object` — see the `anytype_tag` resource.
