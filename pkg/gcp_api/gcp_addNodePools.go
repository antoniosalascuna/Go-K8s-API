package gcp_api

import (
	"context"
	"fmt"
	tools "v1/tools"

	"google.golang.org/api/container/v1"
)

func AddNodePoolsIntoCluster(projectId, region, clusterName string) error {

	context := context.Background()

	creds, err := tools.GetCredentialFromJson(context, tools.JsonCredentialUrl)

	if err != nil {
		return err
	}

	gkeService, err := tools.GetGoogleService(creds, context)

	if err != nil {
		return err
	}

	noderequest := &container.CreateNodePoolRequest{
		NodePool: &container.NodePool{
			Name: "dev",
			Config: &container.NodeConfig{
				DiskSizeGb:  40,
				ImageType:   "COS",
				MachineType: "e2-medium",
				OauthScopes: []string{
					"https://www.googleapis.com/auth/compute",
					"https://www.googleapis.com/auth/devstorage.read_only",
					"https://www.googleapis.com/auth/logging.write",
					"https://www.googleapis.com/auth/monitoring.write",
					"https://www.googleapis.com/auth/servicecontrol",
					"https://www.googleapis.com/auth/service.management.readonly",
					"https://www.googleapis.com/auth/trace.append",
				},
			},
		},
		ForceSendFields: []string{},
		NullFields:      []string{},
	}

	createNodePool, err := gkeService.Projects.Locations.Clusters.NodePools.Create(tools.LocationNodeConfigRRN("us-central1-c", "monkeys"), noderequest).Context(context).Do()
	if err != nil {
		return err
	}

	fmt.Println(createNodePool)

	return err

}
