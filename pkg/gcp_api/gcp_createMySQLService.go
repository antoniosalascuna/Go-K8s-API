package gcp_api

import (
	"context"
	tools "v1/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Add resource config for MySQL
func CreateMySQLService(projectId, location, clusterName, nsName, AppName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	//Create MYSQL Service

	//ServiceClient := clientset.CoreV1().Services(nsName + "-" + AppName)

	ServiceClient := clientset.CoreV1().Services("default")

	PVCStatus := CreateNewPVC(location, clusterName, nsName, AppName, true)

	if PVCStatus != nil {
		return PVCStatus
	}

	serviceRequest := tools.ServiceRequest(AppName, true)

	_, err = ServiceClient.Create(context.TODO(), serviceRequest, metav1.CreateOptions{})

	//ADD Persisten Volume Claim into GCP Cluster

	return err

}

/* func CreateMySQLServiceHandler(w http.ResponseWriter, r *http.Request) {

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

		err := CreateMySQLService(tools.MONKEYSPROJECT, K8sData.GKElocation, K8sData.GKEclusterName, K8sData.EnvType)

		if err != nil {
			response.Message = "Cluster created succesfully"
			response.Result = "Success"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

} */
