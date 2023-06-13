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

	Clusters := fmt.Sprintf("/cloud/project/%s/kube", ServiceName)
	var res []string

	err := client.Get(Clusters, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Getting Clusters ID")
	return res

}
