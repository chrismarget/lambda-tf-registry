locals {
  repo_dir          = "${path.root}/.."
  registry_main_dir = "${local.repo_dir}/cmd/registry"
  registry_main_go  = "${local.registry_main_dir}/main.go"
  all_src_dirs = [
    local.registry_main_dir,
    "${local.repo_dir}/common",
    "${local.repo_dir}/v1_handlers",
  ]
  tmp_dir = "${path.root}/../.temp"
}

# zip of the source tree is used to trigger lambda build when code is changed
data "archive_file" "lambda_source" {
  for_each    = toset(local.all_src_dirs)
  source_dir  = each.key
  output_path = "${local.tmp_dir}/${replace(each.key, "/", "_")}.zip"
  type        = "zip"
}

resource "terraform_data" "build_lambda" {
  triggers_replace = {
    for i in data.archive_file.lambda_source : i.source_dir => i.output_sha256
  }

  provisioner "local-exec" {
    command = "go build -o ${local.tmp_dir}/bootstrap ${local.registry_main_go}"
    environment = {
      CGO_ENABLED = "0"
      GOOS        = "linux"
      GOARCH      = "arm64"
    }
  }
}

data "archive_file" "zip_lambda" {
  source_file = "${path.root}/../.temp/bootstrap"
  output_path = "${path.root}/../.temp/bootstrap.zip"
  type        = "zip"
  depends_on  = [terraform_data.build_lambda]
}

resource "aws_lambda_function" "registry" {
  function_name = "registry"
  role          = aws_iam_role.registry.arn
  filename      = data.archive_file.zip_lambda.output_path
  runtime       = "provided.al2023"
  handler       = "registry"
  architectures = ["arm64"]
  depends_on    = [terraform_data.build_lambda]
  #  publish       = true

  environment {
    variables = {
      DEBUG               = "1"
      PROVIDER_TABLE_NAME = aws_dynamodb_table.registry_providers.name
    }
  }

  lifecycle { replace_triggered_by = [terraform_data.build_lambda] }
}

resource "aws_lambda_permission" "registry_url" {
  function_name          = aws_lambda_function.registry.function_name
  action                 = "lambda:InvokeFunctionUrl"
  principal              = "*"
  function_url_auth_type = "NONE"

  lifecycle { replace_triggered_by = [aws_lambda_function.registry] }
}

resource "aws_lambda_permission" "api_gateway_b" {
  function_name = aws_lambda_function.registry.function_name
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.registry.execution_arn}/*/*"

  lifecycle { replace_triggered_by = [aws_lambda_function.registry] }
}

resource "aws_lambda_function_url" "registry" {
  authorization_type = "NONE"
  function_name      = aws_lambda_function.registry.function_name
}
