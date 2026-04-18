# anytype_space (Data Source)

Look up a single Anytype space by ID.

## Example Usage

```terraform
data "anytype_space" "wiki" {
  id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}
```

## Schema

### Required

- `id` (String) — The ID of the space.

### Read-Only

- `name`, `description`, `network_id`, `gateway_url`, `object` — see the
  `anytype_space` resource for semantics.
