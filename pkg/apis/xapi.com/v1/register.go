package v1

import (
	"github.com/Tasdidur/xcrd/pkg/apis/xapi.com"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)


var SchemeGroupVersion = schema.GroupVersion{Group: xapi_com.GroupName,Version:"v1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource (resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme	=	SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme)  error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Xcrd{},
		&XcrdList{},
	)
	metav1.AddToGroupVersion(scheme,SchemeGroupVersion)
	return nil

}