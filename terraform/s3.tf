resource "aws_s3_bucket" "registry" {
  bucket = "jtaf-registry"
}

data "aws_iam_policy_document" "registry_bucket" {
  statement {
    effect     = "Allow"
    actions    = ["s3:GetObject"]
    resources  = ["${aws_s3_bucket.registry.arn}/*"]

    principals {
      identifiers = ["*"]
      type        = "*"
    }
  }
}

resource "aws_s3_bucket_policy" "my_bucket_policy" {
  bucket = aws_s3_bucket.registry.bucket
  policy = data.aws_iam_policy_document.registry_bucket.json
}
