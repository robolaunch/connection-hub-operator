package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FederationOperatorSpec defines the desired state of FederationOperator
type FederationOperatorSpec struct {
}

// FederationOperatorStatus defines the observed state of FederationOperator
type FederationOperatorStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// FederationOperator is the Schema for the federationoperators API
type FederationOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FederationOperatorSpec   `json:"spec,omitempty"`
	Status FederationOperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FederationOperatorList contains a list of FederationOperator
type FederationOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FederationOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FederationOperator{}, &FederationOperatorList{})
}
