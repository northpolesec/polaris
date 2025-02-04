resource "google_service_account" "default" {
  account_id   = "polaris-svc"
  display_name = "Polaris Service Account"
}

resource "google_service_account_iam_member" "default" {
  service_account_id = google_service_account.default.name
  role               = "roles/iam.serviceAccountUser"
  member             = "principal://iam.googleapis.com/projects/${data.google_project.default.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.default.workload_identity_pool_id}/subject/repo:${var.github_repo}:ref:refs/heads/main"
}

