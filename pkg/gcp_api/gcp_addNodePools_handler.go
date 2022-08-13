package gcp_api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	tools "v1/tools"
)

func AddNodePoolsIntoClusterHandler(w http.ResponseWriter, r *http.Request) {

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

		err := AddNodePoolsIntoCluster(tools.MONKEYSPROJECT, K8sData.Region, K8sData.ClusterName)

		if err != nil {
			response.Message = "Cluster created succesfully"
			response.Result = "Success"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

}
