terraform {
  backend "gcs" {
    bucket  = "terraform-state-603675804309"
    prefix  = "cloud-functions/demo"
  }
  required_providers {
    google  = { source = "hashicorp/google" }
    archive = { source = "hashicorp/archive" }
    random  = { source = "hashicorp/random" }
  }
}

variable "project_id"      { type = string }
variable "region"          { type = string }
variable "service_account" { type = string }

provider "google" {
  project = var.project_id
  region  = var.region
}

# 1. Automatically zip the 'function' folder
data "archive_file" "source" {
  type        = "zip"
  source_dir  = "${path.module}/function"
  output_path = "${path.module}/source.zip"
}

# 2. Generate a random suffix for the bucket name
resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# 3. Create the source bucket with required uniform access
resource "google_storage_bucket" "source_bucket" {
  name                        = "${var.project_id}-function-${random_id.bucket_suffix.hex}"
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}

# 4. Upload the zip (using the hash of the file for the name)
resource "google_storage_bucket_object" "source_archive" {
  # This uses the hash of the actual content to trigger a redeploy when code changes
  name   = "source-${data.archive_file.source.output_md5}.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = data.archive_file.source.output_path
}

# 5. The Cloud Function
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
    available_memory      = "256Mi"
    service_account_email = var.service_account
    ingress_settings      = "ALLOW_ALL"
  }
}

# 6. Make the function public
resource "google_cloud_run_service_iam_member" "public_access" {
  location = var.region
  service  = google_cloudfunctions2_function.helloworld.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "function_url" {
  value = google_cloudfunctions2_function.helloworld.service_config[0].uri
}