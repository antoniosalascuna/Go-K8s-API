package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/container/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetClientSetByJsonFile(location, clusterName string) (kubernetes.Interface, error) {

	ctx := context.Background()

	urlCluster := LocationNodeConfigRRN(location, clusterName)

	creds, err := GetCredentialFromJson(ctx, JsonCredentialUrl)
	if err != nil {
		return nil, err
	}

	gkeService, err := GetGoogleService(creds, ctx)
	if err != nil {
		return nil, err
	}

	cluster, err := container.NewProjectsLocationsClustersService(gkeService).Get(urlCluster).Do()
	if err != nil {
		return nil, err
	}

	clientset, err := getGKEClientset(cluster, creds.TokenSource)

	if err != nil {
		return nil, err
	}

	return clientset, err
}

func getGKEClientset(cluster *container.Cluster, ts oauth2.TokenSource) (kubernetes.Interface, error) {

	capem, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)

	if err != nil {
		return nil, fmt.Errorf("failed to decode cluster CA cert: %s", err)
	}

	config := &rest.Config{
		Host: cluster.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: capem,
		},
	}

	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return &oauth2.Transport{
			Source: ts,
			Base:   rt,
		}
	})

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("failed to initialise clientset from config: %s", err)
	}

	return clientset, nil
}
