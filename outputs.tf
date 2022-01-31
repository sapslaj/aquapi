output "cloudformation_stack_name" {
  value = join("-", [local.config.service, var.stage])
}
