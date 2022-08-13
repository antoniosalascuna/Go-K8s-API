package gcp_api

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"

	//	"k8s.io/client-go/1.5/tools/clientcmd"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/tools/clientcmd/api"
)

/* func ResApi(w http.ResponseWriter, r *http.Request) {

	var K8sData K8sData

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &K8sData); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err := GKEclusterIngress(K8sData.ProjectId, K8sData.GKElocation, K8sData.GKEclusterName, K8sData.GKENamespace)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

} */

//Funcion que ingresa al API-server de K8s cluster pasado por parametro para entrar en el desde aqui
/* func GKEclusterIngress(projectId string, gkeLocation string, gkeClusterName string, gkeNamespace string) error {

ctx := context.Background()





name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", creds.ProjectID, gkeLocation, gkeClusterName)

//Me conecto al cluster

cluster, err := container.NewProjectsLocationsClustersService(gkeService).Get(name).Do()

if err != nil {
	log.Fatalf("Failed to load GKE cluster %q: %s", name, err)
}

//Clientset es el encargado de tomar los credenciales del cluster ingresado para crear Deployments.

clientset, err := getGKEClientset(cluster, creds.TokenSource)

if err != nil {
	log.Fatalf("Failed to initialise Kubernetes clientset: %s", err)
}

//Obtiene los pods de NameSpace pasado por parametro

pods, err := clientset.CoreV1().Pods(gkeNamespace).List(ctx, metav1.ListOptions{})
if err != nil {
	log.Fatalf("Failed to list pods: %s", err)
}

log.Printf("There are %d pods in the namespace", len(pods.Items))

//kubeConfig, err := GcpAuth(ctx, projectId)

/* 	if err != nil {
   		return err
   	}

   	fmt.Println(kubeConfig) */

// Just list all the namespaces found in the project to test the API.
/*for clusterName := range kubeConfig.Clusters {
		cfg, err := clientcmd.NewNonInteractiveClientConfig(*kubeConfig, clusterName, &clientcmd.ConfigOverrides{CurrentContext: clusterName}, nil).ClientConfig()
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes configuration cluster=%s: %w", clusterName, err)
		}

		k8s, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client cluster=%s: %w", clusterName, err)
		}

		 ns,  err := k8s.CoreV1().Namespaces().List(ctx.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list namespaces cluster=%s: %w", clusterName, err)
		}

		log.Printf("Namespaces found in cluster=%s", clusterName)

		for _, item := range ns.Items {
			log.Println(item.Name)
		}
	}

	return nil

} */

func GcpAuth(ctx context.Context, projectId string) (*api.Config, error) {

	//stogeClient, err := storage.NewClient(cxt)

	client, err := container.NewService(ctx, option.WithCredentialsFile("././monkeyskubernetes-5ff80783703d.json"))

	if err != nil {
		log.Fatal(err)
	}

	ret := api.Config{
		APIVersion: "V1",
		Kind:       "Config",
		Clusters:   map[string]*api.Cluster{},
		AuthInfos:  map[string]*api.AuthInfo{},
		Contexts:   map[string]*api.Context{},
	}

	resp, err := client.Projects.Zones.Clusters.List("monkeyskubernetes", "-").Context(ctx).Do()

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range resp.Clusters {

		name := fmt.Sprintf("gke_%s_%s_%s", "monkeyskubernetes", f.Zone, f.Name)
		cert, err := base64.StdEncoding.DecodeString(f.MasterAuth.ClientCertificate)

		if err != nil {
			log.Fatal(err)
		}

		ret.Clusters[name] = &api.Cluster{
			CertificateAuthorityData: cert,
			Server:                   "https://" + f.Endpoint,
		}

		ret.Contexts[name] = &api.Context{
			Cluster:  name,
			AuthInfo: name,
		}

		ret.AuthInfos[name] = &api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "gcp",
				Config: map[string]string{
					"scopes": "https://www.googleapis.com/auth/cloud-platform",
				},
			},
		}

	}

	if err != nil {
		log.Fatal()
	}
	//it := stogeClient.Buckets(cxt, "monkeyskubernetes")

	/* 	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(bucketAttrs.Name)
	} */
	return &ret, nil

}
