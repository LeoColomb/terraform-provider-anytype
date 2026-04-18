# Manage an Anytype space and populate it with a custom type whose shape is
# described by a handful of properties. The property resources below are
# linked to the `anytype_type` by `id` alone — the provider resolves the
# backend-required key / name / format triplet automatically, so the consuming
# `anytype_type` never has to re-assign attributes of `anytype_property`.

resource "anytype_space" "crm" {
  name        = "CRM"
  description = "Managed by Terraform"
}

resource "anytype_property" "status" {
  space_id = anytype_space.crm.id
  name     = "Status"
  format   = "select"

  tags = [
    { name = "Lead", color = "yellow" },
    { name = "Qualified", color = "blue" },
    { name = "Closed", color = "lime" },
  ]
}

resource "anytype_property" "owner" {
  space_id = anytype_space.crm.id
  name     = "Owner"
  format   = "text"
}

resource "anytype_property" "next_followup" {
  space_id = anytype_space.crm.id
  name     = "Next follow-up"
  format   = "date"
}

resource "anytype_type" "account" {
  space_id    = anytype_space.crm.id
  name        = "Account"
  plural_name = "Accounts"
  layout      = "basic"

  properties = [
    { id = anytype_property.status.id },
    { id = anytype_property.owner.id },
    { id = anytype_property.next_followup.id },
  ]
}

output "account_type_id" {
  value = anytype_type.account.id
}
