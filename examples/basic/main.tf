provider "cloudability" {
  api_key          = "${var.api_key}"
  payer_account_id = "${var.payer_account_id}"
}

data "cloudability_account" "cloudability" {
  account_id       = "${var.account_id}"
}

resource "cloudability_account" "cloudability" {
  account_id = "${var.account_id}"
  role_arn   = "${var.role_arn}"
}
