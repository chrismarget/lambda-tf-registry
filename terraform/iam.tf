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
#  statement {
#    effect    = "Allow"
#    actions   = ["logs:CreateLogGroup"]
#    resources = ["arn:aws:logs:${data.aws_region.ours.id}:${data.aws_caller_identity.ours.id}:*"]
#  }
#
#  statement {
#    effect    = "Allow"
#    actions   = ["logs:CreateLogStream", "logs:PutLogEvents"]
#    resources = ["arn:aws:logs:${data.aws_region.ours.id}:${data.aws_caller_identity.ours.id}:log-group:/aws/lambda/${aws_lambda_function.registry.function_name}:*"]
#  }

  statement {
    effect    = "Allow"
    actions   = ["dynamodb:GetItem", "dynamodb:Query"]
    resources = [
      "${aws_dynamodb_table.registry_providers.arn}/index/*",
      aws_dynamodb_table.registry_providers.arn,
    ]
  }

#  statement {
#    effect = "Allow"
#    actions = ["secretsmanager:GetSecretValue"]
#    resources = [ aws_secretsmanager_secret.registry_uploader.arn ]
#  }
}

// arn:aws:dynamodb:us-east-1:086704128018:table/registry-providers
// arn:aws:dynamodb:us-east-1:086704128018:table/registry-providers

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
