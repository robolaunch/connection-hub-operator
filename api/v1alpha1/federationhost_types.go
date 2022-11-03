package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FederationHostSpec defines the desired state of FederationHost
type FederationHostSpec struct {
}

// FederationHostStatus defines the observed state of FederationHost
type FederationHostStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// FederationHost is the Schema for the federationhosts API
type FederationHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FederationHostSpec   `json:"spec,omitempty"`
	Status FederationHostStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FederationHostList contains a list of FederationHost
type FederationHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FederationHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FederationHost{}, &FederationHostList{})
}
