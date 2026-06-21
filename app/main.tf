terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
  backend "local" {}
}

provider "docker" {}

variable "jwt_secret" {
  type      = string
  sensitive = true
  default   = "motherfather"
}

variable "captcha_secret_key" {
  type      = string
  sensitive = true
  default   = "1x0000000000000000000000000000000AA"
}

variable "postgres_password" {
  type      = string
  sensitive = true
  default   = "admin123"
}

variable "rabbitmq_password" {
  type      = string
  sensitive = true
  default   = "admin123"
}


data "docker_image" "user_service" {
  name = "spotnearr_user:latest"
}

resource "docker_container" "user_service" {
  name    = "spotnearr_user_svc"
  image   = data.docker_image.user_service.id
  restart = "unless-stopped"

  env = [
    "PORT=8080",
    "JWT_SECRET=${var.jwt_secret}",
    "DATABASE_URL=postgresql://admin:${var.postgres_password}@spotnearr_postgres:5432/spotnearr_user?sslmode=disable",
    "CACHE_URL=redis://spotnearr_redis:6379",
    "CACHE_PASSWORD=",
    "CAPTCHA_SECRET_KEY=${var.captcha_secret_key}",
    "CAPTCHA_URL=https://challenges.cloudflare.com/turnstile/v0/siteverify",
    "VENDOR_SERVICE_URL=http://spotnearr_vendor_svc:8080",
    "RABBITMQ_URL=amqp://admin:${var.rabbitmq_password}@spotnearr_rabbitmq:5672/",
  ]

  networks_advanced {
    name    = "spotnearr_network"
    aliases = ["user-service"]
  }
}
