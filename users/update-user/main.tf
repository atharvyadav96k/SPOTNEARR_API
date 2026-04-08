terraform {
  backend "gcs" {
    bucket  = "terraform-state-603675804309"
    prefix  = "cloud-functions/update-users"
  }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

variable "project_id" {
  type = string
}

variable "project_number" {
  type = string
}

variable "region" {
  type = string
}

variable "service_account" {
  type = string
}

variable "function_suffix" {
  type    = string
  default = ""
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

resource "google_storage_bucket" "source_bucket" {
  name                        = "${var.project_id}-function-${random_id.bucket_suffix.hex}"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "source_archive" {
  name   = "source-${filesha256("source.zip")}.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = "source.zip"
}

resource "google_cloudfunctions2_function" "update_user" {
  name = "update-user${var.function_suffix}"
  location = var.region

  build_config {
    runtime         = "go122"
    entry_point     = "UpdateUser" 
    service_account = "projects/${var.project_id}/serviceAccounts/${var.service_account}"
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.source_archive.name
      }
    }
  }

  service_config {
    max_instance_count    = 1
    available_memory      = "256M"
    timeout_seconds       = 60
    service_account_email = var.service_account
    ingress_settings      = "ALLOW_ALL"
  }
}

resource "google_cloud_run_service_iam_member" "public_access" {
  location = var.region
  service  = google_cloudfunctions2_function.update_user.name  
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "function_url" {
  value = google_cloudfunctions2_function.update_user.service_config[0].uri
} 