package tools

import "google.golang.org/api/container/v1"

func ClusterRequest(NodePoolName, clusterName, machineType, DiskType string) *container.CreateClusterRequest {

	request := &container.CreateClusterRequest{
		Cluster: &container.Cluster{
			NodePools: []*container.NodePool{
				{

					Name:             NodePoolName,
					InitialNodeCount: 2,
					Management: &container.NodeManagement{
						AutoRepair:  true,
						AutoUpgrade: false,
					},
					MaxPodsConstraint: &container.MaxPodsConstraint{
						MaxPodsPerNode: 20,
					},
					Config: &container.NodeConfig{
						DiskSizeGb:  30,
						ImageType:   "COS_CONTAINERD",
						MachineType: machineType,
						DiskType:    DiskType,
						OauthScopes: []string{
							"https://www.googleapis.com/auth/devstorage.read_only",
							"https://www.googleapis.com/auth/logging.write",
							"https://www.googleapis.com/auth/monitoring",
							"https://www.googleapis.com/auth/servicecontrol",
							"https://www.googleapis.com/auth/service.management.readonly",
							"https://www.googleapis.com/auth/trace.append",
						},
					},
				},
			},
			Name:                  clusterName,
			Description:           "A cluster for e2e testing",
			EnableKubernetesAlpha: false,
		},
	}

	return request
}
