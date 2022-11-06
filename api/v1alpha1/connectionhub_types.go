package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConnectionHubSpec defines the desired state of ConnectionHub
type ConnectionHubSpec struct {
}

// ConnectionHubStatus defines the observed state of ConnectionHub
type ConnectionHubStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// ConnectionHub is the Schema for the connectionhubs API
type ConnectionHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionHubSpec   `json:"spec,omitempty"`
	Status ConnectionHubStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectionHubList contains a list of ConnectionHub
type ConnectionHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectionHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConnectionHub{}, &ConnectionHubList{})
}
