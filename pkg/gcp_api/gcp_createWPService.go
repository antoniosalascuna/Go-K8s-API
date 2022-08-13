package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	tools "v1/tools"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateWPService(projectId, location, clusterName, nsName, AppName, DnsName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	//Create WP Service

	//ServiceClient := clientset.CoreV1().Services(nsName + "-" + AppName)

	ServiceClient := clientset.CoreV1().Services("default")

	service := &apiv1.Service{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "wp-service",
			Labels: map[string]string{
				"run": "wordpress",
			},
		},

		Spec: apiv1.ServiceSpec{

			Ports: []apiv1.ServicePort{

				{
					Port:       *tools.Int32Ptr(80),
					Name:       "http",
					TargetPort: intstr.FromInt(80),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"run": "wordpress",
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}

	//var ips = ""

	ServiceResult, err := ServiceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	IpWpService, err := ServiceClient.Get(context.TODO(), service.ObjectMeta.Name, metav1.GetOptions{})

	ServicePorts := IpWpService.Status.LoadBalancer.Ingress

	time.Sleep(2 * time.Second)

	fmt.Println(ServiceResult)

	fmt.Println(ServicePorts)

	/* for ips == "" {

		IpWpService, err := ServiceClient.Get(context.TODO(), service.ObjectMeta.Name, metav1.GetOptions{})

		if err != nil {
			return err
		}

		ServicePorts := IpWpService.Status.LoadBalancer.Ingress

		for _, i := range ServicePorts {

			ips = i.IP

		}
	}

	if ips != "" {
		_, err := CreateClientDnsName(nsName, DnsName, ips)

		if err != nil {
			return err
		}

	} */

	//ADD Persisten Volume Claim into GCP Cluster

	PVCstatus := CreateNewPVC(location, clusterName, nsName, AppName, false)

	if PVCstatus != nil {
		return PVCstatus
	}

	return err
}

func CreateMyWPServiceHandler(w http.ResponseWriter, r *http.Request) {

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

		err := CreateWPService(tools.MONKEYSPROJECT, K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName, K8sData.DnsName)

		if err != nil {
			response.Message = "Cluster created succesfully"
			response.Result = "Success"
			w.WriteHeader(http.StatusBadRequest)
			encodeData, _ := json.Marshal(err)
			fmt.Fprintf(w, string(encodeData))
			return
		}

	}

}
