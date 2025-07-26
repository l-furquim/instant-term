resource "aws_cloudwatch_log_group" "ecs_logs" {
  name = "/ecs/instant-term"
  retention_in_days = 30

  tags = {
    Name = "instant-term-logs"
    Environment = var.environment
  }
}