terraform {
  required_providers {
    anytype = {
      source  = "LeoColomb/anytype"
      version = "~> 0.1"
    }
  }
}

# The endpoint defaults to the local Anytype desktop app
# (http://127.0.0.1:31009). The api_key is read from the ANYTYPE_API_KEY
# environment variable when omitted from the config.
provider "anytype" {
  endpoint = "http://127.0.0.1:31009"
  api_key  = var.anytype_api_key
}

variable "anytype_api_key" {
  type      = string
  sensitive = true
}
