terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.10.0"
    }
  }
  required_version = ">= 1.2.0"

  backend "gcs" {
    bucket = "nps-terraform-state"
    prefix = "polaris"
  }
}

provider "google" {
  region  = var.region
  project = var.project_id
}

data "google_project" "default" {}

