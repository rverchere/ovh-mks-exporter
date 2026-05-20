package internal

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
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

	S3ContainerCountMetric = "ovh_s3_object_count"
	S3ContainerCountHelp   = "OVH S3 bucket object count"
	S3ContainerCountDesc   = prometheus.NewDesc(S3ContainerCountMetric, S3ContainerCountHelp, []string{"cloud_project_description", "name", "region"}, nil)

	S3ContainerUsageMetric = "ovh_s3_object_bytes"
	S3ContainerUsageHelp   = "OVH S3 bucket object bytes"
	S3ContainerUsageDesc   = prometheus.NewDesc(S3ContainerUsageMetric, S3ContainerUsageHelp, []string{"cloud_project_description", "name", "region"}, nil)

	LBInfoMetric = "ovh_lb_info"
	LBInfoHelp   = "OVH Load Balancer information"
	LBInfoDesc   = prometheus.NewDesc(LBInfoMetric, LBInfoHelp,
		[]string{"id", "name", "region", "operating_status", "provisioning_status", "vip_address"}, nil)

	LBActiveConnectionsMetric = "ovh_lb_active_connections"
	LBActiveConnectionsHelp   = "OVH Load Balancer active connections"
	LBActiveConnectionsDesc   = prometheus.NewDesc(LBActiveConnectionsMetric, LBActiveConnectionsHelp, []string{"id", "name", "region"}, nil)

	LBBytesInMetric = "ovh_lb_bytes_in_total"
	LBBytesInHelp   = "OVH Load Balancer total bytes received"
	LBBytesInDesc   = prometheus.NewDesc(LBBytesInMetric, LBBytesInHelp, []string{"id", "name", "region"}, nil)

	LBBytesOutMetric = "ovh_lb_bytes_out_total"
	LBBytesOutHelp   = "OVH Load Balancer total bytes sent"
	LBBytesOutDesc   = prometheus.NewDesc(LBBytesOutMetric, LBBytesOutHelp, []string{"id", "name", "region"}, nil)

	LBRequestErrorsMetric = "ovh_lb_request_errors_total"
	LBRequestErrorsHelp   = "OVH Load Balancer total request errors"
	LBRequestErrorsDesc   = prometheus.NewDesc(LBRequestErrorsMetric, LBRequestErrorsHelp, []string{"id", "name", "region"}, nil)

	LBTotalConnectionsMetric = "ovh_lb_connections_total"
	LBTotalConnectionsHelp   = "OVH Load Balancer total connections handled"
	LBTotalConnectionsDesc   = prometheus.NewDesc(LBTotalConnectionsMetric, LBTotalConnectionsHelp, []string{"id", "name", "region"}, nil)

	ScrapeDurationMetric = "ovh_mks_scrape_duration_seconds"
	ScrapeDurationHelp   = "Duration of the last OVH MKS metrics scrape in seconds"
	ScrapeDurationDesc   = prometheus.NewDesc(ScrapeDurationMetric, ScrapeDurationHelp, nil, nil)

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
	ch <- S3ContainerCountDesc
	ch <- S3ContainerUsageDesc
	ch <- LBInfoDesc
	ch <- LBActiveConnectionsDesc
	ch <- LBBytesInDesc
	ch <- LBBytesOutDesc
	ch <- LBRequestErrorsDesc
	ch <- LBTotalConnectionsDesc
	ch <- ScrapeDurationDesc
	ch <- InfoDesc
}

func (collector *collector) Collect(ch chan<- prometheus.Metric) {
	client := collector.exporter.Client
	serviceName := collector.exporter.ServiceName
	maxRetries := collector.exporter.MaxRetries

	start := time.Now()
	log.Info("Starting metrics collection")

	// Cloud project global information
	CloudProjectInformation, err := GetCloudProjectInformation(client, serviceName, maxRetries)
	if err != nil {
		log.Error("GetCloudProjectInformation: ", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		CloudProjectInfoDesc,
		prometheus.GaugeValue,
		float64(1),
		CloudProjectInformation.ProjectId, CloudProjectInformation.ProjectName, CloudProjectInformation.Description,
		CloudProjectInformation.Status,
	)

	// Kubernetes Managed Cluster information
	Clusters, err := GetClusters(client, serviceName, maxRetries)
	if err != nil {
		log.Error("GetClusters: ", err)
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Limit to 5 goroutines in parallel

	for _, KubeId := range Clusters {
		wg.Add(1)
		sem <- struct{}{} // block if 5 routines already there

		go func(KubeId string) {
			defer wg.Done()
			defer func() { <-sem }()

			EtcdUsage, err := GetClusterEtcdUsage(client, serviceName, KubeId, maxRetries)
			if err != nil {
				log.Error("GetClusterEtcdUsage: ", err)
				return
			}
			ClusterDescription, err := GetClusterDescription(client, serviceName, KubeId, maxRetries)
			if err != nil {
				log.Error("GetClusterDescription: ", err)
				return
			}

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

			ClusterNodepools, err := GetClusterNodePool(client, serviceName, KubeId, maxRetries)
			if err != nil {
				log.Error("GetClusterNodePool: ", err)
				return
			}
			for _, ClusterNodepool := range ClusterNodepools {
				ClusterNodePoolNodes, err := GetClusterNodePoolNode(client, serviceName, KubeId, ClusterNodepool.Id, maxRetries)
				if err != nil {
					log.Error("GetClusterNodePoolNode: ", err)
					continue
				}
				for _, ClusterNodePoolNode := range ClusterNodePoolNodes {
					ClusterInstance, err := GetClusterInstance(client, serviceName, ClusterNodePoolNode.InstanceId, maxRetries)
					if err != nil {
						log.Error("GetClusterInstance: ", err)
						continue
					}
					ch <- prometheus.MustNewConstMetric(
						ClusterInstanceInfoDesc,
						prometheus.GaugeValue,
						float64(1),
						KubeId, ClusterDescription.Region, ClusterDescription.Name, ClusterNodepool.Name,
						ClusterInstance.Name, ClusterInstance.Status, ClusterInstance.MonthlyBilling.Status,
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
	StorageContainers, err := GetStorageContainers(client, serviceName, maxRetries)
	if err != nil {
		log.Error("GetStorageContainers: ", err)
		return
	}
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

	// S3 Containers information
	s3Regions := collector.exporter.S3Regions
	if len(s3Regions) == 0 {
		s3Regions, err = GetRegions(client, serviceName, maxRetries)
		if err != nil {
			log.Error("GetRegions: ", err)
		}
	}
	if len(s3Regions) > 0 {
		var s3wg sync.WaitGroup
		s3sem := make(chan struct{}, 5)

		for _, region := range s3Regions {
			s3wg.Add(1)
			s3sem <- struct{}{}

			go func(region string) {
				defer s3wg.Done()
				defer func() { <-s3sem }()

				s3Containers, err := GetS3Containers(client, serviceName, region, maxRetries)
				if err != nil {
					log.Debugf("GetS3Containers %s: %v", region, err)
					return
				}
				for _, container := range s3Containers {
					ch <- prometheus.MustNewConstMetric(
						S3ContainerCountDesc,
						prometheus.GaugeValue,
						float64(container.ObjectsCount),
						CloudProjectInformation.Description, container.Name, container.Region,
					)
					ch <- prometheus.MustNewConstMetric(
						S3ContainerUsageDesc,
						prometheus.GaugeValue,
						float64(container.ObjectsSize),
						CloudProjectInformation.Description, container.Name, container.Region,
					)
				}
			}(region)
		}
		s3wg.Wait()
	}

	// Load Balancer information
	lbRegions := collector.exporter.LBRegions
	if len(lbRegions) == 0 {
		lbRegions, err = GetRegions(client, serviceName, maxRetries)
		if err != nil {
			log.Error("GetRegions: ", err)
		}
	}
	if len(lbRegions) > 0 {
		var lbwg sync.WaitGroup
		lbsem := make(chan struct{}, 5)

		for _, region := range lbRegions {
			lbwg.Add(1)
			lbsem <- struct{}{}

			go func(region string) {
				defer lbwg.Done()
				defer func() { <-lbsem }()

				lbs, err := GetLoadBalancers(client, serviceName, region, maxRetries)
				if err != nil {
					log.Debugf("GetLoadBalancers %s: %v", region, err)
					return
				}
				for _, lb := range lbs {
					ch <- prometheus.MustNewConstMetric(
						LBInfoDesc,
						prometheus.GaugeValue,
						float64(1),
						lb.ID, lb.Name, lb.Region, lb.OperatingStatus, lb.ProvisioningStatus, lb.VipAddress,
					)
					stats, err := GetLoadBalancerStats(client, serviceName, region, lb.ID, maxRetries)
					if err != nil {
						log.Errorf("GetLoadBalancerStats %s/%s: %v", region, lb.ID, err)
						continue
					}
					ch <- prometheus.MustNewConstMetric(LBActiveConnectionsDesc, prometheus.GaugeValue, float64(stats.ActiveConnections), lb.ID, lb.Name, lb.Region)
					ch <- prometheus.MustNewConstMetric(LBBytesInDesc, prometheus.GaugeValue, float64(stats.BytesIn), lb.ID, lb.Name, lb.Region)
					ch <- prometheus.MustNewConstMetric(LBBytesOutDesc, prometheus.GaugeValue, float64(stats.BytesOut), lb.ID, lb.Name, lb.Region)
					ch <- prometheus.MustNewConstMetric(LBRequestErrorsDesc, prometheus.GaugeValue, float64(stats.RequestErrors), lb.ID, lb.Name, lb.Region)
					ch <- prometheus.MustNewConstMetric(LBTotalConnectionsDesc, prometheus.GaugeValue, float64(stats.TotalConnections), lb.ID, lb.Name, lb.Region)
				}
			}(region)
		}
		lbwg.Wait()
	}

	// Application Information
	ch <- prometheus.MustNewConstMetric(
		InfoDesc,
		prometheus.GaugeValue,
		float64(1),
	)

	duration := time.Since(start).Seconds()
	ch <- prometheus.MustNewConstMetric(ScrapeDurationDesc, prometheus.GaugeValue, duration)
	log.Infof("Metrics collection completed in %.2fs", duration)
}
