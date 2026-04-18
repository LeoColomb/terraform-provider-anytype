data "anytype_members" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "owners" {
  value = [for m in data.anytype_members.all.members : m.name if m.role == "owner"]
}
