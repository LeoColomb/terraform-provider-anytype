data "anytype_types" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "type_names" {
  value = [for t in data.anytype_types.all.types : t.name]
}
