package internal

import (
	"runtime"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	exporter *Exporter
}

var (
	/* Prometheus related vars */

	CloudProjectInfoMetric = "ovh_mks_cloud_info"
	CloudProjectInfoHelp   = "OVH Public Cloud Information"
	CloudProjectInfoDesc   = prometheus.NewDesc(CloudProjectInfoMetric, CloudProjectInfoHelp,
		[]string{"id", "name", "description", "status"}, nil)

	EtcdUsageUsageMetric = "ovh_mks_etcd_usage_usage_bytes"
	EtcdUsageUsageHelp   = "OVH Managed Kubernetes Service ETCD Usage"
	EtcdUsageUsageDesc   = prometheus.NewDesc(EtcdUsageUsageMetric, EtcdUsageUsageHelp, []string{"id", "region", "name", "version"}, nil)

	EtcdUsageQuotaMetric = "ovh_mks_etcd_usage_quota_bytes"
	EtcdUsageQuotaHelp   = "OVH Managed Kubernetes Service ETCD Quota"
	EtcdUsageQuotaDesc   = prometheus.NewDesc(EtcdUsageQuotaMetric, EtcdUsageQuotaHelp, []string{"id", "region", "name", "version"}, nil)

	ClusterIsUpToDateMetric = "ovh_mks_cluster_isuptodate"
	ClusterIsUpToDateHelp   = "OVH Managed Kubernetes Service has a pending security/patch upgrade"
	ClusterIsUpToDateDesc   = prometheus.NewDesc(ClusterIsUpToDateMetric, ClusterIsUpToDateHelp, []string{"id", "region", "name", "version"}, nil)

	ClusterInfoMetric = "ovh_mks_cluster_info"
	ClusterInfoHelp   = "OVH Managed Kubernetes Service Information"
	ClusterInfoDesc   = prometheus.NewDesc(ClusterInfoMetric, ClusterInfoHelp,
		[]string{"id", "region", "name", "version", "status", "update_policy", "is_up_to_date", "control_plane_is_up_to_date"}, nil)

	ClusterNodepoolInfoMetric = "ovh_mks_cluster_nodepool_info"
	ClusterNodepoolInfoHelp   = "OVH Managed Kubernetes Nodepool information"
	ClusterNodepoolInfoDesc   = prometheus.NewDesc(ClusterNodepoolInfoMetric, ClusterNodepoolInfoHelp,
		[]string{"id", "region", "name", "version", "nodepool_name", "current_nodes", "desired_nodes", "flavor", "max_nodes", "min_nodes", "monthly_billed", "status"}, nil)

	ClusterInstanceInfoMetric = "ovh_mks_cluster_instance_info"
	ClusterInstanceInfoHelp   = "OVH Managed Kubernetes Instances information"
	ClusterInstanceInfoDesc   = prometheus.NewDesc(ClusterInstanceInfoMetric, ClusterInstanceInfoHelp,
		[]string{"id", "region", "name", "nodepool_name", "node_name", "status", "monthly_billed"}, nil)

	StorageContainerCountMetric = "ovh_storage_object_count"
	StorageContainerCountHelp   = "OVH storage containers object count"
	StorageContainerCountDesc   = prometheus.NewDesc(StorageContainerCountMetric, StorageContainerCountHelp, []string{"cloud_project_description", "id", "region", "name"}, nil)

	StorageContainerUsageMetric = "ovh_storage_object_bytes"
	StorageContainerUsageHelp   = "OVH storage containers object bytes"
	StorageContainerUsageDesc   = prometheus.NewDesc(StorageContainerUsageMetric, StorageContainerUsageHelp, []string{"cloud_project_description", "id", "region", "name"}, nil)

	InfoMetric      = "ovh_mks_exporter_build_info"
	InfoHelp        = "A metric with a constant '1' value labeled with version, revision, build date, Go version, Go OS, and Go architecture"
	InfoConstLabels = prometheus.Labels{
		"version":   Version,
		"goversion": runtime.Version(),
		"goos":      runtime.GOOS,
		"goarch":    runtime.GOARCH,
	}
	InfoDesc = prometheus.NewDesc(InfoMetric, InfoHelp, nil, InfoConstLabels)
)

func Bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (collector *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- CloudProjectInfoDesc
	ch <- EtcdUsageUsageDesc
	ch <- EtcdUsageQuotaDesc
	ch <- ClusterIsUpToDateDesc
	ch <- ClusterInfoDesc
	ch <- ClusterNodepoolInfoDesc
	ch <- ClusterInstanceInfoDesc
	ch <- StorageContainerCountDesc
	ch <- StorageContainerUsageDesc
	ch <- InfoDesc
}

func (collector *collector) Collect(ch chan<- prometheus.Metric) {
	// Cloud project global information
	CloudProjectInformation := GetCloudProjectInformation(Client, ServiceName)
	ch <- prometheus.MustNewConstMetric(
		CloudProjectInfoDesc,
		prometheus.GaugeValue,
		float64(1),
		CloudProjectInformation.ProjectId, CloudProjectInformation.ProjectName, CloudProjectInformation.Description,
		CloudProjectInformation.Status,
	)

	// Kubernetes Managed Cluster information
	Clusters := GetClusters(Client, ServiceName)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Limit to 5 goroutines in parallel

	for _, KubeId := range Clusters {
		wg.Add(1)
		sem <- struct{}{} // block if 5 routines already there

		go func(KubeId string) {
			defer wg.Done()
			defer func() { <-sem }() // libère une place

			EtcdUsage := GetClusterEtcdUsage(Client, ServiceName, KubeId)
			ClusterDescription := GetClusterDescription(Client, ServiceName, KubeId)

			ch <- prometheus.MustNewConstMetric(
				EtcdUsageUsageDesc,
				prometheus.GaugeValue,
				float64(EtcdUsage.Usage),
				KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterDescription.Version,
			)
			ch <- prometheus.MustNewConstMetric(
				EtcdUsageQuotaDesc,
				prometheus.GaugeValue,
				float64(EtcdUsage.Quota),
				KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterDescription.Version,
			)
			ch <- prometheus.MustNewConstMetric(
				ClusterIsUpToDateDesc,
				prometheus.GaugeValue,
				float64(Bool2int(ClusterDescription.IsUpToDate)),
				KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterDescription.Version,
			)

			ClusterNodepools := GetClusterNodePool(Client, ServiceName, KubeId)
			for _, ClusterNodepool := range ClusterNodepools {
				ClusterNodePoolNodes := GetClusterNodePoolNode(Client, ServiceName, KubeId, ClusterNodepool.Id)
				for _, ClusterNodePoolNode := range ClusterNodePoolNodes {
					ClusterInstance := GetClusterInstance(Client, ServiceName, ClusterNodePoolNode.InstanceId)
					ch <- prometheus.MustNewConstMetric(
						ClusterInstanceInfoDesc,
						prometheus.GaugeValue,
						float64(1),
						KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterNodepool.Name,
						ClusterInstance.Name, ClusterInstance.Status, ClusterInstance.MonthyBilling.Status,
					)
				}
				ch <- prometheus.MustNewConstMetric(
					ClusterNodepoolInfoDesc,
					prometheus.GaugeValue,
					float64(1),
					KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterDescription.Version,
					ClusterNodepool.Name, strconv.Itoa(ClusterNodepool.CurrentNodes), strconv.Itoa(ClusterNodepool.DesiredNodes),
					ClusterNodepool.Flavor, strconv.Itoa(ClusterNodepool.MaxNodes), strconv.Itoa(ClusterNodepool.MinNodes),
					strconv.FormatBool(ClusterNodepool.MonthlyBilled), ClusterNodepool.Status,
				)
			}
			ch <- prometheus.MustNewConstMetric(
				ClusterInfoDesc,
				prometheus.GaugeValue,
				float64(1),
				KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterDescription.Version,
				ClusterDescription.Status, ClusterDescription.UpdatePolicy,
				strconv.FormatBool(ClusterDescription.IsUpToDate), strconv.FormatBool(ClusterDescription.ControlPlaneIsUpToDate),
			)
		}(KubeId)
	}
	wg.Wait()

	// Storage Containers (Swift) information
	StorageContainers := GetStorageContainers(Client, ServiceName)
	for _, StorageContainer := range StorageContainers {
		ch <- prometheus.MustNewConstMetric(
			StorageContainerCountDesc,
			prometheus.GaugeValue,
			float64(StorageContainer.StoredObjects),
			CloudProjectInformation.Description, StorageContainer.ID, StorageContainer.Region, StorageContainer.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			StorageContainerUsageDesc,
			prometheus.GaugeValue,
			float64(StorageContainer.StoredBytes),
			CloudProjectInformation.Description, StorageContainer.ID, StorageContainer.Region, StorageContainer.Name,
		)
	}

	// Application Information
	ch <- prometheus.MustNewConstMetric(
		InfoDesc,
		prometheus.GaugeValue,
		float64(1),
	)
}
