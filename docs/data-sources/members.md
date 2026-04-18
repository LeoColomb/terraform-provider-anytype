# anytype_members (Data Source)

List all members of an Anytype space.

## Example Usage

```terraform
data "anytype_members" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "owners" {
  value = [for m in data.anytype_members.all.members : m.name if m.role == "owner"]
}
```

## Schema

### Required

- `space_id` (String).

### Read-Only

- `members` (List of Object) — Each item has `id`, `identity`, `global_name`,
  `name`, `role`, `status`, `object`.
