package tools

import (
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ServiceRequest(AppName string, isMySQL bool) *apiv1.Service {

	if strings.ToLower(AppName) == "wordpress" && !isMySQL {
		return WordpressServiceConfig(AppName)
	}
	if isMySQL {
		return MySQLServiceConfig(AppName)
	}

	if strings.ToLower(AppName) == "drupal" && !isMySQL {
		return DrupalServiceConfig(AppName)
	}

	if (strings.ToLower(AppName) == "nginx") && !isMySQL {

		return NginxServiceConfig(AppName)

	}
	return nil

}

func WordpressServiceConfig(AppName string) *apiv1.Service {

	service := &apiv1.Service{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "service",
			Labels: map[string]string{
				"app": AppName + "-" + "deployment",
			},
		},

		Spec: apiv1.ServiceSpec{

			Ports: []apiv1.ServicePort{

				{
					Port:       *Int32Ptr(80),
					Name:       "http",
					TargetPort: intstr.FromInt(80),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": AppName + "-" + "deployment",
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	return service
}

func MySQLServiceConfig(AppName string) *apiv1.Service {
	service := &apiv1.Service{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql" + "-" + "service",
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{

				{
					Port: 3306,
				},
			},
			Selector: map[string]string{
				"app":  AppName,
				"tier": "backend",
			},
		},
	}

	return service
}

func DrupalServiceConfig(AppName string) *apiv1.Service {
	service := &apiv1.Service{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "service",
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{

				{
					Port:       *Int32Ptr(80),
					Name:       "web",
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app":  AppName,
				"tier": "frontend",
			},

			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	return service
}

func NginxServiceConfig(AppName string) *apiv1.Service {
	service := &apiv1.Service{

		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName,
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{

				{
					Port:       *Int32Ptr(80),
					Name:       "nginx",
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app":  AppName,
				"type": "frontend",
			},

			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	return service
}
