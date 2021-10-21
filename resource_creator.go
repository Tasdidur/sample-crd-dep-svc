package main

import (
	xapiv1 "github.com/Tasdidur/xcrd/pkg/apis/xapi.com/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(mycrd *xapiv1.Xcrd) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mycrd.Spec.Name + "-dep",
			Namespace: mycrd.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"finder": mycrd.Spec.Finder,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"finder": mycrd.Spec.Finder,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  mycrd.Spec.Name + "-contnr",
							Image: mycrd.Spec.Image,
						},
					},
				},
			},
		},
	}
}

func newService(mycrd *xapiv1.Xcrd) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: mycrd.Spec.Name + "-svc",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"finder": mycrd.Spec.Finder,
			},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol: "TCP",
					Port:     8003,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: 8081,
						StrVal: "8081",
					},
				},
			},
		},
	}
}

func newIngress(mycrd *xapiv1.Xcrd) *netv1.Ingress {
	return &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      mycrd.Spec.Name + "-ing",
			Namespace: mycrd.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/use-regex": "true",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				netv1.IngressRule{
					Host: mycrd.Spec.Domain,
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: funca(mycrd),
						},
					},
				},
			},
		},
	}
}

func funca(mycrd *xapiv1.Xcrd) []netv1.HTTPIngressPath {
	var ara []netv1.HTTPIngressPath
	var ele netv1.HTTPIngressPath
	ele.PathType = func() *netv1.PathType {
		pt := netv1.PathTypePrefix
		return &pt
	}()
	ele.Backend = netv1.IngressBackend{
		Service: &netv1.IngressServiceBackend{
			Name: mycrd.Spec.Name + "-svc",
			Port: netv1.ServiceBackendPort{
				Number: 8001,
			},
		},
	}
	for _, i := range mycrd.Spec.Paths {
		ele.Path = i
		ara = append(ara, ele)
	}
	return ara
}

func int32Ptr(i int32) *int32 { return &i }
