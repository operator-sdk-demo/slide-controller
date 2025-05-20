package mdparser

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

const presConfigName string = "presentation-config"
const presConfigPath string = "presentation"
const appLabelValue string = "md-parser"
const appLabel string = "app"

func CreateMarkdownParser(
	name string,
	namespace string,
	markdown string,
) (*corev1.ConfigMap, *appsv1.Deployment, *corev1.Service) {
	return createConfigMap(
			name,
			namespace,
			markdown,
		), createDeployment(
			name,
			namespace,
		), createService(
			name,
			namespace,
		)
}

func createConfigMap(name string, namespace string, markdown string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      presConfigName,
			Namespace: namespace,
		},
		Data: map[string]string{
			"presentation.md": markdown,
		},
	}
}

func createDeployment(name string, namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To[int32](1),
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					appLabel: appLabelValue,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{Labels: map[string]string{appLabel: appLabelValue}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "python",
							Image:   "python:3.13-alpine",
							Command: []string{"/bin/sh", "-c"},
							Args: []string{
								fmt.Sprintf(
									"pip install mkslides && mkslides serve %s",
									presConfigPath,
								),
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: presConfigName, MountPath: presConfigPath},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: presConfigName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: presConfigName,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func createService(name string, namespace string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				appLabel: appLabelValue,
			},
			Ports: []corev1.ServicePort{
				{Name: "ui", Port: int32(80), TargetPort: intstr.FromInt(8000)},
			},
		},
	}
}
