resource "aws_security_group" "server_ecs_tasks_sg" {
  name        = "instant-term-server-ecs-taks-sg"
  description = "Allow inbound traffic from 9090 port"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "instant-term-server-ecs-tasks-sg"
    Environment = var.environment
  }
}

resource "aws_security_group" "cli_ecs_tasks_sg" {
  name        = "instant-term-cli-ecs-taks-sg"
  description = "Allow inbound traffic from http port"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "instant-term-cli-ecs-tasks-sg"
    Environment = var.environment
  }
}
