resource "cloudability_account" "cloudability" {
  account_id = "${var.account_id}"
  role_arn   = "${var.role_arn}"
}

variable "role_arn" { type = string }
variable "account_id" { type = string }
