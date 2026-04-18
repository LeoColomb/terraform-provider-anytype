data "anytype_spaces" "all" {}

output "space_names" {
  value = [for s in data.anytype_spaces.all.spaces : s.name]
}
