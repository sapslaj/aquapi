data "aws_iam_policy_document" "role" {
  statement {
    actions = [
      "s3:ListBucket",
      "s3:GetBucketLocation",
    ]
    resources = [sensitive(var.images_bucket_arn)]
  }
}

locals {
  config = {
    service          = "aquapi"
    frameworkVersion = "3"

    provider = {
      name    = "aws"
      runtime = "go1.x"
      stage   = var.stage
      region  = var.aws_region
      iam = {
        role = {
          statements = jsondecode(data.aws_iam_policy_document.role.json)["Statement"]
        }
      }
      environment = {
        AQUAPI_IMAGES_HOST   = var.images_domain
        AQUAPI_IMAGES_BUCKET = sensitive(var.images_bucket_id)
      }
      httpApi = {
        cors = true
      }
    }

    package = {
      patterns = [
        "!./**",
        "./out/**",
      ]
    }

    functions = {
      images = {
        handler = "out/images"
        events = [{
          httpApi = {
            path   = "/images"
            method = "get"
          }
        }]
      }
      lucky = {
        handler = "out/lucky"
        events = [{
          httpApi = {
            path   = "/images/lucky"
            method = "get"
          }
          }, {
          httpApi = {
            path   = "/"
            method = "get"
          }
        }]
      }
    }
  }
}
