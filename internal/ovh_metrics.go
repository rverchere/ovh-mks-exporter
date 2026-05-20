package internal

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ovh/go-ovh/ovh"
)

type CloudProjectInformation struct {
	// https://transform.tools/json-to-go
	ProjectId    string    `json:"project_id"`
	ProjectName  string    `json:"projectName,omitempty"`
	Description  string    `json:"description,omitempty"`
	PlanCode     string    `json:"planCode"`
	Unleash      bool      `json:"unleash"`
	Expiration   time.Time `json:"expiration,omitempty"`
	CreationDate time.Time `json:"creationDate"`
	OrderId      int       `json:"orderId,omitempty"`
	Access       string    `json:"access"`
	Status       string    `json:"status"`
	ManualQuota  bool      `json:"manualQuota"`
}

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

type S3Container struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	ObjectsCount int64  `json:"objectsCount"`
	ObjectsSize  int64  `json:"objectsSize"`
	VirtualHost  string `json:"virtualHost"`
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
	MonthlyBilling struct {
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
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeployedAt  time.Time `json:"deployedAt"`
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

func GetCloudProjectInformation(client *ovh.Client, ServiceName string, maxRetries int) (CloudProjectInformation, error) {
	log.Info(fmt.Sprintf("Getting cloud project information for %s", ServiceName))
	url := fmt.Sprintf("/cloud/project/%s", ServiceName)
	return retry(maxRetries, func() (CloudProjectInformation, error) {
		var res CloudProjectInformation
		if err := client.Get(url, &res); err != nil {
			return res, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusterNodePool(client *ovh.Client, ServiceName string, KubeId string, maxRetries int) ([]NodePool, error) {
	log.Info(fmt.Sprintf("Getting cluster nodepools for cluster %s", KubeId))
	url := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", ServiceName, KubeId)
	return retry(maxRetries, func() ([]NodePool, error) {
		var res []NodePool
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusterNodePoolNode(client *ovh.Client, ServiceName string, KubeId string, NodepoolId string, maxRetries int) ([]Node, error) {
	log.Info(fmt.Sprintf("Getting cluster nodepool node for cluster %s, nodepool %s", KubeId, NodepoolId))
	url := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s/nodes", ServiceName, KubeId, NodepoolId)
	return retry(maxRetries, func() ([]Node, error) {
		var res []Node
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusterInstance(client *ovh.Client, ServiceName string, InstanceId string, maxRetries int) (Instance, error) {
	log.Info(fmt.Sprintf("Getting cluster instance information %s", InstanceId))
	url := fmt.Sprintf("/cloud/project/%s/instance/%s", ServiceName, InstanceId)
	return retry(maxRetries, func() (Instance, error) {
		var res Instance
		if err := client.Get(url, &res); err != nil {
			return res, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusterEtcdUsage(client *ovh.Client, ServiceName string, KubeId string, maxRetries int) (EtcdUsage, error) {
	log.Info(fmt.Sprintf("Getting ETCD usage for cluster %s", KubeId))
	url := fmt.Sprintf("/cloud/project/%s/kube/%s/metrics/etcdUsage", ServiceName, KubeId)
	return retry(maxRetries, func() (EtcdUsage, error) {
		var res EtcdUsage
		if err := client.Get(url, &res); err != nil {
			return res, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusterDescription(client *ovh.Client, ServiceName string, KubeId string, maxRetries int) (ClusterDescription, error) {
	log.Info(fmt.Sprintf("Getting cluster description for cluster %s", KubeId))
	url := fmt.Sprintf("/cloud/project/%s/kube/%s", ServiceName, KubeId)
	return retry(maxRetries, func() (ClusterDescription, error) {
		var res ClusterDescription
		if err := client.Get(url, &res); err != nil {
			return res, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetClusters(client *ovh.Client, ServiceName string, maxRetries int) ([]string, error) {
	log.Info(fmt.Sprintf("Getting clusters ID for service %s", ServiceName))
	url := fmt.Sprintf("/cloud/project/%s/kube", ServiceName)
	return retry(maxRetries, func() ([]string, error) {
		var res []string
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetStorageContainers(client *ovh.Client, ServiceName string, maxRetries int) ([]StorageContainers, error) {
	log.Info(fmt.Sprintf("Getting storage containers information for service %s", ServiceName))
	url := fmt.Sprintf("/cloud/project/%s/storage", ServiceName)
	return retry(maxRetries, func() ([]StorageContainers, error) {
		var res []StorageContainers
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetRegions(client *ovh.Client, ServiceName string, maxRetries int) ([]string, error) {
	log.Info(fmt.Sprintf("Getting regions for service %s", ServiceName))
	url := fmt.Sprintf("/cloud/project/%s/region", ServiceName)
	return retry(maxRetries, func() ([]string, error) {
		var res []string
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}

func GetS3Containers(client *ovh.Client, ServiceName string, RegionName string, maxRetries int) ([]S3Container, error) {
	log.Info(fmt.Sprintf("Getting S3 containers for service %s in region %s", ServiceName, RegionName))
	url := fmt.Sprintf("/cloud/project/%s/region/%s/storage", ServiceName, RegionName)
	return retry(maxRetries, func() ([]S3Container, error) {
		var res []S3Container
		if err := client.Get(url, &res); err != nil {
			return nil, err
		}
		log.Debug(res)
		return res, nil
	})
}
