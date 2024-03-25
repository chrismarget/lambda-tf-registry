resource "aws_apigatewayv2_api" "registry" {
  name          = "registry"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "all" {
  #  description = "created by terraform"
  api_id      = aws_apigatewayv2_api.registry.id
  name        = "all"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.registry_api.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
}

resource "aws_apigatewayv2_integration" "registry" {
  description        = "registry lambda"
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  api_id             = aws_apigatewayv2_api.registry.id
  integration_uri    = aws_lambda_function.registry.invoke_arn
}

#resource "aws_apigatewayv2_route" "default" {
#  api_id = aws_apigatewayv2_api.registry.id
#  route_key = "$default"
#  target    = "integrations/${aws_apigatewayv2_integration.registry.id}"
#}

resource "aws_apigatewayv2_route" "well_known" {
  api_id = aws_apigatewayv2_api.registry.id
  route_key = "GET /.well-known/terraform.json"
  target    = "integrations/${aws_apigatewayv2_integration.registry.id}"
}

resource "aws_apigatewayv2_route" "v1" {
  api_id = aws_apigatewayv2_api.registry.id
    route_key = "GET /v1/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.registry.id}"
}

resource "aws_apigatewayv2_domain_name" "tf_registry_click" {
  domain_name = data.aws_route53_zone.tf_registry_click.name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.tf_registry_click.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "tf_registry_click" {
  api_id      = aws_apigatewayv2_api.registry.id
  domain_name = aws_apigatewayv2_domain_name.tf_registry_click.id
  stage       = aws_apigatewayv2_stage.all.id
}

locals {
  dnc = aws_apigatewayv2_domain_name.tf_registry_click.domain_name_configuration
  dnc_names = local.dnc[*].target_domain_name
  dnc_map = zipmap(local.dnc_names, local.dnc)
}

resource "aws_route53_record" "tf_registry_click" {
  for_each = local.dnc_map
  zone_id = data.aws_route53_zone.tf_registry_click.zone_id
  name    = data.aws_route53_zone.tf_registry_click.name
  type    = "A"

  alias {
    evaluate_target_health = false
    name                   = each.value["target_domain_name"]
    zone_id                = each.value["hosted_zone_id"]
  }
}
