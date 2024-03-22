resource "aws_cloudwatch_log_group" "registry_lambda" {
  name = "/aws/lambda/${aws_lambda_function.registry.function_name}"
  retention_in_days = 30
}

resource "aws_cloudwatch_log_group" "registry_api" {
  name = "/aws/api-gw/${aws_apigatewayv2_api.registry.name}"
  retention_in_days = 30
}
