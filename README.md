# OVH Prometheus Exporter for Managed Kubernetes Clusters.

⚠️ This is a work in progress project, and I'm learning GO langage with that simple small project. ⚠️

This exporter retrieves 3 information from the OVHcloud API, which are not handled directly with k8s internal metrics:
- etcd quota usage
- up-to-date cluster version (to check if a security/patch upgrade is available)
- general information of the clusters

It retrieves metrics from all clusters defined in a public cloud project.

This exporter is inspired by the https://github.com/enix/x509-certificate-exporter project, thanks to Enix team!

## Prerequisites

You must generate OVHcloud API token here : https://eu.api.ovh.com/createToken/

The application uses only environment variables:
- OVH_ENDPOINT (default to `ovh-eu`)
- OVH_APPLICATION_KEY
- OVH_APPLICATION_SECRET
- OVH_CONSUMER_KEY
- OVH_CLOUDPROJECT_SERVICENAME: the service name of the OVHcloud Public Cloud Project


## Installation

On a kubernetes environment, you have to:

1. Create a secret `
```
kubectl create secret generic ovh-mks-exporter \
    --from-literal=OVH_ENDPOINT=ovh-eu \
    --from-literal=OVH_APPLICATION_KEY=${OVH_APPLICATION_KEY} \
    --from-literal=OVH_APPLICATION_SECRET=${OVH_APPLICATION_SECRET} \
    --from-literal=OVH_CONSUMER_KEY=${OVH_CONSUMER_KEY} \
    --from-literal=OVH_CLOUDPROJECT_SERVICENAME=${OVH_CLOUDPROJECT_SERVICENAME}
``` 
2. Deploy application, service and servicemonitor in the ̀deployment` folder (helm chart soon)
```
kubectl apply -f deployment/
```

The **servicemonitor** should be changed to match your prometheus installation (`release: prom` should differ).

## Metrics

3 metrics are exported :

| Name | Description | Values |
|------|-------------|--------|
| ovh_mks_cluster_isuptodate | Cluster is up to date (patch/security version) | 0 (no), 1 (yes) |
| ovh_mks_etcd_usage_quota_bytes | ETCD quota  max usage | bytes |
| ovh_mks_etcd_usage_usage_bytes | ETCD current usage | bytes |
| ovh_mks_cluster_info | Cluster Information (id, name, region, status, etc) | 1 |
## Example

![Grafana Dashboard](docs/grafana.png)


