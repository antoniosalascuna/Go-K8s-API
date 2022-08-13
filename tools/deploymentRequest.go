package tools

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploymentRequest(AppName string, isMySQL bool) *appsv1.Deployment {

	if strings.ToLower(AppName) == "wordpress" && !isMySQL {
		return WPDeploymentConfig(AppName)
	}
	if strings.ToLower(AppName) == "drupal" && !isMySQL {
		return DPDeploymentConfig(AppName)
	}

	if strings.ToLower(AppName) == "nginx" && !isMySQL {
		return NGINXDeploymentConfig(AppName)
	}

	if isMySQL {
		return DeploymentMySQLRequest(AppName)
	}
	return nil
}

func WPDeploymentConfig(AppName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "deployment",
			Labels: map[string]string{
				"app": AppName + "-" + "deployment",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": AppName + "-" + "deployment",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": AppName + "-" + "deployment",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "wordpress",
							Image: "wordpress:4.8-apache",

							Env: []apiv1.EnvVar{
								{
									Name:  "WORDPRESS_DB_HOST",
									Value: "mysql-service",
								},
								{
									Name: "WORDPRESS_DB_PASSWORD",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: AppName + "-" + "mysql" + "-" + "pass",
											},
											Key: "password",
										},
									},
								},
							},

							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: *int32Ptr(80),
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      AppName + "-" + "volumeclaim",
									MountPath: "/var/www/html",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: AppName + "-" + "volumeclaim",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: AppName + "-" + "volumeclaim",
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

func DPDeploymentConfig(AppName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "deployment",
			Labels: map[string]string{
				"app":  AppName,
				"tier": "frontend",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  AppName,
					"tier": "frontend",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  AppName,
						"tier": "frontend",
					},
				},
				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{
						{
							Name:  "init-sites-volume",
							Image: "drupal:8.6.15-apache",
							Command: []string{
								"/bin/bash",
								"-c",
							},
							Args: []string{
								"cp -r /var/www/html/sites /data",
								"chown www-data:www-data /data/ -R",
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/data",
									Name:      AppName + "-" + "volumeclaim",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:  AppName,
							Image: "drupal:8.9.20-php7.4-apache",

							Env: []apiv1.EnvVar{
								{
									Name:  "DRUPAL_DB_HOST",
									Value: "mysql-service",
								},
								{
									Name: "DRUPAL_DB_PASSWORD",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: AppName + "-" + "mysql-pass",
											},
											Key: "password",
										},
									},
								},
							},

							Ports: []apiv1.ContainerPort{
								{
									Name:          "drupal",
									ContainerPort: *int32Ptr(80),
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      AppName + "-" + "volumeclaim",
									MountPath: "/var/www/html/modules",
									SubPath:   "modules",
								},
								{
									Name:      AppName + "-" + "volumeclaim",
									MountPath: "/var/www/html/profiles",
									SubPath:   "profiles",
								},
								{
									Name:      AppName + "-" + "volumeclaim",
									MountPath: "/var/www/html/sites",
									SubPath:   "sites",
								},
								{
									Name:      AppName + "-" + "volumeclaim",
									MountPath: "/var/www/html/themes",
									SubPath:   "themes",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: AppName + "-" + "volumeclaim",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: AppName + "-" + "volumeclaim",
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

func NGINXDeploymentConfig(AppName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: AppName + "-" + "deployment",
			Labels: map[string]string{
				"app":  AppName,
				"type": "frontend",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  AppName,
					"type": "frontend",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  AppName,
						"type": "frontend",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  AppName,
							Image: "nginx:latest",

							Ports: []apiv1.ContainerPort{
								{
									Name:          "nginx",
									ContainerPort: *int32Ptr(80),
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

func DeploymentMySQLRequest(AppName string) *appsv1.Deployment {

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql" + "-" + "deployment",
			Labels: map[string]string{
				"app": AppName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  AppName,
					"tier": "backend",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  AppName,
						"tier": "backend",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "mysql",
							Image: "mysql:5.6",
							Env: []apiv1.EnvVar{
								{
									Name: "MYSQL_ROOT_PASSWORD",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: AppName + "-" + "mysql" + "-" + "pass",
											},
											Key: "password",
										},
									},
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "mysql",
									ContainerPort: 3306,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "mysql-persistent-storage",
									MountPath: "/var/lib/mysql",
									SubPath:   "dbdata",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "mysql-persistent-storage",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mysql-volumeclaim",
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

func int32Ptr(i int32) *int32 { return &i }
