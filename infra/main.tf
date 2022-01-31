data "aws_region" "current" {}

data "cloudflare_zone" "this" {
  name = var.cloudflare_zone
}

module "images_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> v2.13.0"

  acl           = "public-read"
  bucket_prefix = "aquapi-images"
  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}

resource "aws_s3_bucket_object" "root_files" {
  for_each = {
    "index.html" = {
      content_type = "text/html"
    }
    "error.html" = {
      content_type = "text/html"
    }
    "favicon.ico" = {
      content_type = "image/x-icon"
    }
  }

  key          = each.key
  source       = "${path.module}/${each.key}"
  content_type = lookup(each.value, "content_type", null)

  bucket = module.images_bucket.s3_bucket_id
  acl    = "public-read"
}

data "aws_acm_certificate" "wildcard" {
  domain = "*.sapslaj.com"
}

# Have to use this because _FUCKING_ S3 doesn't support TLS with bare website
# hosting. WHY, AMAZON????
module "cloudfront" {
  source  = "terraform-aws-modules/cloudfront/aws"
  version = "~> 2.9.2"

  enabled     = true
  price_class = "PriceClass_100"
  aliases     = [var.images_domain]

  default_root_object = "index.html"

  origin = {
    s3 = {
      domain_name      = module.images_bucket.s3_bucket_bucket_regional_domain_name
      s3_origin_config = {}
    }
  }

  default_cache_behavior = {
    target_origin_id       = "s3"
    viewer_protocol_policy = "allow-all"
  }

  viewer_certificate = {
    acm_certificate_arn = data.aws_acm_certificate.wildcard.arn
    ssl_support_method  = "sni-only"
  }
}

resource "cloudflare_record" "images" {
  zone_id = data.cloudflare_zone.this.id
  name    = var.images_domain
  value   = module.cloudfront.cloudfront_distribution_domain_name
  type    = "CNAME"
  ttl     = 1
  proxied = true
}

module "serverless" {
  source = "../"

  aws_region        = data.aws_region.current.name
  api_domain        = var.api_domain
  images_bucket_arn = module.images_bucket.s3_bucket_arn
  images_bucket_id  = module.images_bucket.s3_bucket_id
  images_domain     = var.images_domain
  stage             = var.stage
}

data "aws_cloudformation_stack" "serverless" {
  name = module.serverless.cloudformation_stack_name
}

resource "aws_apigatewayv2_domain_name" "api" {
  domain_name = var.api_domain

  domain_name_configuration {
    certificate_arn = data.aws_acm_certificate.wildcard.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api" {
  api_id      = data.aws_cloudformation_stack.serverless.outputs["HttpApiId"]
  domain_name = aws_apigatewayv2_domain_name.api.id
  stage       = "$default"
}

resource "cloudflare_record" "api" {
  zone_id = data.cloudflare_zone.this.id
  name    = var.api_domain
  value   = replace(data.aws_cloudformation_stack.serverless.outputs["HttpApiUrl"], "https://", "")
  type    = "CNAME"
  ttl     = 1
  proxied = true
}
