package gcp_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	tools "v1/tools"

	networking "k8s.io/api/networking/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createIngress(location, clusterName, nsName, AppName, DnsName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	var pathPrefix networking.PathType = networking.PathTypeImplementationSpecific

	if err != nil {
		return err
	}

	//	ingressClient := clientset.NetworkingV1().Ingresses(nsName + "-" + AppName)
	ingressClient := clientset.NetworkingV1().Ingresses("default")

	ingressReq := &networking.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "wp-ingress",
			Annotations: map[string]string{

				"kubernetes.io/ingress.class": "gce",
			},
		},
		Spec: networking.IngressSpec{
			Rules: []networking.IngressRule{
				{
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path:     "",
									PathType: &pathPrefix,
									Backend: networking.IngressBackend{
										Service: &networking.IngressServiceBackend{
											Name: "wp-service",
											Port: networking.ServiceBackendPort{
												Number: *tools.Int32Ptr(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			/* 	DefaultBackend: &networking.IngressBackend{
				Service: &networking.IngressServiceBackend{
					Name: "wp-service",
					Port: networking.ServiceBackendPort{
						Number: int32(4001),
					},
				},
			}, */
		},
	}

	result, err := ingressClient.Create(context.TODO(), ingressReq, metav1.CreateOptions{})

	fmt.Println(result)

	return err

}

func CreateIngressHandler(w http.ResponseWriter, r *http.Request) {

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

		resultingress := createIngress(K8sData.Region, K8sData.ClusterName, K8sData.EnvType, K8sData.AppName, K8sData.DnsName)

		fmt.Println(resultingress)

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
