package gcp_api

import (
	"context"
	"fmt"
	"log"
	"time"

	tools "v1/tools"

	"google.golang.org/api/container/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func CreateCluster(region, clusterName, planType, envType, AppName string) error {

	/* var machinetype string

	SE UTILIZA CUANDO SE TIENE QUE TENER PRESENTE LOS PAQUETES A OFRECER


	if planType == "professional" {
		machinetype = "g1-small"
	} */

	context := context.Background()

	creds, err := tools.GetCredentialFromJson(context, tools.JsonCredentialUrl)

	if err != nil {
		return err
	}

	gkeService, err := tools.GetGoogleService(creds, context)

	if err != nil {
		return err
	}

	//Func to create the cluster

	clusterReq := tools.ClusterRequest("pool-001", clusterName, "g1-small", "pd-standard")

	result, err := gkeService.Projects.Locations.Clusters.Create(tools.LocationClusterZone(tools.MONKEYSPROJECT, region), clusterReq).Context(context).Do()

	if err != nil {
		return err
	}
	fmt.Println(result)

	time.Sleep(6 * time.Minute)

	/* 	if getStatusOfCluster(clusterName, region) == "RUNNING" {

		CreateNewNamespace(region, clusterName, envType, systemType)
	} */

	return err

}

func getStatusOfCluster(clusterName, region string) string {

	context := context.Background()

	creds, err := tools.GetCredentialFromJson(context, tools.JsonCredentialUrl)

	if err != nil {
		return "Error"
	}

	gkeService, err := tools.GetGoogleService(creds, context)

	if err != nil {
		return "Error"
	}

	clientset, err := tools.GetClientSetByJsonFile(region, clusterName)
	if err != nil {
		return err.Error()
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)

	informer := factory.Core().V1().Pods().Informer()

	stopper := make(chan struct{})

	defer close(stopper)
	const ByIP = "IndexByIP"

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)

			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
	})

	informer.Run(stopper)

	getCluster, err := gkeService.Projects.Locations.Clusters.Get(tools.LocationNodeConfigRRN(region, clusterName)).Context(context).Do()

	if err != nil {
		return err.Error()
	}

	return getCluster.Status
}

func UpdateNodeCountInPoolCluster(projectId, region, clusterName, nodePoolName string, nodeCount int) error {

	context := context.Background()

	creds, err := tools.GetCredentialFromJson(context, tools.JsonCredentialUrl)

	if err != nil {
		return err
	}

	gkeService, err := tools.GetGoogleService(creds, context)

	if err != nil {
		return err
	}

	noderequest := &container.SetNodePoolSizeRequest{
		Name:      tools.LocationNodePools(region, clusterName, nodePoolName),
		NodeCount: int64(nodeCount),
	}

	resp, err := gkeService.Projects.Locations.Clusters.NodePools.SetSize(tools.LocationNodePools(region, clusterName, nodePoolName), noderequest).Context(context).Do()

	fmt.Println(resp)
	fmt.Println(gkeService)

	if err != nil {
		return err
	}

	return err
}
