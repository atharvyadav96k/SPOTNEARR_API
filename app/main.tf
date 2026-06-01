terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }

  backend "gcs" {}
}

variable "project_id"      {}
variable "project_number"  {}
variable "region"          {}
variable "service_account" {}
variable "image"           {}
variable "database_url"    {}
variable "cache_url"       {}
variable "captcha_url"     {}
variable "captcha_secret_key" {}
variable "jwt_secret"      {}
variable "cache_password"  { default = "" }
variable "port"            { default = "8080" }

provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_client_config" "default" {}

provider "docker" {
  registry_auth {
    address  = "${var.region}-docker.pkg.dev"
    username = "oauth2accesstoken"
    password = data.google_client_config.default.access_token
  }
}

data "google_artifact_registry_repository" "app" {
  project       = var.project_id
  location      = var.region
  repository_id = "gcf-artifacts"
}

resource "docker_image" "app" {
  name = var.image

  build {
    context    = path.module
    dockerfile = "${path.module}/Dockerfile"
    platform   = "linux/amd64"
  }

  depends_on = [data.google_artifact_registry_repository.app]
}

resource "docker_registry_image" "app" {
  name          = docker_image.app.name
  keep_remotely = true

  depends_on = [docker_image.app]
}

resource "google_cloud_run_v2_service" "app" {
  name     = "app"
  location = var.region
  project  = var.project_id

  template {
    service_account = var.service_account

    containers {
      image = "${var.image}@${docker_registry_image.app.sha256_digest}"

      env {
        name  = "DATABASE_URL"
        value = var.database_url
      }

      env {
        name = "CACHE_URL"
        value = var.cache_url
      }


      env {
        name = "CAPTCHA_URL"
        value = var.captcha_url
      }


      env {
        name  = "CAPTCHA_SECRET_KEY"
        value = var.captcha_secret_key
      }

      env {
        name  = "JWT_SECRET"
        value = var.jwt_secret
      }

      env {
        name  = "CACHE_PASSWORD"
        value = var.cache_password
      }

      env {
        name  = "PORT"
        value = var.port
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }
  }

  depends_on = [docker_registry_image.app]
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}


output "service_url" {
  value = google_cloud_run_v2_service.app.uri
}

output "image_digest" {
  value = docker_registry_image.app.sha256_digest
}