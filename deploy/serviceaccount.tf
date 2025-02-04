resource "google_service_account" "default" {
  account_id   = "polaris-svc"
  display_name = "Polaris Service Account"
}

