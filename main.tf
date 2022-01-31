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
