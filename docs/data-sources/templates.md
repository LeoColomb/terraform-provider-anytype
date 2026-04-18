# anytype_templates (Data Source)

List all templates defined on a given Anytype type.

## Example Usage

```terraform
data "anytype_templates" "page_templates" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  type_id  = "bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu"
}
```

## Schema

### Required

- `space_id` (String), `type_id` (String).

### Read-Only

- `templates` (List of Object) — Each item has `id`, `name`, `layout`,
  `snippet`, `object`, `archived`.
