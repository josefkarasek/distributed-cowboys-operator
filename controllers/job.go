package controllers

import (
	"fmt"

	examplecomv1alpha1 "github.com/josefkarasek/distributed-cowboys-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	img             = "quay.io/josefkarasek/distributed-cowboys:v0.0.1"
	cmd             = "/distributed-cowboys"
	configMountPath = "/var/run/cowboys"
	podName         = cowboy
)

func desiredJob(instance *examplecomv1alpha1.Shootout, jobName, cowboyName, neighbor string, coordinator bool) *batchv1.Job {
	labels := map[string]string{selectorKey: jobName}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: instance.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{{
				Name:       instance.Name,
				Kind:       instance.Kind,
				UID:        instance.UID,
				APIVersion: instance.APIVersion,
			}},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:            podName,
						Image:           img,
						ImagePullPolicy: corev1.PullAlways,
						Command:         []string{cmd},
						Args: []string{
							fmt.Sprintf("--name=%s", cowboyName),
							fmt.Sprintf("--shootout-name=%s", instance.Name),
							fmt.Sprintf("--neighbor=%s", neighbor),
							fmt.Sprintf("--coordinator=%t", coordinator),
						},
						Ports: []corev1.ContainerPort{
							{
								Protocol:      corev1.ProtocolTCP,
								ContainerPort: port,
							},
						},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      configKey,
							MountPath: configMountPath,
							ReadOnly:  true,
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: configKey,
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: instance.Name,
								},
							},
						},
					}},
				},
			},
		},
	}

	return job
}
