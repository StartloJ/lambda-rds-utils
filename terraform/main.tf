locals {
  region = "ap-southeast-1"

  # For S3
  S3_prefix = "lambda-lake-"

  # For Lambda components
  lambda_func_name = "rds_copy_snap"
  lambda_s3_key    = "lambda-hello.zip"
  handler_bin_name = "hello"
  func_env_var = {
    OPT_SRC_REGION        = "ap-southeast-1"
    OPT_TARGET_REGION     = "ap-southeast-2"
    OPT_BEDUG             = "True"
    OPT_DB_NAME           = ""
    OPT_OPTION_GROUP_NAME = ""
    OPT_KMS_KEY_ID        = ""
  }

  # For EventBridge
  bus_name = "rds-bus" #Leave with `default` to use AWS service event bus.
  rds_events = {
    description = "Capture for RDS snapshot event"
    event_pattern = jsonencode({
      "source" : ["aws.rds", "demo.event"],
      "detail" : {
        "EventID" : ["RDS-EVENT-0042", "RDS-EVENT-0091"]
      }
    })
  }

  tags = {
    Owner   = "watcharin"
    Project = "local-dev"
  }
}

data "aws_caller_identity" "current" {}
