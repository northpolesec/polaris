# Polaris

Polaris is the statistics collection server for the Santa client. It receives
requests over gRPC from the Santa client and publishes the data into BigQuery.

### Stats collection documentation

See https://northpole.dev/deployment/stats for more details

### Deployments

Terraform is used to manage the configuration of the artifact registry, Cloud
Run service, and BigQuery table. The terraform state is stored in a shared GCS
bucket.

Deployments of new versions are automated through GitHub with the container
image being created through GitHub actions, with attestation. This also triggers
an update to the Cloud Run deployment.
