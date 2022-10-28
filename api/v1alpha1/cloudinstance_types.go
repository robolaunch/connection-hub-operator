package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CloudInstanceSpec defines the desired state of CloudInstance
type CloudInstanceSpec struct {
}

// CloudInstanceStatus defines the observed state of CloudInstance
type CloudInstanceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// CloudInstance is the Schema for the cloudinstances API
type CloudInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudInstanceSpec   `json:"spec,omitempty"`
	Status CloudInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CloudInstanceList contains a list of CloudInstance
type CloudInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudInstance{}, &CloudInstanceList{})
}
