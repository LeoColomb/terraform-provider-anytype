data "anytype_objects" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "object_names" {
  value = [for o in data.anytype_objects.all.objects : o.name]
}
