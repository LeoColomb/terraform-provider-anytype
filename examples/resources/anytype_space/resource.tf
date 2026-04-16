resource "anytype_space" "wiki" {
  name        = "Engineering Wiki"
  description = "The local-first engineering wiki"
}

output "wiki_id" {
  value = anytype_space.wiki.id
}
