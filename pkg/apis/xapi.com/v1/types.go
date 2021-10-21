package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Xcrd struct{
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XcrdSpec   `json:"spec"`
	Status XcrdStatus `json:"status"`
}

type XcrdSpec struct{

	Name string `json:"name"`
	Finder string `json:"finder"`
	Domain string `json:"domain"`
	Image string `json:"image"`
	Port int `json:"port"`
	TargetPort int `json:"target-port"`
	Paths []string `json:"paths"`

}

type XcrdStatus struct {
	AllReady bool `json:"all_ready"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type XcrdList struct{
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items []Xcrd `json:"items"`
}