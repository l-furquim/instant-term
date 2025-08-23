terraform {
  backend "s3" {
    bucket = "tfstate-710271919573"
    key = "instant-term"
    region = "us-east-1"
  }
}