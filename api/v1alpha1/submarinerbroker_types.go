package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubmarinerBrokerSpec defines the desired state of SubmarinerBroker
type SubmarinerBrokerSpec struct {
}

// SubmarinerBrokerStatus defines the observed state of SubmarinerBroker
type SubmarinerBrokerStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SubmarinerBroker is the Schema for the submarinerbrokers API
type SubmarinerBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubmarinerBrokerSpec   `json:"spec,omitempty"`
	Status SubmarinerBrokerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubmarinerBrokerList contains a list of SubmarinerBroker
type SubmarinerBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubmarinerBroker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubmarinerBroker{}, &SubmarinerBrokerList{})
}
