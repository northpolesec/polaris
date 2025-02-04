resource "google_cloud_run_v2_service" "default" {
  project  = var.project_id
  name     = "polaris"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.default.email
    containers {
      name  = "polaris-1"
      image = "us-central1-docker.pkg.dev/${var.project_id}/polaris/polaris:latest"
      env {
        name  = "POLARIS_PROJECT_ID"
        value = google_bigquery_table.default.project
      }
      env {
        name  = "POLARIS_DATASET_ID"
        value = google_bigquery_table.default.dataset_id
      }
      env {
        name  = "POLARIS_TABLE_ID"
        value = google_bigquery_table.default.table_id
      }
      env {
        name  = "POLARIS_STREAM_ID"
        value = "_default"
      }
    }
  }
}

resource "google_cloud_run_domain_mapping" "default" {
  location = var.region
  name     = var.domain

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.default.name
  }
}

