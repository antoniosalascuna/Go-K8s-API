package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	tools "v1/tools"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Nsname Is DEV,TEST, LIVE Enviroment
func CreateMySQLSecret(location, clusterName, AppName, nsName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	secretyaml := &apiv1.Secret{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		Data: map[string][]byte{
			"password": []byte("1234"),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "mysql" + "-" + "pass",
			/*  Name: "mysql" + "-" + "pass", */
		},
	}

	//If nsName not exists set Default namespace
	if nsName == "" {

		secretClient := clientset.CoreV1().Secrets(metav1.NamespaceDefault)

		_, err := secretClient.Create(context.TODO(), secretyaml, metav1.CreateOptions{})

		return err

	} else {
		secretClient := clientset.CoreV1().Secrets(nsName + "-" + AppName)

		_, err := secretClient.Create(context.TODO(), secretyaml, metav1.CreateOptions{})

		return err

	}

}

func CreateSecretHandler(w http.ResponseWriter, r *http.Request) {

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

		err := CreateMySQLSecret(K8sData.Region, K8sData.ClusterName, K8sData.AppName, K8sData.EnvType)

		if err != nil {
			response.Message = "Secret created succesfully"
			response.Result = "Success"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

}
