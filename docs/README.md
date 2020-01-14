# Cloudability Provider

The Cloudability provider is used to interact with Cloudability AWS resources.

The provider allows you to manage your Cloudability AWS connections. It needs to be configured with the proper credentials before it can be used.

Use the navigation below to read about the available resources.

## Example Usage

```hcl
# Configure the Cloudability Provider
provider "cloudability" {}

# Get the account from cloudability
data "cloudability_account" "cloudability" {
  account_id = "987654321012"
}

```

## Authentication

The Cloudability provider can be provided credentials for authentication using environment variables or static credentials.

### Static credentials

| :warning: Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system :warning: |
| --- |

Static credentials can be provided by adding an `api_key` and `payer_account_id` in-line in the AWS provider block or by variables:

Usage:

```hcl
provider "cloudability" {
  api_key          = "my-api-key"
  payer_account_id = "123456789012"
}
```
### Environment variables

You can provide your credentials via the `CLOUDABILITY_TOKEN` and `CLOUDABILITY_PAYER_ACCOUNT_ID`, environment variables, representing your Cloudability API Token and AWS Billing Account ID, respectively.

```hcl
provider "aws" {}
```

Usage:

```bash
export CLOUDABILITY_TOKEN="my-api-key"
export CLOUDABILITY_PAYER_ACCOUNT_ID="123456789012"
terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the Cloudability provider block:

- `api_key` - (Optional) This is the Cloudability API token. It must be provided, but it can also be sourced from the `CLOUDABILITY_PAYER_ACCOUNT_ID` environment variable.
- `payer_account_id` - (Optional) The 12-digit AWS identifier of the payer account. It must be provided, but it can also be sourced from the `CLOUDABILITY_PAYER_ACCOUNT_ID` environment variable.