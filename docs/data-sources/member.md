# anytype_member (Data Source)

Look up a single member of an Anytype space by ID or identity.

## Example Usage

```terraform
data "anytype_member" "me" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  id       = "_participant_bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1_AAjEaEwPF4nkEh9AWkqEnzcQ8HziBB4ETjiTpvRCQvWnSMDZ"
}
```

## Schema

### Required

- `space_id` (String), `id` (String).

### Read-Only

- `identity`, `global_name`, `name`, `role`, `status`, `object`.
