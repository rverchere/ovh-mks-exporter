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

type Instance struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IpAddresses []struct {
		Ip        string `json:"ip"`
		Type      string `json:"type"`
		Version   int    `json:"version"`
		NetworkId string `json:"networkId"`
		GatewayIp string `json:"gatewayIp"`
	} `json:"ipAddresses"`
	FlavorId      string    `json:"flavorId"`
	ImageId       string    `json:"imageId"`
	SshKeyId      string    `json:"sshKeyId"`
	CreatedAt     time.Time `json:"createdAt"`
	Region        string    `json:"region"`
	MonthyBilling struct {
		Since  string `json:"since"`
		Status string `json:"status"`
	} `json:"monthlyBilling,omitempty"`
	Status                      string   `json:"status"`
	PlanCode                    string   `json:"planCode"`
	OperationIds                []string `json:"operationIds,omitempty"`
	CurrentMonthOutgoingTraffic int64    `json:"currentMonthOutgoingTraffic,omitempty"`
}

type Node struct {
	Id          string    `json:"id"`
	ProjectId   string    `json:"projectId"`
	InstanceId  string    `json:"instanceId"`
	NodePoolId  string    `json:"nodePoolId"`
	Name        string    `json:"name"`
	Flavor      string    `json:"flavor"`
	Status      string    `json:"status"`
	IsUpToDate  bool      `json:"isUpToDate"`
	Version     string    `json:"Version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeployeddAt time.Time `json:"deployedAt"`
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
	Taints []struct {
		Effect string `json:"effect,omitempty"`
		Key    string `json:"key,omitempty"`
		Value  string `json:"value,omitempty"`
	} `json:"taints"`
	Unschedulable bool `json:"unschedulable"`
}

func GetClusterNodePool(client *ovh.Client, ServiceName string, KubeId string) []NodePool {

	log.Info(fmt.Sprintf("Getting cluster nodepools for cluster %s", KubeId))

	NodePoolUrl := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", ServiceName, KubeId)

	var res []NodePool

	err := client.Get(NodePoolUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res
}

func GetClusterNodePoolNode(client *ovh.Client, ServiceName string, KubeId string, NodepoolId string) []Node {

	log.Info(fmt.Sprintf("Getting cluster nodepool node for cluster %s, nodepool %s", KubeId, NodepoolId))

	NodePoolNodeUrl := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s/nodes", ServiceName, KubeId, NodepoolId)

	var res []Node

	err := client.Get(NodePoolNodeUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res

}

func GetClusterInstance(client *ovh.Client, ServiceName string, InstanceId string) Instance {

	log.Info(fmt.Sprintf("Getting cluster instance information %s", InstanceId))

	InstanceUrl := fmt.Sprintf("/cloud/project/%s/instance/%s", ServiceName, InstanceId)

	var res Instance

	err := client.Get(InstanceUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res

}

func GetClusterEtcdUsage(client *ovh.Client, ServiceName string, KubeId string) EtcdUsage {

	log.Info(fmt.Sprintf("Getting ETCD usage for cluster %s", KubeId))

	EtcdUsageUrl := fmt.Sprintf("/cloud/project/%s/kube/%s/metrics/etcdUsage", ServiceName, KubeId)

	var res EtcdUsage

	err := client.Get(EtcdUsageUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res
}

func GetClusterDescription(client *ovh.Client, ServiceName string, KubeId string) ClusterDescription {

	log.Info(fmt.Sprintf("Getting cluster description for cluster %s", KubeId))

	ClusterDescriptionUrl := fmt.Sprintf("/cloud/project/%s/kube/%s", ServiceName, KubeId)

	var res ClusterDescription
	err := client.Get(ClusterDescriptionUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res
}

func GetClusters(client *ovh.Client, ServiceName string) []string {

	log.Info(fmt.Sprintf("Getting clusters ID for service %s", ServiceName))

	ClustersUrl := fmt.Sprintf("/cloud/project/%s/kube", ServiceName)
	var res []string

	err := client.Get(ClustersUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res
}

func GetStorageContainers(client *ovh.Client, ServicName string) []StorageContainers {

	log.Info(fmt.Sprintf("Getting storage containers information for service %s", ServicName))

	StorageContainersUrl := fmt.Sprintf("/cloud/project/%s/storage", ServiceName)

	var res []StorageContainers

	err := client.Get(StorageContainersUrl, &res)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug(res)
	return res
}
