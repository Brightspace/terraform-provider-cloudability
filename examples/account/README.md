# Cloudability Account

This example describes how to configure an AWS account in Cloudability.

## How to Run

You will need to have setup Evident API credentials, which you learn about in the [usage section](../../README.md). You will also need an IAM role setup for this to work.

```bash
terraform apply \
        -var="write=0" \
        -var="account_id=123412341234" \
        -var="role_arn=arn:aws:iam::123412341234:role/CloudabilityRole"
```

The `write` variable is used. The example should be run with `write=0`, which will output the expected `external_id`. You can then set this in the IAM Role `CloudabilityRole`. When you run with `write=1`, it'll configure the account.
