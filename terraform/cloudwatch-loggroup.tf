resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${local.lambda_func_name}"
  retention_in_days = 1
  lifecycle {
    prevent_destroy = false
  }
}