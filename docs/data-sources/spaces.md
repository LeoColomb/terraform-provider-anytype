# anytype_spaces (Data Source)

List all Anytype spaces accessible by the authenticated user.

## Example Usage

```terraform
data "anytype_spaces" "all" {}

output "space_names" {
  value = [for s in data.anytype_spaces.all.spaces : s.name]
}
```

## Schema

### Read-Only

- `spaces` (List of Object) — The list of accessible spaces. Each item has
  `id`, `name`, `description`, `network_id`, `gateway_url`, and `object`.
