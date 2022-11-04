package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MemberInfo struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	MemberSpec FederationMemberSpec `json:"memberSpec"`
}

type MemberStatus struct {
	Name    string `json:"name"`
	Created bool   `json:"created,omitempty"`
}

// FederationHostSpec defines the desired state of FederationHost
type FederationHostSpec struct {
	FederationMembers []MemberInfo `json:"members,omitempty"`
}

// FederationHostStatus defines the observed state of FederationHost
type FederationHostStatus struct {
	SelfJoined     bool           `json:"selfJoined,omitempty"`
	MemberStatuses []MemberStatus `json:"memberStatuses,omitempty"`
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
