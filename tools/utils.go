package tools

import "fmt"

type Response struct {
	Message string `json:"message"`
	Result  string `json:"result"`
	Error   string `json:"error"`
}

/*
--> Region = this value is region of cluster
--> ClusterName = this value is the Name of cluster
--> PlanType = this value has reference type of plan (Basic, professional, etc)
--> EnvType = this value are of type enviroment( into cluster is the namespace) Example = dev, test, live
--> AppName = this value has reference to which is the AppName(Wordpress, Python, Drupal, NGINX)
--> DNSname = Domain name
*/
type K8sData struct {
	Region      string `json:"Region,omitempty"`
	ClusterName string `json:"ClusterName,omitempty"`
	PlanType    string `json:"PlanType,omitempty"`
	EnvType     string `json:"envType,omitempty"`
	AppName     string `json:"AppName,omitempty"`
	DnsName     string `json:"DnsName,omitempty"`
}

type ClusterConfig struct {
	ClusterName         string `json:"ClusterName,omitempty"`
	ClusterNodeCount    int    `json:"ClusterNodeCount,omitempty"`
	ClusterNodePoolName string `json:"ClusterNodePoolName,omitempty"`
	ClusterRegion       string `json:"ClusterRegion,omitempty"`
	PVCName             string `json:"PVCName,omitempty"`
	NsName              string `json:"NsName,omitempty"`
	AppName             string `json:"AppName,omitempty"`
}

type ServiceConfig struct {
	AppName     string `json:"AppName,omitempty"`
	NsName      string `json:"NsName,omitempty"`
	ClusterName string `json:"ClusterName,omitempty"`
	Region      string `json:"Region,omitempty"`
	Port        int    `json:"Port,omitempty"`
}

//ENV VARIABLES
const MONKEYSPROJECT = "monkeyscloudclients"
const JsonCredentialUrl = "././monkeyscloudclients-020ae8c51c7b.json"
const DNSCLIENTZONE = "monkeys-clients-zone"
const DNSNAME = ".monkeys.cloud."

func LocationRRN(location string) string {
	return fmt.Sprintf("projects/%s/locations/%s", MONKEYSPROJECT, location)
}
func LocationNodeConfigRRN(location, cluster string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", MONKEYSPROJECT, location, cluster)
}

func LocationClusterZone(projectID, zone string) string {
	return fmt.Sprintf("projects/%s/zones/%s", projectID, zone)
}

func LocationNodePools(zone, clusterName, nodePoolName string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", MONKEYSPROJECT, zone, clusterName, nodePoolName)
}

func Int32Ptr(i int32) *int32 { return &i }
