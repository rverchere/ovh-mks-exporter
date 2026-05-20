package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ovh/go-ovh/ovh"
)

const testService = "test-svc"

func newTestClient(t *testing.T, mux *http.ServeMux) *ovh.Client {
	t.Helper()
	mux.HandleFunc("/1.0/auth/time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(time.Now().Unix())
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	client, err := ovh.NewClient(srv.URL+"/1.0", "app-key", "app-secret", "consumer-key")
	if err != nil {
		t.Fatalf("ovh.NewClient: %v", err)
	}
	return client
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func TestGetCloudProjectInformation(t *testing.T) {
	want := CloudProjectInformation{
		ProjectId:   "proj-abc",
		ProjectName: "my-project",
		Description: "test",
		Status:      "ok",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetCloudProjectInformation(client, testService, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.ProjectId != want.ProjectId || got.Description != want.Description || got.Status != want.Status {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetClusters(t *testing.T) {
	want := []string{"cluster-1", "cluster-2"}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/kube", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusters(client, testService, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetClusterDescription(t *testing.T) {
	clusterID := "kube-1"
	want := ClusterDescription{
		ID:      clusterID,
		Name:    "my-cluster",
		Region:  "GRA5",
		Version: "1.29",
		Status:  "READY",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/kube/"+clusterID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusterDescription(client, testService, clusterID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != want.Name || got.Region != want.Region || got.Version != want.Version {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetClusterEtcdUsage(t *testing.T) {
	clusterID := "kube-1"
	want := EtcdUsage{Quota: 4096, Usage: 1024}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/kube/"+clusterID+"/metrics/etcdUsage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusterEtcdUsage(client, testService, clusterID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Quota != want.Quota || got.Usage != want.Usage {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetClusterNodePool(t *testing.T) {
	clusterID := "kube-1"
	want := []NodePool{
		{Id: "np-1", Name: "workers", CurrentNodes: 3, DesiredNodes: 3, Flavor: "b2-7"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/kube/"+clusterID+"/nodepool", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusterNodePool(client, testService, clusterID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Id != "np-1" || got[0].CurrentNodes != 3 {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetClusterNodePoolNode(t *testing.T) {
	clusterID := "kube-1"
	nodepoolID := "np-1"
	want := []Node{
		{Id: "node-1", InstanceId: "inst-1", Status: "READY"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/kube/"+clusterID+"/nodepool/"+nodepoolID+"/nodes", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusterNodePoolNode(client, testService, clusterID, nodepoolID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Id != "node-1" || got[0].InstanceId != "inst-1" {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetClusterInstance(t *testing.T) {
	instanceID := "inst-1"
	want := Instance{
		ID:     instanceID,
		Name:   "worker-1",
		Status: "ACTIVE",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/instance/"+instanceID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetClusterInstance(client, testService, instanceID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != want.Name || got.Status != want.Status {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetStorageContainers(t *testing.T) {
	want := []StorageContainers{
		{ID: "sc-1", Name: "container1", StoredObjects: 10, StoredBytes: 1024, Region: "GRA"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/storage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetStorageContainers(client, testService, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].StoredObjects != 10 || got[0].StoredBytes != 1024 {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetRegions(t *testing.T) {
	want := []string{"GRA1", "SBG1", "BHS1"}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/region", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetRegions(client, testService, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) || got[0] != "GRA1" {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetS3Containers(t *testing.T) {
	region := "GRA1"
	want := []S3Container{
		{Name: "s3-bucket", Region: region, ObjectsCount: 5, ObjectsSize: 2048},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/region/"+region+"/storage", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetS3Containers(client, testService, region, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "s3-bucket" || got[0].ObjectsCount != 5 {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetLoadBalancers(t *testing.T) {
	region := "GRA1"
	want := []LoadBalancer{
		{ID: "lb-1", Name: "my-lb", Region: region, OperatingStatus: "ONLINE", ProvisioningStatus: "ACTIVE", VipAddress: "1.2.3.4"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/region/"+region+"/loadbalancing/loadbalancer", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetLoadBalancers(client, testService, region, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "lb-1" || got[0].OperatingStatus != "ONLINE" {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetLoadBalancerStats(t *testing.T) {
	region := "GRA1"
	lbID := "lb-1"
	want := LoadBalancerStats{
		ActiveConnections: 10,
		BytesIn:           100,
		BytesOut:          200,
		RequestErrors:     1,
		TotalConnections:  500,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/1.0/cloud/project/"+testService+"/region/"+region+"/loadbalancing/loadbalancer/"+lbID+"/stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, want)
	})
	client := newTestClient(t, mux)

	got, err := GetLoadBalancerStats(client, testService, region, lbID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.ActiveConnections != want.ActiveConnections || got.BytesIn != want.BytesIn || got.TotalConnections != want.TotalConnections {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
