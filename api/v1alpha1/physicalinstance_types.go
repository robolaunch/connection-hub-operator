package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PhysicalInstanceSpec defines the desired state of PhysicalInstance
type PhysicalInstanceSpec struct {
}

// PhysicalInstanceStatus defines the observed state of PhysicalInstance
type PhysicalInstanceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// PhysicalInstance is the Schema for the physicalinstances API
type PhysicalInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalInstanceSpec   `json:"spec,omitempty"`
	Status PhysicalInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PhysicalInstanceList contains a list of PhysicalInstance
type PhysicalInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalInstance{}, &PhysicalInstanceList{})
}
