data "aws_iam_policy_document" "registry_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "registry" {
  name               = "terraform_registry"
  assume_role_policy = data.aws_iam_policy_document.registry_assume_role.json
}

data "aws_iam_policy_document" "dynamodb" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:GetItem", "dynamodb:Query"]
    resources = [
      "${aws_dynamodb_table.registry_providers.arn}/index/*",
      aws_dynamodb_table.registry_providers.arn,
    ]
  }
}

resource "aws_iam_policy" "read_dynamodb" {
  name   = "terraform_registry"
  policy = data.aws_iam_policy_document.dynamodb.json
}

resource "aws_iam_role_policy_attachment" "our_policy" {
  policy_arn = aws_iam_policy.read_dynamodb.arn
  role       = aws_iam_role.registry.name
}

resource "aws_iam_role_policy_attachment" "AWSLambdaBasicExecutionRole" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.registry.name
}

data "aws_iam_policy_document" "registry_uploader"{
  statement {
    effect = "Allow"
    actions = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.registry_providers.arn]
  }
  statement {
    effect = "Allow"
    actions = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.registry.arn}/*"]
  }
}

locals {
  registry_uploaders = toset([
    "chris",
  ])
}

resource "aws_iam_user" "registry_uploader" {
  for_each = local.registry_uploaders
  name = "registry_uploader-${each.key}"
}

resource "aws_iam_user_policy" "registry_uploader" {
  for_each = aws_iam_user.registry_uploader
  policy = data.aws_iam_policy_document.registry_uploader.json
  user   = each.value.id
}
