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
