package internal

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type collector struct {
	exporter *Exporter
}

var (
	/* Prometheus related vars */
	EtcdUsageUsageMetric = "ovh_mks_etcd_usage_usage_bytes"
	EtcdUsageUsageHelp   = "OVH Managed Kubernetes Service ETCD Usage"
	EtcdUsageUsageDesc   = prometheus.NewDesc(EtcdUsageUsageMetric, EtcdUsageUsageHelp, nil, nil)

	EtcdUsageQuotaMetric = "ovh_mks_etcd_usage_quota_bytes"
	EtcdUsageQuotaHelp   = "OVH Managed Kubernetes Service ETCD Quota"
	EtcdUsageQuotaDesc   = prometheus.NewDesc(EtcdUsageQuotaMetric, EtcdUsageQuotaHelp, nil, nil)

	ClusterIsUpToDateMetric = "ovh_mks_cluster_isuptodate"
	ClusterIsUpToDateHelp   = "OVH Managed Kubernetes Service has a pending security/patch upgrade"
	ClusterIsUpToDateDesc   = prometheus.NewDesc(ClusterIsUpToDateMetric, ClusterIsUpToDateHelp, nil, nil)

	InfoMetric      = "ovh_mks_exporter_build_info"
	InfoHelp        = "A metric with a constant '1' value labeled with version, revision, build date, Go version, Go OS, and Go architecture"
	InfoConstLabels = prometheus.Labels{
		//"version":   Version,
		//"revision":  Revision,
		//"built":     BuildDateTime,
		"goversion": runtime.Version(),
		"goos":      runtime.GOOS,
		"goarch":    runtime.GOARCH,
	}
	InfoDesc = prometheus.NewDesc(InfoMetric, InfoHelp, nil, InfoConstLabels)
)

func Bool2int(b bool) int {
	// The compiler currently only optimizes this form.
	// See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

func (collector *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- EtcdUsageUsageDesc
	ch <- EtcdUsageQuotaDesc
	ch <- ClusterIsUpToDateDesc
	ch <- InfoDesc
}

func (collector *collector) Collect(ch chan<- prometheus.Metric) {
	EtcdUsage := GetClusterEtcdUsage(Client, ServiceName, KubeId)
	ClusterDescription := GetClusterDescription(Client, ServiceName, KubeId)

	ch <- prometheus.MustNewConstMetric(
		EtcdUsageUsageDesc,
		prometheus.GaugeValue,
		float64(EtcdUsage.Usage),
	)

	ch <- prometheus.MustNewConstMetric(
		EtcdUsageQuotaDesc,
		prometheus.GaugeValue,
		float64(EtcdUsage.Quota),
	)

	ch <- prometheus.MustNewConstMetric(
		ClusterIsUpToDateDesc,
		prometheus.GaugeValue,
		float64(Bool2int(ClusterDescription.IsUpToDate)),
	)

	ch <- prometheus.MustNewConstMetric(
		InfoDesc,
		prometheus.GaugeValue,
		float64(1),
	)
}
