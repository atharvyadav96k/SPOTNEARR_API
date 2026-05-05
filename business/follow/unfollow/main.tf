variable "function_name" {
  type = string
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

variable "bucket_name" {
  type = string
}

variable "topic_name" {
  type = string
}

terraform {
  backend "gcs" {}

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

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

resource "google_storage_bucket" "source_bucket" {
  name                        = "${var.function_name}-${var.project_id}-source-${random_id.bucket_suffix.hex}"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "source_archive" {
  name   = "${var.function_name}-${filesha256("source.zip")}.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = "source.zip"
}

resource "google_pubsub_topic" "topic" {
  name = var.topic_name
}

resource "google_cloudfunctions2_function" "function" {
  name     = var.function_name
  location = var.region

  build_config {
    runtime         = "go122"
    entry_point     = "Unfollow"
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
  }

  event_trigger {
    trigger_region        = var.region
    event_type            = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic          = google_pubsub_topic.topic.id
    retry_policy          = "RETRY_POLICY_RETRY"
    service_account_email = var.service_account
  }
}

output "topic_id" {
  value = google_pubsub_topic.topic.id
}

output "function_name" {
  value = google_cloudfunctions2_function.function.name
}