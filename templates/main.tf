terraform {
  required_providers {
    aws-custom = {
      source = "local/custom/aws-custom"
      version = "1.0.0"
    }
  }
}

provider "aws-custom" {}

resource "aws_custom_ec2_instance" "example" {
  provider      = aws-custom
  instance_type = "t2.micro"
  ami           = "ami-12345678"
}