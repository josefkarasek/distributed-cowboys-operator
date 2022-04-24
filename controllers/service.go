package controllers

import (
	examplecomv1alpha1 "github.com/josefkarasek/distributed-cowboys-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	port        = 8080
	selectorKey = cowboy
)

func desiredService(instance *examplecomv1alpha1.Shootout, name string) *corev1.Service {
	labels := map[string]string{selectorKey: name}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       instance.Name,
					Kind:       instance.Kind,
					UID:        instance.UID,
					APIVersion: instance.APIVersion,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       port,
					TargetPort: intstr.FromInt(port),
				},
			},
		},
	}

	return svc
}
