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

variable "fn_user_register" {
  type    = string
  default = ""
}
variable "fn_user_login" {
  type    = string
  default = ""
}

variable "fn_biz_register" {
  type    = string
  default = ""
}
variable "fn_add_categories" {
  type    = string
  default = ""
}
variable "fn_get_categories" {
  type    = string
  default = ""
}
variable "fn_create_location" {
  type    = string
  default = ""
}
variable "fn_new_product" {
  type    = string
  default = ""
}
variable "fn_create_inventory" {
  type    = string
  default = ""
}
variable "fn_get_business_products" {
  type    = string
  default = ""
}
variable "fn_get_business_inventory" {
  type    = string
  default = ""
}
variable "fn_create_spotlight" {
  type    = string
  default = ""
}
variable "fn_create_offer" {
  type    = string
  default = ""
}
variable "fn_get_business_offers" {
  type    = string
  default = ""
}
variable "fn_update_offer" {
  type    = string
  default = ""
}
variable "fn_delete_offer" {
  type    = string
  default = ""
}
variable "fn_nearby_products" {
  type    = string
  default = ""
}
variable "fn_search_products" {
  type    = string
  default = ""
}
variable "fn_spotlight_by_location" {
  type    = string
  default = ""
}
variable "fn_get_offers_by_business" {
  type    = string
  default = ""
}
variable "fn_like_business" {
  type    = string
  default = ""
}
variable "fn_dislike_business" {
  type    = string
  default = ""
}
variable "fn_get_liked_businesses" {
  type    = string
  default = ""
}
variable "fn_follow_business" {
  type    = string
  default = ""
}
variable "fn_get_following_businesses" {
  type    = string
  default = ""
}
variable "fn_followed_spotlights" {
  type    = string
  default = ""
}
variable "fn_claim_product" {
  type    = string
  default = ""
}
variable "fn_get_my_claims" {
  type    = string
  default = ""
}
variable "fn_cancel_claim" {
  type    = string
  default = ""
}
variable "fn_get_incoming_claims" {
  type    = string
  default = ""
}
variable "fn_accept_claim" {
  type    = string
  default = ""
}
variable "fn_reject_claim" {
  type    = string
  default = ""
}
variable "fn_mark_received" {
  type    = string
  default = ""
}
variable "fn_add_review" {
  type    = string
  default = ""
}
variable "fn_get_reviews" {
  type    = string
  default = ""
}
variable "fn_get_notifications" {
  type    = string
  default = ""
}
variable "fn_mark_notification_read" {
  type    = string
  default = ""
}
variable "fn_user_profile_info" {
  type    = string
  default = ""
}
variable "fn_biz_profile_info" {
  type    = string
  default = ""
}
variable "product_category_id" {
  type    = string
  default = ""
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

resource "google_cloudfunctions2_function" "function" {
  name     = var.function_name
  location = var.region

  build_config {
    runtime         = "go122"
    entry_point     = "Function"
    service_account = "projects/${var.project_id}/serviceAccounts/${var.service_account}"
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.source_archive.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    available_memory   = "512M"
    timeout_seconds    = 300
    service_account_email = var.service_account
    ingress_settings      = "ALLOW_ALL"

    environment_variables = {
      FN_USER_REGISTER          = var.fn_user_register
      FN_USER_LOGIN             = var.fn_user_login
      FN_BIZ_REGISTER           = var.fn_biz_register
      FN_ADD_CATEGORIES         = var.fn_add_categories
      FN_GET_CATEGORIES         = var.fn_get_categories
      FN_CREATE_LOCATION        = var.fn_create_location
      FN_NEW_PRODUCT            = var.fn_new_product
      FN_CREATE_INVENTORY       = var.fn_create_inventory
      FN_GET_BUSINESS_PRODUCTS  = var.fn_get_business_products
      FN_GET_BUSINESS_INVENTORY = var.fn_get_business_inventory
      FN_CREATE_SPOTLIGHT       = var.fn_create_spotlight
      FN_CREATE_OFFER           = var.fn_create_offer
      FN_GET_BUSINESS_OFFERS    = var.fn_get_business_offers
      FN_UPDATE_OFFER           = var.fn_update_offer
      FN_DELETE_OFFER           = var.fn_delete_offer
      FN_NEARBY_PRODUCTS        = var.fn_nearby_products
      FN_SEARCH_PRODUCTS        = var.fn_search_products
      FN_SPOTLIGHT_BY_LOCATION  = var.fn_spotlight_by_location
      FN_GET_OFFERS_BY_BUSINESS = var.fn_get_offers_by_business
      FN_LIKE_BUSINESS          = var.fn_like_business
      FN_DISLIKE_BUSINESS       = var.fn_dislike_business
      FN_GET_LIKED_BUSINESSES   = var.fn_get_liked_businesses
      FN_FOLLOW_BUSINESS        = var.fn_follow_business
      FN_GET_FOLLOWING_BUSINESSES = var.fn_get_following_businesses
      FN_FOLLOWED_SPOTLIGHTS    = var.fn_followed_spotlights
      FN_CLAIM_PRODUCT          = var.fn_claim_product
      FN_GET_MY_CLAIMS          = var.fn_get_my_claims
      FN_CANCEL_CLAIM           = var.fn_cancel_claim
      FN_GET_INCOMING_CLAIMS    = var.fn_get_incoming_claims
      FN_ACCEPT_CLAIM           = var.fn_accept_claim
      FN_REJECT_CLAIM           = var.fn_reject_claim
      FN_MARK_RECEIVED          = var.fn_mark_received
      FN_ADD_REVIEW             = var.fn_add_review
      FN_GET_REVIEWS            = var.fn_get_reviews
      FN_GET_NOTIFICATIONS      = var.fn_get_notifications
      FN_MARK_NOTIFICATION_READ = var.fn_mark_notification_read
      FN_USER_PROFILE_INFO      = var.fn_user_profile_info
      FN_BIZ_PROFILE_INFO       = var.fn_biz_profile_info
      PRODUCT_CATEGORY_ID       = var.product_category_id
    }
  }
}

resource "google_cloud_run_service_iam_member" "public_access" {
  location   = var.region
  service    = google_cloudfunctions2_function.function.service_config[0].service
  role       = "roles/run.invoker"
  member     = "allUsers"
  depends_on = [google_cloudfunctions2_function.function]
}

output "function_url" {
  value = google_cloudfunctions2_function.function.service_config[0].uri
}
