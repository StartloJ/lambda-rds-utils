data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = [
      "sts:AssumeRole"
    ]
  }
}

data "aws_iam_policy_document" "allow_logging" {
  policy_id = "AllowLambdaPushLog"
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["arn:aws:logs:*:*:*"]
  }
}

data "aws_iam_policy_document" "allow_copy_rds_snapshot" {
  policy_id = "AllowLambdaCopyRdsSnapshot"
  statement {
    effect = "Allow"
    sid    = "AllowLambdaAccessToRdsSnapshot"
    actions = [
      "rds:CopyDBSnapshot",
      "rds:ModifyDBSnapshot",
      "rds:DescribeDBSnapshots",
      "rds:ModifyDBSnapshotAttribute"
    ]
    resources = ["*"]
  }

  statement {
    sid    = "AllowLambdaAccessToKMSKey"
    effect = "Allow"
    actions = [
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:ReEncrypt*",
      "kms:GenerateDataKey*"
    ]
    resources = ["${aws_kms_key.rds_snap.arn}"]
  }
}

resource "aws_iam_policy" "lambda_logging" {
  name   = "lambda-logging-cloudwatch"
  policy = data.aws_iam_policy_document.allow_logging.json
}

resource "aws_iam_role" "lambda_iam_role" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
}

resource "aws_iam_role_policy_attachment" "func_log_policy" {
  role       = aws_iam_role.lambda_iam_role.id
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_lambda_function" "rds_snap" {
  function_name = local.lambda_func_name
  description   = "Func to handle event from EventBridge to request RDS copy snapshot."
  role          = aws_iam_role.lambda_iam_role.arn

  handler   = "hello"
  runtime   = "go1.x"
  s3_bucket = module.s3_lambda.s3_bucket_id
  s3_key    = local.lambda_s3_key

  package_type = "Zip"
  memory_size  = 128

  environment {
    variables = {
      OPT_SRC_REGION        = "ap-southeast-1"
      OPT_TARGET_REGION     = "ap-southeast-2"
      OPT_BEDUG             = "True"
      OPT_DB_NAME           = ""
      OPT_OPTION_GROUP_NAME = ""
      OPT_KMS_KEY_ID        = ""
    }
  }

  tags = local.tags
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rds_snap.function_name
  principal     = "events.amazonaws.com"
  source_arn    = module.eventbridge.eventbridge_rule_arns.orders
}