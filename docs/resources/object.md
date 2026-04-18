# anytype_object (Resource)

Manages an [Anytype object](https://anytype.io) inside a space. Objects are
concrete instances of an `anytype_type` — this resource manages their name,
markdown body, and type-key mapping. The polymorphic `properties` array (a
`oneOf` in the OpenAPI) is not yet exposed by the provider.

## Example Usage

```terraform
resource "anytype_type" "note" {
  space_id    = anytype_space.wiki.id
  name        = "Note"
  plural_name = "Notes"
  layout      = "note"
}

resource "anytype_object" "welcome" {
  space_id = anytype_space.wiki.id
  type_key = anytype_type.note.key
  name     = "Welcome"
  body     = "# Welcome\nManaged by Terraform."
}
```

## Schema

### Required

- `space_id` (String) — The ID of the space. Changing this forces a new
  resource.
- `type_key` (String) — The key of the type of the object. Changing this
  forces a new resource.

### Optional

- `name` (String) — The name of the object.
- `body` (String) — The markdown body of the object. Passed as `body` on
  create and as `markdown` on update.
- `template_id` (String) — Optional template ID to seed the object body on
  create. Changing this forces a new resource.

### Read-Only

- `id` (String) — The ID of the object.
- `layout` (String) — The layout of the object, inherited from the type.
- `markdown` (String) — The current markdown body as returned by the API.
- `snippet` (String) — A short snippet of the object body.
- `object` (String) — The data model of the object (`object`).
- `archived` (Bool) — Whether the object is archived.

## Import

Import using the composite ID `<space_id>/<object_id>`:

```shell
terraform import anytype_object.welcome bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1/bafyreie6n5l5nkbjal37su54cha4coy7qzuhrnajluzv5qd5jvtsrxkequ
```
