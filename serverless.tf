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
    }

    package = {
      patterns = [
        "!./**",
        "./out/**",
      ]
    }

    functions = {
      random = {
        handler = "out/random"
        events = [{
          httpApi = {
            path   = "/"
            method = "get"
          }
        }]
      }
    }
  }
}

resource "local_file" "serverless_config" {
  content  = jsonencode(local.config)
  filename = "${path.module}/serverless.json"
}

resource "null_resource" "make_build" {
  triggers = {
    always = timestamp()
  }

  provisioner "local-exec" {
    command     = "make build"
    working_dir = path.module
  }
}

resource "null_resource" "serverless_deploy" {
  triggers = {
    always = timestamp()
  }

  provisioner "local-exec" {
    command     = "serverless deploy --verbose"
    working_dir = path.module
  }

  depends_on = [
    local_file.serverless_config,
    null_resource.make_build,
  ]
}
