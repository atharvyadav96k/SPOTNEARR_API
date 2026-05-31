terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }

  backend "gcs" {}
}

variable "project_id"      {}
variable "project_number"  {}
variable "region"          {}
variable "service_account" {}
variable "backend_url"     {}

provider "google" {
  project = var.project_id
  region  = var.region
}


provider "google-beta" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "apigateway" {
  project            = var.project_id
  service            = "apigateway.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "servicemanagement" {
  project            = var.project_id
  service            = "servicemanagement.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "servicecontrol" {
  project            = var.project_id
  service            = "servicecontrol.googleapis.com"
  disable_on_destroy = false
}

resource "google_api_gateway_api" "app" {
  provider = google-beta
  project  = var.project_id
  api_id   = "spotnearr-api"

  depends_on = [
    google_project_service.apigateway,
    google_project_service.servicemanagement,
    google_project_service.servicecontrol,
  ]
}

resource "google_api_gateway_api_config" "app" {
  provider      = google-beta
  project       = var.project_id
  api           = google_api_gateway_api.app.api_id
  api_config_id = "spotnearr-config-${substr(sha256(var.backend_url), 0, 8)}"

  openapi_documents {
    document {
      path = "openapi.yaml"
      contents = base64encode(templatefile("${path.module}/openapi.yaml", {
        backend_url = var.backend_url
        project_id  = var.project_id
      }))
    }
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [google_api_gateway_api.app]
}

resource "google_api_gateway_gateway" "app" {
  provider   = google-beta
  project    = var.project_id
  region     = var.region
  api_config = google_api_gateway_api_config.app.id
  gateway_id = "spotnearr-gateway"

  depends_on = [google_api_gateway_api_config.app]
}

output "gateway_url" {
  value = "https://${google_api_gateway_gateway.app.default_hostname}"
}