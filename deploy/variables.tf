variable "project_id" {
  description = "The ID of the target GCP project."
  type        = string
  default     = "polaris-449516"
}

variable "region" {
  description = "The GCP region to host in."
  type        = string
  default     = "us-central1"
}

variable "domain" {
  description = "Domain name to run the service on."
  default     = "polaris.northpole.security"
}

variable "github_repo" {
  description = "GitHub repo hosting this repository, to give federated access."
  default     = "northpolesec/polaris"
}

