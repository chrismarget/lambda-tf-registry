data "aws_route53_zone" "tf_registry_click" {
  name = "tf-registry.click"
}

resource "aws_acm_certificate" "tf_registry_click" {
  domain_name       = data.aws_route53_zone.tf_registry_click.name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.tf_registry_click.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.tf_registry_click.zone_id
}

resource "aws_acm_certificate_validation" "tf_registry_click" {
  certificate_arn         = aws_acm_certificate.tf_registry_click.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}
