data "cloudability_account" "cloudability" {
  account_id       = "${var.account_id}"
}

output "id" { value = data.cloudability_account.cloudability.id }