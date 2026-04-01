terraform {
  backend "gcs" {
    bucket  = "terraform-state-603675804309"
    prefix  = "cloud-functions/demo"
  }
  required_providers {
    google = { source = "hashicorp/google" }
    archive = { source = "hashicorp/archive" }
  }
}

variable "project_id" { type = string }
variable "region" { type = string }
variable "service_account" { type = string }

provider "google" {
  project = var.project_id
  region  = var.region
}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = "${path.module}/function"
  output_path = "${path.module}/source.zip"
}

resource "google_storage_bucket" "source_bucket" {
  name     = "${var.project_id}-gcf-source"
  location = var.region
  force_destroy = true
}

resource "google_storage_bucket_object" "source_archive" {
  name   = "source-${data.archive_file.source.output_md5}.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = data.archive_file.source.output_path
}

resource "google_cloudfunctions2_function" "helloworld" {
  name     = "HelloWorld"
  location = var.region

  build_config {
    runtime     = "go122" 
    entry_point = "HelloWorld"
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
    service_account_email = var.service_account
    ingress_settings      = "ALLOW_ALL"
  }
}

resource "google_cloud_run_service_iam_member" "public_access" {
  location = var.region
  service  = google_cloudfunctions2_function.helloworld.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "function_url" {
  value = google_cloudfunctions2_function.helloworld.service_config[0].uri
}