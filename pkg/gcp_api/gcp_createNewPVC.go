package gcp_api

import (
	"context"

	tools "v1/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNewPVC(location, clusterName, nsName, AppName string, isMySQL bool) error {

	clientset, err := tools.GetClientSetByJsonFile(location, clusterName)

	if err != nil {
		return err
	}

	if isMySQL == true {
		pv := clientset.CoreV1().PersistentVolumeClaims(nsName + "-" + AppName)

		pvRequest := tools.PvcRequest("mysql")

		_, err = pv.Create(context.TODO(), pvRequest, metav1.CreateOptions{})

		if err != nil {
			return err
		}

	} else {

		pv := clientset.CoreV1().PersistentVolumeClaims(nsName + "-" + AppName)

		pvRequest := tools.PvcRequest(AppName)

		_, err = pv.Create(context.TODO(), pvRequest, metav1.CreateOptions{})

		if err != nil {
			return err
		}

	}

	//pv := clientset.CoreV1().PersistentVolumeClaims("default")

	return err
}
