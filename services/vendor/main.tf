terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
  backend "gcs" {}
}

variable "project_id"           {}
variable "project_number"       {}
variable "region"               {}
variable "service_account"      {}
variable "image"                {}
variable "vendor_database_url"  {}
variable "cache_url"            {}
variable "captcha_url"          {}
variable "captcha_secret_key"   {}
variable "jwt_secret"           {}
variable "cache_password"       { default = "" }
variable "environment"          { default = "prod" }

provider "google" {
  project = var.project_id
  region  = var.region
}

data "google_cloud_run_v2_service" "user" {
  name     = "app-${var.environment}"
  location = var.region
  project  = var.project_id
}

resource "google_cloud_run_v2_service" "vendor" {
  name     = "vendor-${var.environment}"
  location = var.region
  project  = var.project_id

  template {
    service_account = var.service_account

    containers {
      image = var.image

      env {
        name  = "DATABASE_URL"
        value = var.vendor_database_url
      }
      env {
        name  = "CACHE_URL"
        value = var.cache_url
      }
      env {
        name  = "CAPTCHA_URL"
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
        name  = "APP_ENV"
        value = var.environment
      }
      env {
        name  = "USER_SERVICE_URL"
        value = data.google_cloud_run_v2_service.user.uri
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
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.vendor.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "service_url" {
  value = google_cloud_run_v2_service.vendor.uri
}
