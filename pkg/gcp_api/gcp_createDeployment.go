package gcp_api

import (
	"context"
	"fmt"
	"v1/tools"

	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(location, clusterName, nsName, AppName string) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	deploymentclient := clientset.AppsV1().Deployments(nsName + "-" + AppName)

	//deploymentclient := clientset.AppsV1().Deployments("default")
	/**Esta funcion es la que sirve para crear un service**/
	//clientset.CoreV1().Services(apiv1.NamespaceDefault)

	if AppName == "wordpress" || AppName == "drupal" {

		//My App deployment
		WPDeploymentConfig := tools.DeploymentRequest(strings.ToLower(AppName), false)

		result, err := deploymentclient.Create(context.TODO(), WPDeploymentConfig, metav1.CreateOptions{})

		if err != nil {
			return err
		}
		fmt.Println(result)

		//MySQL deployment
		MYSQLDeploymentConfig := tools.DeploymentRequest(strings.ToLower(AppName), true)

		MYSQLresult, err := deploymentclient.Create(context.TODO(), MYSQLDeploymentConfig, metav1.CreateOptions{})

		if err != nil {
			return err
		}
		fmt.Println(MYSQLresult)

	} else if AppName == "nginx" {

		NginxDeploymentConfig := tools.DeploymentRequest(strings.ToLower(AppName), false)

		result, err := deploymentclient.Create(context.TODO(), NginxDeploymentConfig, metav1.CreateOptions{})

		if err != nil {
			return err
		}

		fmt.Println(result)

	}

	return err
}
