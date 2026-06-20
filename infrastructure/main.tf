variable "state_path" {
  type = string
}

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
  backend "local" {
    path = var.state_path
  }
}

provider "docker" {}


resource "docker_network" "app_network" {
  name = "spotnearr_network"
}

resource "docker_image" "nginx" {
  name = "nginx:alpine"
}

resource "docker_image" "rabbitmq" {
  name = "rabbitmq:alpine"
}

resource "docker_image" "postgres" {
  name = "postgres:16-alpine"
}

resource "docker_image" "redis" {
  name = "redis:7-alpine"
}

resource "docker_image" "typesense" {
  name = "typesense/typesense:26.0"
}

