# anytype_template (Data Source)

Look up a single Anytype template (an `ObjectWithBody` attached to a type).

## Example Usage

```terraform
data "anytype_template" "default" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  type_id  = "bafyreigyb6l5szohs32ts26ku2j42yd65e6hqy2u3gtzgdwqv6hzftsetu"
  id       = "bafyreictrp3obmnf6dwejy5o4p7bderaaia4bdg2psxbfzf44yya5uutge"
}
```

## Schema

### Required

- `space_id` (String), `type_id` (String), `id` (String).

### Read-Only

- `name`, `markdown`, `snippet`, `layout`, `object`, `archived`.
