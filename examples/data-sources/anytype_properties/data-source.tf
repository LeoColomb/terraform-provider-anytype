data "anytype_properties" "all" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "property_names" {
  value = [for p in data.anytype_properties.all.properties : p.name]
}
