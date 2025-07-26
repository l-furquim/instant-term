resource "aws_ecr_repository" "server_ecr_repository" {
  name = "instant-term-server"

  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "cli_ecr_repository" {
  name = "instant-term-cli"

  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}