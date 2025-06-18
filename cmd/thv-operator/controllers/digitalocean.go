package controllers

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mcpv1alpha1 "github.com/stacklok/toolhive/cmd/thv-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const (
	DigitalOceanProvider      = "digitalocean"
	AnnotationClusterProvider = "k8s.io/cluster-provider"
)

func IsDigitalOceanProvider(mcpServer *mcpv1alpha1.MCPServer) bool {
	return mcpServer.Annotations[AnnotationClusterProvider] == DigitalOceanProvider
}

func GetDigitalOceanProvider(mcpServer *mcpv1alpha1.MCPServer) string {
	return mcpServer.Annotations[AnnotationClusterProvider]
}

func buildIngress(mcpServer *mcpv1alpha1.MCPServer, svc *corev1.Service) *networkingv1.Ingress {
	if !IsDigitalOceanProvider(mcpServer) {
		return nil
	}
	// ingress to the service
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mcpServer.Name,
			Namespace: mcpServer.Namespace,
			Annotations: map[string]string{
				AnnotationClusterProvider: DigitalOceanProvider,
			},
		},
		// spec, select by labels
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: mcpServer.Name + "." + mcpServer.Namespace + ".nip.io",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: func() *networkingv1.PathType { pt := networkingv1.PathTypePrefix; return &pt }(),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: svc.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
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
