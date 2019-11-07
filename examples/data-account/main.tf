data "cloudability_account" "cloudability" {
  account_id = "${var.account_id}"
}

variable "account_id" { type = string }

output "id" { value = data.cloudability_account.cloudability.id }
output "role_name" { value = data.cloudability_account.cloudability.role_name }
output "external_id" { value = data.cloudability_account.cloudability.external_id }