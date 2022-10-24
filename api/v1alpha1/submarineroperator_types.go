package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubmarinerOperatorSpec defines the desired state of SubmarinerOperator
type SubmarinerOperatorSpec struct {
}

// SubmarinerOperatorStatus defines the observed state of SubmarinerOperator
type SubmarinerOperatorStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// SubmarinerOperator is the Schema for the submarineroperators API
type SubmarinerOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubmarinerOperatorSpec   `json:"spec,omitempty"`
	Status SubmarinerOperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubmarinerOperatorList contains a list of SubmarinerOperator
type SubmarinerOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubmarinerOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubmarinerOperator{}, &SubmarinerOperatorList{})
}
