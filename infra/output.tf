output "server_ecr_repository_url" {
  value = aws_ecr_repository.server_ecr_repository.repository_url
}

output "cli_ecr_repository_url" {
  value = aws_ecr_repository.cli_ecr_repository.repository_url
}

output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "server_ecs_service_name" {
  value = aws_ecs_service.server-service.name
}

output "cli_ecs_service_name" {
  value = aws_ecs_service.cli-service.name
}

output "server_task_definition_arn" {
  value = aws_ecs_task_definition.server-task-definition.arn
}

output "cli_task_definition_arn" {
  value = aws_ecs_task_definition.cli-task-definition.arn
}

