terraform {
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "sapslaj"

    workspaces {
      name = "aquapi-production"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

module "infra" {
  source = "../"

  api_domain    = "aquapi.sapslaj.com"
  images_domain = "aquapic.sapslaj.com"
  stage         = "production"
}
