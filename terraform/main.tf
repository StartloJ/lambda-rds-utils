locals {
  region = "ap-southeast-1"

  lambda_func_name = "rds_copy_snap"
  lambda_s3_key    = "lambda-hello.zip"

  tags = {
    Owner   = "watcharin"
    Project = "local-dev"
  }
}

data "aws_caller_identity" "current" {}
