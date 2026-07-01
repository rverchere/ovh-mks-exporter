# OVH Prometheus Exporter for Managed Kubernetes Clusters.

[![ci](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/docker-publish.yml)
[![ci](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/go-releaser.yml/badge.svg)](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/go-releaser.yml)
[![tests](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/go-test.yml/badge.svg)](https://github.com/rverchere/ovh-mks-exporter/actions/workflows/go-test.yml)
[![Plumber Score](https://score.getplumber.io/github.com/rverchere/ovh-mks-exporter.svg?style=flat)](https://score.getplumber.io/github.com/rverchere/ovh-mks-exporter)

[![Docker Stars](https://img.shields.io/docker/stars/rverchere/ovh-mks-exporter.svg?style=flat)](https://hub.docker.com/r/rverchere/ovh-mks-exporter/)
[![Docker Pulls](https://img.shields.io/docker/pulls/rverchere/ovh-mks-exporter.svg?style=flat)](https://hub.docker.com/r/rverchere/ovh-mks-exporter/)

This exporter retrieves some information from the OVHcloud API, which are not handled directly with k8s internal metrics:
- etcd quota usage
- up-to-date cluster version (to check if a security/patch upgrade is available)
- general information of the clusters
- instances information of the clusters
- swift storage (objects and usage)

It retrieves metrics from all clusters and swift containers defined in a public cloud project.

This exporter is inspired by the https://github.com/enix/x509-certificate-exporter project, thanks to Enix team!

## Prerequisites

You must generate OVHcloud API token here : https://eu.api.ovh.com/createToken/

The application uses only environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OVH_ENDPOINT` | No | `ovh-eu` | OVHcloud API endpoint |
| `OVH_APPLICATION_KEY` | Yes | | OVHcloud application key |
| `OVH_APPLICATION_SECRET` | Yes | | OVHcloud application secret |
| `OVH_CONSUMER_KEY` | Yes | | OVHcloud consumer key |
| `OVH_CLOUD_PROJECT_SERVICE` | Yes | | Service name of the OVHcloud Public Cloud Project |
| `OVH_LOG_LEVEL` | No | `info` | Log level: `trace`, `debug`, `info`, `warn`, `error` |
| `OVH_MAX_RETRIES` | No | `3` | Number of attempts on transient API errors (e.g. EOF). Each retry waits 1s, 2s, 4s… (exponential backoff) |
| `OVH_S3_REGIONS` | No | _(all regions)_ | Comma-separated list of regions to scrape for S3 metrics (e.g. `GRA,SBG,BHS`). If unset, all regions are queried |
| `OVH_LB_REGIONS` | No | _(all regions)_ | Comma-separated list of regions to scrape for Load Balancer metrics (e.g. `GRA,SBG`). If unset, all regions are queried |

## Installation

On a kubernetes environment, you have to:

1. Create a secret `
```
kubectl create secret generic ovh-mks-exporter \
    --from-literal=OVH_ENDPOINT=ovh-eu \
    --from-literal=OVH_APPLICATION_KEY=${OVH_APPLICATION_KEY} \
    --from-literal=OVH_APPLICATION_SECRET=${OVH_APPLICATION_SECRET} \
    --from-literal=OVH_CONSUMER_KEY=${OVH_CONSUMER_KEY} \
    --from-literal=OVH_CLOUD_PROJECT_SERVICE=${OVH_CLOUD_PROJECT_SERVICE}
``` 
2. Deploy application, service and servicemonitor using the helm chart in the `deployment`folder:
```
helm repo add ovh-mks-exporter https://rverchere.github.io/ovh-mks-exporter
helm repo update
helm upgrade --install ovh-mks-exporter ovh-mks-exporter/ovh-mks-exporter
```

The **servicemonitor** must be changed to match your prometheus installation (see `prometheusReleaseName`in the values file).

## Metrics

The following metrics are exported:

| Name | Description | Values |
|------|-------------|--------|
| ovh_mks_cloud_info | Public cloud projet information (id, name, description, status) | 1 |
| ovh_mks_cluster_info | Cluster information (id, name, region, status, etc) | 1 |
| ovh_mks_cluster_isuptodate | Cluster is up to date (patch/security version) | 0 (no), 1 (yes) |
| ovh_mks_etcd_usage_quota_bytes | ETCD quota  max usage | bytes |
| ovh_mks_etcd_usage_usage_bytes | ETCD current usage | bytes |
| ovh_mks_cluster_nodepool_info | Nodepool information (id, name, nodes number, nodes flavor, etc) | 1 |
| ovh_mks_cluster_instance_info | Instance information (id, name, billing) | 1 |
| ovh_storage_object_count | Swift container object count | count |
| ovh_storage_object_bytes | Swift container object usage | bytes |
| ovh_s3_object_count | S3 bucket object count | count |
| ovh_s3_object_bytes | S3 bucket object bytes | bytes |
| ovh_lb_info | Load Balancer information (id, name, region, operating_status, provisioning_status, vip_address) | 1 |
| ovh_lb_active_connections | Load Balancer active connections | count |
| ovh_lb_bytes_in_total | Load Balancer total bytes received | bytes |
| ovh_lb_bytes_out_total | Load Balancer total bytes sent | bytes |
| ovh_lb_request_errors_total | Load Balancer total request errors | count |
| ovh_lb_connections_total | Load Balancer total connections handled | count |
| ovh_mks_scrape_duration_seconds | Duration of the last metrics scrape | seconds |
| ovh_mks_exporter_build_info | Exporter build information (version, goversion, goos, goarch) | 1 |

## Example

![Grafana Dashboard](docs/grafana.png)


