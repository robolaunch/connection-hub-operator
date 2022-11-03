package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FederationMemberCredentials struct {
	// +kubebuilder:validation:Required
	CertificateAuthorityData string `json:"certificateAuthorityData"`
	// +kubebuilder:validation:Required
	ClientCertificateData string `json:"clientCertificateData"`
	// +kubebuilder:validation:Required
	ClientKey string `json:"clientKey"`
}

// FederationMemberSpec defines the desired state of FederationMember
type FederationMemberSpec struct {
	// +kubebuilder:validation:Required
	Server string `json:"server"`
	// +kubebuilder:validation:Required
	Credentials FederationMemberCredentials `json:"credentials"`
}

// FederationMemberStatus defines the observed state of FederationMember
type FederationMemberStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// FederationMember is the Schema for the federationmembers API
type FederationMember struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FederationMemberSpec   `json:"spec,omitempty"`
	Status FederationMemberStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FederationMemberList contains a list of FederationMember
type FederationMemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FederationMember `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FederationMember{}, &FederationMemberList{})
}
