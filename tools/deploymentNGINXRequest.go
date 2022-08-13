package tools

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploymentNGINXRequest(AppName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{

		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": AppName,
				},
			},
			Replicas: Int32Ptr(1),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": AppName,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{

						{
							Name:  "nginx",
							Image: "nginx:latest",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 81,
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}
