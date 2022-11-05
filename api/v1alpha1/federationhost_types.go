package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MemberInfo struct {
	// +kubebuilder:validation:Required
	MemberSpec FederationMemberSpec `json:"memberSpec"`
}

type MemberResourcePhase string

const (
	MemberResourcePhaseActive MemberResourcePhase = "Active"
	MemberResourcePhaseIdle   MemberResourcePhase = "Idle"
)

type MemberStatus struct {
	Created       bool                   `json:"created,omitempty"`
	Status        FederationMemberStatus `json:"status,omitempty"`
	ResourcePhase MemberResourcePhase    `json:"resourcePhase,omitempty"`
}

// FederationHostSpec defines the desired state of FederationHost
type FederationHostSpec struct {
	FederationMembers map[string]MemberInfo `json:"members,omitempty"`
}

type FederationHostPhase string

const (
	FederationHostPhaseJoiningSelf     FederationHostPhase = "JoiningSelf"
	FederationHostPhaseReady           FederationHostPhase = "Ready"
	FederationHostPhaseDeletingMembers FederationHostPhase = "DeletingMembers"
)

// FederationHostStatus defines the observed state of FederationHost
type FederationHostStatus struct {
	SelfJoined     bool                    `json:"selfJoined,omitempty"`
	MemberStatuses map[string]MemberStatus `json:"memberStatuses,omitempty"`
	Phase          FederationHostPhase     `json:"phase,omitempty"`
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
