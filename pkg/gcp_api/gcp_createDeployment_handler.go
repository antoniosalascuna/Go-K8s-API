package gcp_api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	tools "v1/tools"
)

func CreateDeploymentHandler(w http.ResponseWriter, r *http.Request) {

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

		result, err := CreateNewNamespace(K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName)

		if err != nil {
			response.Message = "Namespace failed"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return

		}

		time.Sleep(4 * time.Second)

		//StateFulApps

		//Comprueba si la ns esta listo para utilizarse
		if result.Status.Phase == "Active" {

			if strings.ToLower(K8sData.AppName) == "drupal" || strings.ToLower(K8sData.AppName) == "wordpress" {

				secretStatus := CreateMySQLSecret(K8sData.Region, K8sData.ClusterName, K8sData.AppName, K8sData.EnvType)

				if secretStatus != nil {
					response.Message = "Error when secret has been created"
					response.Result = "Error"
					response.Error = secretStatus.Error()
					w.WriteHeader(http.StatusBadRequest)
					encodeData, _ := json.Marshal(response)
					fmt.Fprintf(w, string(encodeData))
					return

				}

				serviceStatus := CreateService(K8sData.AppName, K8sData.EnvType, K8sData.ClusterName, K8sData.Region)

				if serviceStatus != nil {
					response.Message = "Error when Service has been created"
					response.Result = "Error"
					response.Error = serviceStatus.Error()
					w.WriteHeader(http.StatusBadRequest)
					encodeData, _ := json.Marshal(response)
					fmt.Fprintf(w, string(encodeData))
					return

				}

				DeploymentStatus := CreateDeployment(K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName)

				if DeploymentStatus != nil {
					response.Message = "Error when deployment has been created"
					response.Result = "Error"
					response.Error = DeploymentStatus.Error()
					w.WriteHeader(http.StatusBadRequest)
					encodeData, _ := json.Marshal(response)
					fmt.Fprintf(w, string(encodeData))
					return

				}

				response.Message = "Deployment created succesfully"
				response.Result = "Success"
				w.WriteHeader(http.StatusCreated)
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))
				return
				//CreateMYSQLDeployment(tools.MONKEYSPROJECT, K8sData.GKElocation, K8sData.GKEclusterName, K8sData.EnvType, K8sData.AppName)
				//StatelssApps
			} else if strings.ToLower(K8sData.AppName) == "nginx" {

				serviceStatus := CreateService(K8sData.AppName, K8sData.EnvType, K8sData.ClusterName, K8sData.Region)

				if serviceStatus != nil {
					response.Message = "Error when Service has been created"
					response.Result = "Error"
					response.Error = serviceStatus.Error()
					w.WriteHeader(http.StatusBadRequest)
					encodeData, _ := json.Marshal(response)
					fmt.Fprintf(w, string(encodeData))
					return

				}

				DeploymentStatus := CreateDeployment(K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName)

				if DeploymentStatus != nil {
					response.Message = "Error when deployment has been created"
					response.Result = "Error"
					response.Error = DeploymentStatus.Error()
					w.WriteHeader(http.StatusBadRequest)
					encodeData, _ := json.Marshal(response)
					fmt.Fprintf(w, string(encodeData))
					return

				}

				response.Message = "Deployment created succesfully"
				response.Result = "Success"
				w.WriteHeader(http.StatusCreated)
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))

			}

		}

	}

}
