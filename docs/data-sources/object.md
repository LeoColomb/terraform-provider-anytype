# anytype_object (Data Source)

Look up a single Anytype object by ID in a given space.

## Example Usage

```terraform
data "anytype_object" "welcome" {
  space_id = "bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1"
  id       = "bafyreie6n5l5nkbjal37su54cha4coy7qzuhrnajluzv5qd5jvtsrxkequ"
}
```

## Schema

### Required

- `space_id` (String), `id` (String).

### Read-Only

- `name`, `markdown`, `snippet`, `layout`, `object`, `archived` — see the
  `anytype_object` resource.
