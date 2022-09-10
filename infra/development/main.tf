provider "aws" {
  region = "us-east-1"
}

module "infra" {
  source = "../"

  api_domain    = "aquapi-development.sapslaj.com"
  images_domain = "aquapic-development.sapslaj.com"
  stage         = "development"
}
