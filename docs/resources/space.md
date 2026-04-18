# anytype_space (Resource)

Manages an [Anytype space](https://anytype.io).

~> **Note** The Anytype API does not currently support deleting spaces. Running
`terraform destroy` removes the space from state but leaves the space intact
in your Anytype account.

## Example Usage

```terraform
resource "anytype_space" "wiki" {
  name        = "Engineering Wiki"
  description = "The local-first engineering wiki"
}
```

## Schema

### Required

- `name` (String) — The name of the space.

### Optional

- `description` (String) — The description of the space.

### Read-Only

- `id` (String) — The ID of the space.
- `network_id` (String) — The Anytype network the space belongs to.
- `gateway_url` (String) — Gateway URL used to serve files and media.
- `object` (String) — Data model of the object (`space` or `chat`).

## Import

Import using the space ID:

```shell
terraform import anytype_space.wiki bafyreigyfkt6rbv24sbv5aq2hko3bhmv5xxlf22b4bypdu6j7hnphm3psq.23me69r569oi1
```
