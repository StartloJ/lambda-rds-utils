module "s3_lambda" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket_prefix     = local.S3_prefix
  create_bucket     = true
  block_public_acls = true

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  tags = local.tags
}
