package gcp_api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"v1/tools"
)

func CreateClusterHandler(w http.ResponseWriter, r *http.Request) {

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

		err := CreateCluster(K8sData.Region, K8sData.ClusterName, K8sData.PlanType, K8sData.EnvType, K8sData.AppName)

		if err != nil {
			response.Message = "Cluster fail to create"
			response.Result = "Error"
			response.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		response.Message = "Cluster created succesfully"
		response.Result = "Success"
		w.WriteHeader(http.StatusOK)
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

}

func GetClusterHandler(w http.ResponseWriter, r *http.Request) {

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

		Status := getStatusOfCluster(K8sData.ClusterName, K8sData.Region)

		/* 	if err != nil {
			response.Message = "Cluster fail to create"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		} */

		response.Message = "Cluster created succesfully"
		response.Result = "Success"
		w.WriteHeader(http.StatusOK)
		encodeData, _ := json.Marshal(Status)
		fmt.Fprintf(w, string(encodeData))
		return
	}

}

func UpdateNodeCountClusterHandler(w http.ResponseWriter, r *http.Request) {

	var ClusterData tools.ClusterConfig

	var response tools.Response

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &ClusterData); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err := UpdateNodeCountInPoolCluster(tools.MONKEYSPROJECT, ClusterData.ClusterRegion, ClusterData.ClusterName, ClusterData.ClusterNodePoolName, ClusterData.ClusterNodeCount)

		if err != nil {
			response.Message = "Cluster fail to create"
			response.Result = "Error"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		response.Message = "Cluster created succesfully"
		response.Result = "Success"
		w.WriteHeader(http.StatusOK)
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

}
