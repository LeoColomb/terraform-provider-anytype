data "anytype_space" "wiki" {
  id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
}

output "wiki_name" {
  value = data.anytype_space.wiki.name
}
