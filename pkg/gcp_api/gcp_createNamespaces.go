package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	tools "v1/tools"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNewNamespace(location, clusterName, envType, AppName string) (*v1.Namespace, error) {

	context := context.Background()

	var Ns = envType + "-" + AppName

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return nil, err
	}

	nsName := &v1.Namespace{

		ObjectMeta: metav1.ObjectMeta{
			Name: envType + "-" + AppName,
		},
	}

	getNs, err := clientset.CoreV1().Namespaces().Get(context, Ns, metav1.GetOptions{})

	if getNs.Status.Phase == "Active" {

		return getNs, err

	} else {

		NsResult, err := clientset.CoreV1().Namespaces().Create(context, nsName, metav1.CreateOptions{})

		return NsResult, err

	}
}

func CreateCreateNewNamespaceHandler(w http.ResponseWriter, r *http.Request) {

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

		_, err := CreateNewNamespace(K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName)

		if err != nil {
			response.Message = "Namespaces fail to create"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		response.Message = "Namespaces created succesfully"
		response.Result = "Success"
		w.WriteHeader(http.StatusOK)
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

}
