package main

import (
	"fmt"
	"log"
	"net/http"
	"v1/pkg/gcp_api"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	//"net/http"
	//"log"
)

func main() {

	r := mux.NewRouter()

	//Routes

	//r.HandleFunc("/gcpAuth", gcp_api.ResApi).Methods("POST")
	r.HandleFunc("/CreateSecret", gcp_api.CreateSecretHandler).Methods("POST")
	//r.HandleFunc("/CreateServiceMysql", gcp_api.CreateMySQLServiceHandler).Methods("POST")

	r.HandleFunc("/CreateCluster", gcp_api.CreateClusterHandler).Methods("POST")
	r.HandleFunc("/AddNodePools", gcp_api.AddNodePoolsIntoClusterHandler).Methods("POST")
	//Deployments
	r.HandleFunc("/CreateDeployment", gcp_api.CreateDeploymentHandler).Methods("POST")
	//Secrets

	r.HandleFunc("/CreateNameSpace", gcp_api.CreateCreateNewNamespaceHandler).Methods("POST")
	r.HandleFunc("/UpdateCountNodes", gcp_api.UpdateNodeCountClusterHandler).Methods("POST")
	r.HandleFunc("/NewService", gcp_api.CreateServiceHandler).Methods("POST")
	r.HandleFunc("/CreateServiceWP", gcp_api.CreateMyWPServiceHandler).Methods("POST")

	r.HandleFunc("/CreateIngress", gcp_api.CreateIngressHandler).Methods("POST")

	r.HandleFunc("/GetCluster", gcp_api.GetClusterHandler).Methods("GET")

	//r.HandleFunc("/CreateRecordSet", gcp_api.CreateRecordSetHandler).Methods("POST")

	handler := cors.Default().Handler(r)

	fmt.Println("Server on!!!!!!")

	//Cambiar siempre a 4001 por la config en el cluster
	log.Fatal(http.ListenAndServe(":4000", handler))
}
