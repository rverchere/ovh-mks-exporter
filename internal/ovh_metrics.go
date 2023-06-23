package internal

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ovh/go-ovh/ovh"
)

var (
	// OVH related vars
	Client      *ovh.Client
	ServiceName string
	KubeId      string
)

type EtcdUsage struct {
	Quota int64 `json:"quota"`
	Usage int64 `json:"usage"`
}

type ClusterDescription struct {
	// https://transform.tools/json-to-go
	ID                  string   `json:"id"`
	Region              string   `json:"region"`
	Name                string   `json:"name"`
	URL                 string   `json:"url"`
	NodesURL            string   `json:"nodesUrl"`
	Version             string   `json:"version"`
	NextUpgradeVersions []string `json:"nextUpgradeVersions"`
	KubeProxyMode       string   `json:"kubeProxyMode"`
	Customization       struct {
		APIServer struct {
			AdmissionPlugins struct {
				Enabled  []string      `json:"enabled"`
				Disabled []interface{} `json:"disabled"`
			} `json:"admissionPlugins"`
		} `json:"apiServer"`
	} `json:"customization"`
	Status                 string      `json:"status"`
	UpdatePolicy           string      `json:"updatePolicy"`
	IsUpToDate             bool        `json:"isUpToDate"`
	ControlPlaneIsUpToDate bool        `json:"controlPlaneIsUpToDate"`
	PrivateNetworkID       interface{} `json:"privateNetworkId"`
	CreatedAt              time.Time   `json:"createdAt"`
	UpdatedAt              time.Time   `json:"updatedAt"`
}

type StorageContainers struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StoredObjects int    `json:"storedObjects"`
	StoredBytes   int    `json:"storedBytes"`
	Region        string `json:"region"`
}

type NodePool struct {
	Id             string                `json:"id"`
	ProjectId      string                `json:"projectId"`
	Name           string                `json:"name"`
	Autoscale      bool                  `json:"autoscale"`
	AntiAffinity   bool                  `json:"antiAffinity"`
	AvailableNodes int                   `json:"availableNodes"`
	CreatedAt      string                `json:"createdAt"`
	CurrentNodes   int                   `json:"currentNodes"`
	DesiredNodes   int                   `json:"desiredNodes"`
	Flavor         string                `json:"flavor"`
	MaxNodes       int                   `json:"maxNodes"`
	MinNodes       int                   `json:"minNodes"`
	MonthlyBilled  bool                  `json:"monthlyBilled"`
	SizeStatus     string                `json:"sizeStatus"`
	Status         string                `json:"status"`
	UpToDateNodes  int                   `json:"upToDateNodes"`
	UpdatedAt      string                `json:"updatedAt"`
	Template       *KubeNodePoolTemplate `json:"template,omitempty"`
}

type TaintEffectType int

type Taint struct {
	Effect TaintEffectType `json:"effect,omitempty"`
	Key    string          `json:"key,omitempty"`
	Value  string          `json:"value,omitempty"`
}

type KubeNodePoolTemplate struct {
	Metadata *KubeNodePoolTemplateMetadata `json:"metadata,omitempty"`
	Spec     *KubeNodePoolTemplateSpec     `json:"spec,omitempty"`
}

type KubeNodePoolTemplateMetadata struct {
	Annotations map[string]string `json:"annotations"`
	Finalizers  []string          `json:"finalizers"`
	Labels      map[string]string `json:"labels"`
}

type KubeNodePoolTemplateSpec struct {
	Taints        []Taint `json:"taints"`
	Unschedulable bool    `json:"unschedulable"`
}

func GetClusterNodePool(client *ovh.Client, ServiceName string, KubeId string) []NodePool {

	NodePoolUrl := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", ServiceName, KubeId)

	var res []NodePool

	err := client.Get(NodePoolUrl, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(res)
	return res
}

func GetClusterEtcdUsage(client *ovh.Client, ServiceName string, KubeId string) EtcdUsage {

	EtcdUsageUrl := fmt.Sprintf("/cloud/project/%s/kube/%s/metrics/etcdUsage", ServiceName, KubeId)

	var res EtcdUsage

	err := client.Get(EtcdUsageUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func GetClusterDescription(client *ovh.Client, ServiceName string, KubeId string) ClusterDescription {

	ClusterDescriptionUrl := fmt.Sprintf("/cloud/project/%s/kube/%s", ServiceName, KubeId)

	var res ClusterDescription
	err := client.Get(ClusterDescriptionUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func GetClusters(client *ovh.Client, ServiceName string) []string {

	ClustersUrl := fmt.Sprintf("/cloud/project/%s/kube", ServiceName)
	var res []string

	err := client.Get(ClustersUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Getting Clusters ID")
	return res

}

func GetStorageContainers(client *ovh.Client, ServicName string) []StorageContainers {

	StorageContainersUrl := fmt.Sprintf("/cloud/project/%s/storage", ServiceName)

	var res []StorageContainers

	err := client.Get(StorageContainersUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Getting Storage Containers information")

	return res

}
