resource "aws_ecs_cluster" "main" {
  name = "instant-term-cluster"

  tags = {
    Name        = "instant-term-ecs-cluster"
    Environment = var.environment
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "instant-term-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "server-task-definition" {
  family                   = "instant-term-server-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  depends_on = [aws_cloudwatch_log_group.ecs_logs]

  container_definitions = jsonencode([
    {
      name  = "instant-term-server-container"
      image = "${aws_ecr_repository.server_ecr_repository.repository_url}:latest"
      portMappings = [
        {
          containerPort = 9090
          hostPort      = 9090
          protocol      = "tcp"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/instant-term"
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
  tags = {
    Name        = "instant-term-server-task-definition"
    Environment = var.environment
  }
}

resource "aws_ecs_task_definition" "cli-task-definition" {
  family                   = "instant-term-cli-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  depends_on = [aws_cloudwatch_log_group.ecs_logs]

  container_definitions = jsonencode([
    {
      name  = "instant-term-cli-container"
      image = "${aws_ecr_repository.cli_ecr_repository.repository_url}:latest"
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
          protocol      = "tcp"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/instant-term"
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
  tags = {
    Name        = "instant-term-cli-task-definition"
    Environment = var.environment
  }
}

resource "aws_ecs_service" "server_service" {
  name            = "instant-term-server-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.server-task-definition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.server_ecs_tasks_sg.id]
    assign_public_ip = true
  }

  tags = {
    Name        = "instant-term-server-service"
    Environment = var.environment
  }

}

resource "aws_ecs_service" "cli_service" {
  name            = "instant-term-cli-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.cli-task-definition.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.cli_ecs_tasks_sg.id]
    assign_public_ip = true
  }

  tags = {
    Name        = "instant-term-cli-service"
    Environment = var.environment
  }

}
