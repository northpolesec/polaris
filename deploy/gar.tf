resource "google_artifact_registry_repository" "default" {
  location      = var.region
  repository_id = "polaris"
  format        = "DOCKER"

  docker_config {
    immutable_tags = false
  }

  cleanup_policies {
    id     = "KeepN-2"
    action = "KEEP"
    most_recent_versions {
      keep_count = 3
    }
  }
}

resource "google_artifact_registry_repository_iam_member" "default" {
  project    = google_artifact_registry_repository.default.project
  location   = google_artifact_registry_repository.default.location
  repository = google_artifact_registry_repository.default.name
  role       = "roles/artifactregistry.writer"
  member     = "principal://iam.googleapis.com/projects/${data.google_project.default.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.default.workload_identity_pool_id}/subject/repo:${var.github_repo}:ref:refs/heads/main"
}
