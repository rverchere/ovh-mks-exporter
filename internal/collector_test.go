package internal

import (
	"net/http"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
)

func TestBool2int(t *testing.T) {
	tests := []struct {
		in   bool
		want int
	}{
		{true, 1},
		{false, 0},
	}
	for _, tt := range tests {
		if got := Bool2int(tt.in); got != tt.want {
			t.Errorf("Bool2int(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestDescribe(t *testing.T) {
	exp := &Exporter{Client: nil, ServiceName: "svc", MaxRetries: 1}
	col := &collector{exporter: exp}

	ch := make(chan *prometheus.Desc, 30)
	col.Describe(ch)
	close(ch)

	var descs []*prometheus.Desc
	for d := range ch {
		descs = append(descs, d)
	}

	// Verify all 19 descriptors are sent
	if len(descs) != 19 {
		t.Errorf("Describe sent %d descriptors, want 19", len(descs))
	}
}

func TestCollect(t *testing.T) {
	const svc = "test-svc"
	const clusterID = "cluster-1"
	const nodepoolID = "np-1"
	const instanceID = "inst-1"
	const region = "GRA1"
	const lbID = "lb-1"

	mux := http.NewServeMux()

	// Cloud project
	mux.HandleFunc("/1.0/cloud/project/"+svc, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, CloudProjectInformation{
			ProjectId:   "proj-1",
			Description: "My Project",
			Status:      "ok",
		})
	})

	// Cluster list
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/kube", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []string{clusterID})
	})

	// Cluster description
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/kube/"+clusterID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, ClusterDescription{
			ID:                     clusterID,
			Name:                   "k8s-1",
			Region:                 "GRA5",
			Version:                "1.29",
			Status:                 "READY",
			UpdatePolicy:           "ALWAYS_UPDATE",
			IsUpToDate:             true,
			ControlPlaneIsUpToDate: true,
		})
	})

	// ETCD usage
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/kube/"+clusterID+"/metrics/etcdUsage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, EtcdUsage{Quota: 4096, Usage: 1024})
	})

	// Node pools
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/kube/"+clusterID+"/nodepool", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []NodePool{
			{Id: nodepoolID, Name: "workers", CurrentNodes: 2, DesiredNodes: 2, Flavor: "b2-7", MaxNodes: 5, MinNodes: 1},
		})
	})

	// Node pool nodes
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/kube/"+clusterID+"/nodepool/"+nodepoolID+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []Node{{Id: "node-1", InstanceId: instanceID, Status: "READY"}})
	})

	// Instance
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/instance/"+instanceID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, Instance{ID: instanceID, Name: "worker-1", Status: "ACTIVE"})
	})

	// Swift storage
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/storage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []StorageContainers{
			{ID: "sc-1", Name: "container1", StoredObjects: 5, StoredBytes: 1024, Region: "GRA"},
		})
	})

	// S3 storage (region-scoped)
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/region/"+region+"/storage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []S3Container{
			{Name: "bucket1", Region: region, ObjectsCount: 3, ObjectsSize: 512},
		})
	})

	// Load balancers
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/region/"+region+"/loadbalancing/loadbalancer", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, []LoadBalancer{
			{ID: lbID, Name: "my-lb", Region: region, OperatingStatus: "ONLINE", ProvisioningStatus: "ACTIVE", VipAddress: "1.2.3.4"},
		})
	})

	// Load balancer stats
	mux.HandleFunc("/1.0/cloud/project/"+svc+"/region/"+region+"/loadbalancing/loadbalancer/"+lbID+"/stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, LoadBalancerStats{ActiveConnections: 10, BytesIn: 100, BytesOut: 200, RequestErrors: 1, TotalConnections: 500})
	})

	client := newTestClient(t, mux)
	exp := &Exporter{
		Client:      client,
		ServiceName: svc,
		MaxRetries:  1,
		S3Regions:   []string{region},
		LBRegions:   []string{region},
	}
	col := &collector{exporter: exp}

	ch := make(chan prometheus.Metric, 50)
	col.Collect(ch)
	close(ch)

	var metrics []prometheus.Metric
	for m := range ch {
		metrics = append(metrics, m)
	}

	// 1 cloud info + 1 etcd usage + 1 etcd quota + 1 isuptodate + 1 instance info +
	// 1 nodepool info + 1 cluster info + 2 swift + 2 s3 + 1 lb info + 5 lb stats +
	// 1 build info + 1 scrape duration = 19
	if len(metrics) != 19 {
		t.Errorf("Collect produced %d metrics, want 19", len(metrics))
	}

	// Spot-check etcd usage value
	found := false
	for _, m := range metrics {
		var d dto.Metric
		if err := m.Write(&d); err != nil {
			t.Fatal(err)
		}
		if d.Gauge != nil && d.Gauge.GetValue() == 1024 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a gauge metric with value 1024 (etcd usage)")
	}
}
