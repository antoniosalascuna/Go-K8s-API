package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	tools "v1/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateService(AppName, nsName, clusterName, region string) error {

	clientset, err := tools.GetClientSetByJsonFile(region, clusterName)

	if err != nil {
		return err
	}

	//Create Service
	ServiceClient := clientset.CoreV1().Services(nsName + "-" + AppName)

	if AppName == "wordpress" || AppName == "drupal" {
		//MyAppPVC
		pvcStatus := CreateNewPVC(region, clusterName, nsName, AppName, false)

		if pvcStatus != nil {
			return pvcStatus
		}

		//MySQLPVC
		MySQLPVCStatus := CreateNewPVC(region, clusterName, nsName, AppName, true)

		if MySQLPVCStatus != nil {
			return MySQLPVCStatus
		}
		//APP SERVICE (WP or DRUPAL)
		serviceRequest := tools.ServiceRequest(strings.ToLower(AppName), false)

		_, err = ServiceClient.Create(context.TODO(), serviceRequest, metav1.CreateOptions{})

		if err != nil {
			return err
		}

		//MySQL SERVICE
		MySQLserviceRequest := tools.ServiceRequest(strings.ToLower(AppName), true)

		_, err = ServiceClient.Create(context.TODO(), MySQLserviceRequest, metav1.CreateOptions{})

		if err != nil {
			return err
		}

		return err

		//NGINX SERVICE (STATELESS APP)
	} else if AppName == "nginx" {

		serviceRequest := tools.ServiceRequest(strings.ToLower(AppName), false)

		_, err = ServiceClient.Create(context.TODO(), serviceRequest, metav1.CreateOptions{})

		if err != nil {
			return err
		}

	}

	return nil

}

func CreateServiceHandler(w http.ResponseWriter, r *http.Request) {

	var ServiceData tools.ServiceConfig

	var response tools.Response

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &ServiceData); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err := CreateService(ServiceData.AppName, ServiceData.NsName, ServiceData.ClusterName, ServiceData.Region)

		if err != nil {
			response.Message = "Service created succesfully"
			response.Result = "Success"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

}
