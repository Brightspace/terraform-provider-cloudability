data "cloudability_account" "cloudability" {
  account_id = "${var.account_id}"
}

resource "cloudability_account" "cloudability" {
  count      = "${var.write}"
  account_id = "${var.account_id}"
  role_arn   = "${var.role_arn}"

  depends_on = [data.cloudability_account.cloudability]
}

variable "write" {
  type    = string
  default = 0
}
variable "role_arn" { type = string }
variable "account_id" { type = string }

output "id" { value = data.cloudability_account.cloudability.id }
output "role_name" { value = data.cloudability_account.cloudability.role_name }
output "external_id" { value = data.cloudability_account.cloudability.external_id }
