package controllers

import (
	examplecomv1alpha1 "github.com/josefkarasek/distributed-cowboys-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func desiredConfigMap(instance *examplecomv1alpha1.Shootout) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       instance.Name,
					Kind:       instance.Kind,
					UID:        instance.UID,
					APIVersion: instance.APIVersion,
				},
			},
		},
		Data: map[string]string{
			configKey: instance.Spec.Cowboys,
		},
	}
	return cm
}
