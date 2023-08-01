module "eventbridge" {
  source  = "terraform-aws-modules/eventbridge/aws"
  version = "2.3.0"

  bus_name = local.bus_name
  rules = {
    orders = local.rds_events
  }

  # attach_lambda_policy = true
  # lambda_target_arns   = ["${aws_lambda_function.rds_snap.arn}"]

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