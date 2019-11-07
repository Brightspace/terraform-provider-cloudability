# Cloudability Account

This example describes how to configure an AWS account in Cloudability.

## How to Run

You will need to have setup Evident API credentials, which you learn about in the [usage section](../../README.md). You will also need an IAM role setup for this to work.

```bash
terraform apply \
        -var="account_id=123412341234" \
        -var="role_arn=arn:aws:iam::123412341234:role/CloudabilityRole"
````