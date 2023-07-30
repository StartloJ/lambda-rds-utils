module "eventbridge" {
  source  = "terraform-aws-modules/eventbridge/aws"
  version = "2.3.0"

  bus_name = "rds-bus"
  rules = {
    orders = {
      description = "Capture for RDS snapshot event"
      event_pattern = jsonencode({
        "source" : ["aws.rds", "demo.event"]
      })
    }
  }

  attach_lambda_policy = true
  lambda_target_arns   = ["${aws_lambda_function.rds_snap.arn}"]

  targets = {
    orders = [
      {
        name = "event-to-lambda"
        arn  = "${aws_lambda_function.rds_snap.arn}"
      }
    ]
  }

  tags = local.tags

  depends_on = [aws_lambda_function.rds_snap]
}