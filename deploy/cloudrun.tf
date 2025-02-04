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

  lifecycle {
    # Images are updated by the CI/CD pipeline.
    ignore_changes = [
      template.0.containers.0.image,
      template.0.labels,
    ]
  }
}

resource "google_cloud_run_v2_service_iam_member" "public_access" {
  project  = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_v2_service_iam_member" "federated_run_admin" {
  project  = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name
  role     = "roles/run.admin"
  member   = "principal://iam.googleapis.com/projects/${data.google_project.default.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.default.workload_identity_pool_id}/subject/repo:${var.github_repo}:ref:refs/heads/main"
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

