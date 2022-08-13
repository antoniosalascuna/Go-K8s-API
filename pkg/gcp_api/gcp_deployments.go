package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"v1/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NGXDeployment(projectId, location, clusterName, nsName, AppName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	deploymentclient := clientset.AppsV1().Deployments(nsName)

	WPDeploymentYAMLConfig := tools.DeploymentNGINXRequest(AppName)

	result, err := deploymentclient.Create(context.TODO(), WPDeploymentYAMLConfig, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	fmt.Println(result)

	return err
}

func CreateNGIXDeploymentHandler(w http.ResponseWriter, r *http.Request) {

	var K8sData tools.K8sData

	var response tools.Response

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

		err := NGXDeployment(tools.MONKEYSPROJECT, K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName)

		if err != nil {
			response.Message = "Deployment failed"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		response.Message = "Deployment created succesfully"
		response.Result = "Success"
		w.WriteHeader(http.StatusCreated)
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return

	}

}
