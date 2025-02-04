resource "google_iam_workload_identity_pool" "default" {
  workload_identity_pool_id = "github"
  display_name              = "GitHub"
}

resource "google_iam_workload_identity_pool_provider" "default" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.default.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
  display_name                       = "GitHub: ${var.github_repo}"
  description                        = "GitHub Actions identity pool provider for: ${var.github_repo}"
  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }
  attribute_condition = "google.subject == repo:${var.github_repo}:ref:refs/heads/main"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}
