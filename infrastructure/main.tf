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

variable "typesense_api_key" {
  type      = string
  sensitive = true
  default   = "dev-api-key"
}


resource "docker_network" "app_network" {
  name = "spotnearr_network"
}

resource "docker_volume" "pgdata" {
  name = "spotnearr_pgdata"
}

resource "docker_volume" "redis_data" {
  name = "spotnearr_redis_data"
}

resource "docker_volume" "rabbitmq_data" {
  name = "spotnearr_rabbitmq_data"
}

resource "docker_volume" "typesense_data" {
  name = "spotnearr_typesense_data"
}

resource "docker_image" "postgres" {
  name = "postgres:16-alpine"
}

resource "docker_image" "redis" {
  name = "redis:7-alpine"
}

resource "docker_image" "rabbitmq" {
  name = "rabbitmq:3.13-management-alpine"
}

resource "docker_image" "typesense" {
  name = "typesense/typesense:26.0"
}

resource "docker_container" "postgres" {
  name    = "spotnearr_postgres"
  image   = docker_image.postgres.image_id
  restart = "unless-stopped"

  env = [
    "POSTGRES_USER=admin",
    "POSTGRES_PASSWORD=${var.postgres_password}",
    "POSTGRES_DB=spotnearr_user",
  ]

  volumes {
    volume_name    = docker_volume.pgdata.name
    container_path = "/var/lib/postgresql/data"
  }

  volumes {
    host_path      = abspath("${path.module}/../docker/init-db.sql")
    container_path = "/docker-entrypoint-initdb.d/init.sql"
    read_only      = true
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  healthcheck {
    test     = ["CMD-SHELL", "pg_isready -U admin -d spotnearr_user"]
    interval = "5s"
    timeout  = "5s"
    retries  = 10
  }
}

resource "docker_container" "redis" {
  name    = "spotnearr_redis"
  image   = docker_image.redis.image_id
  restart = "unless-stopped"

  volumes {
    volume_name    = docker_volume.redis_data.name
    container_path = "/data"
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  healthcheck {
    test     = ["CMD", "redis-cli", "ping"]
    interval = "5s"
    timeout  = "3s"
    retries  = 10
  }
}


resource "docker_container" "rabbitmq" {
  name    = "spotnearr_rabbitmq"
  image   = docker_image.rabbitmq.image_id
  restart = "unless-stopped"

  env = [
    "RABBITMQ_DEFAULT_USER=admin",
    "RABBITMQ_DEFAULT_PASS=${var.rabbitmq_password}",
  ]

  volumes {
    volume_name    = docker_volume.rabbitmq_data.name
    container_path = "/var/lib/rabbitmq"
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  healthcheck {
    test         = ["CMD", "rabbitmq-diagnostics", "check_port_connectivity"]
    interval     = "10s"
    timeout      = "5s"
    retries      = 15
    start_period = "20s"
  }
}

resource "docker_container" "typesense" {
  name    = "spotnearr_typesense"
  image   = docker_image.typesense.image_id
  restart = "unless-stopped"

  command = [
    "--data-dir=/data",
    "--api-key=${var.typesense_api_key}",
    "--enable-cors",
  ]

  volumes {
    volume_name    = docker_volume.typesense_data.name
    container_path = "/data"
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  healthcheck {
    test         = ["CMD", "bash", "-c", "exec 3<>/dev/tcp/127.0.0.1/8108"]
    interval     = "5s"
    timeout      = "3s"
    retries      = 15
    start_period = "10s"
  }
}