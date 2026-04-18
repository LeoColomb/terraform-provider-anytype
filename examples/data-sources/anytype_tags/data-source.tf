data "anytype_tags" "status" {
  space_id    = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  property_id = "bafyreids36kpw5ppuwm3ce2p4ezb3ab7cihhkq6yfbwzwpp4mln7rcgw7a"
}

output "tag_names" {
  value = [for t in data.anytype_tags.status.tags : t.name]
}
