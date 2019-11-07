data "cloudability_account" "cloudability" {
  account_id = "${var.account_id}"
}

variable "account_id" { type = string }
output "id" { value = data.cloudability_account.cloudability.id }
