package internal

import (
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	exporter *Exporter
}

var (
	/* Prometheus related vars */
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
	ClusterInfoHelp   = "OVH Managed Kubernetes Service Informations"
	ClusterInfoDesc   = prometheus.NewDesc(ClusterInfoMetric, ClusterInfoHelp,
		[]string{"id", "region", "name", "version", "status", "update_policy", "is_up_to_date", "control_plane_is_up_to_date"}, nil)

	ClusterNodepoolInfoMetric = "ovh_mks_cluster_nodepool_info"
	ClusterNodepoolInfoHelp   = "OVH Managed Kubernetes Nodepool informations"
	ClusterNodepoolInfoDesc   = prometheus.NewDesc(ClusterNodepoolInfoMetric, ClusterNodepoolInfoHelp,
		[]string{"id", "region", "name", "version", "nodepool_name", "current_nodes", "desired_nodes", "flavor", "max_nodes", "min_nodes", "monthly_billed", "status"}, nil)

	StorageContainerCountMetric = "ovh_storage_object_count"
	StorageContainerCountHelp   = "OVH storage containers object count"
	StorageContainerCountDesc   = prometheus.NewDesc(StorageContainerCountMetric, StorageContainerCountHelp, []string{"id", "region", "name"}, nil)

	StorageContainerUsageMetric = "ovh_storage_object_bytes"
	StorageContainerUsageHelp   = "OVH storage containers object bytes"
	StorageContainerUsageDesc   = prometheus.NewDesc(StorageContainerUsageMetric, StorageContainerUsageHelp, []string{"id", "region", "name"}, nil)

	InfoMetric      = "ovh_mks_exporter_build_info"
	InfoHelp        = "A metric with a constant '1' value labeled with version, revision, build date, Go version, Go OS, and Go architecture"
	InfoConstLabels = prometheus.Labels{
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
	ch <- EtcdUsageUsageDesc
	ch <- EtcdUsageQuotaDesc
	ch <- ClusterIsUpToDateDesc
	ch <- ClusterInfoDesc
	ch <- ClusterNodepoolInfoDesc
	ch <- StorageContainerCountDesc
	ch <- StorageContainerUsageDesc
	ch <- InfoDesc
}

func (collector *collector) Collect(ch chan<- prometheus.Metric) {

	/* Kubernetes Managed Cluster information */
	var Clusters []string = GetClusters(Client, ServiceName)

	for _, KubeId := range Clusters {
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

		var ClusterNodepools []NodePool = GetClusterNodePool(Client, ServiceName, KubeId)
		for _, ClusterNodepool := range ClusterNodepools {
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
			strconv.FormatBool(lusterDescription.IsUpToDate), strconv.FormatBool(ClusterDescription.ControlPlaneIsUpToDate),
		)
	}

	/* Storage Containers (Swift) information */
	var StorageContainers []StorageContainers = GetStorageContainers(Client, ServiceName)

	for _, StorageContainer := range StorageContainers {
		ch <- prometheus.MustNewConstMetric(
			StorageContainerCountDesc,
			prometheus.GaugeValue,
			float64(StorageContainer.StoredObjects),
			StorageContainer.ID, StorageContainer.Region, StorageContainer.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			StorageContainerUsageDesc,
			prometheus.GaugeValue,
			float64(StorageContainer.StoredBytes),
			StorageContainer.ID, StorageContainer.Region, StorageContainer.Name,
		)

	}

	/* Application Information */
	ch <- prometheus.MustNewConstMetric(
		InfoDesc,
		prometheus.GaugeValue,
		float64(1),
	)
}
