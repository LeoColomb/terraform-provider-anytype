data "anytype_templates" "page_templates" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  type_id  = "bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu"
}

output "template_names" {
  value = [for t in data.anytype_templates.page_templates.templates : t.name]
}
