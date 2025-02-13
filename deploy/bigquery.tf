resource "google_bigquery_dataset" "default" {
  dataset_id  = "polaris"
  description = "Dataset containing stats submitted from Santa clients"
}

resource "google_bigquery_table" "default" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "stats"

  time_partitioning {
    type = "DAY"
  }

  clustering = ["org_id", "santa_version"]

  schema = <<EOF
[
  {
    "name": "machine_id_hash",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "64"
  },
  {
    "name": "org_id",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "36"
  },
  {
    "name": "santa_version",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "16"
  },
  {
    "name": "macos_version",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "36"
  },
  {
    "name": "macos_build",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "36"
  },
  {
    "name": "mac_model",
    "type": "STRING",
    "mode": "NULLABLE",
    "maxLength": "36"
  }
]
EOF
}

resource "google_bigquery_table_iam_binding" "default" {
  project    = google_bigquery_table.default.project
  dataset_id = google_bigquery_table.default.dataset_id
  table_id   = google_bigquery_table.default.table_id
  role       = "roles/bigquery.dataOwner"
  members    = ["serviceAccount:${google_service_account.default.email}"]
}

