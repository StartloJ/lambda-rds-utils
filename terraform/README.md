# Terraform support this lab

This terraform should to provision resources to served this lambda function.
It provide related resources like S3, KMS, EventBridge, and other.

## Resources created

1. S3 bucket
2. KMS multi-region
3. Lambda functions
4. Cloudwatch Log group for Lambda function
5. EventBridge bus and rules

## Step to use

1. Change locals parameter in `main.tf` to served.
2. Run `terraform apply`
3. Check resources on AWS console

## Tip note

> For the events from AWS services, you can use default bus instead.
> Because the AWS services will use a default bus to push their events into
> the EventBridge. But this lab will create a new event bus cause to proof
> my concept without effect to my current project running.

## Ref
- https://dev.to/aws-builders/creating-a-multi-region-key-using-terraform-51o4
- https://registry.terraform.io/providers/hashicorp/aws/5.10.0/docs/resources
- https://stackoverflow.com/questions/59949808/write-aws-lambda-logs-to-cloudwatch-log-group-with-terraform
- https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html
