package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type FederationMemberCredentials struct {
	// +kubebuilder:validation:Required
	CertificateAuthority string `json:"certificateAuthority"`
	// +kubebuilder:validation:Required
	ClientCertificate string `json:"clientCertificate"`
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
	Joined bool `json:"joined,omitempty"`
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

func (federationmember *FederationMember) GetOwnerMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: federationmember.OwnerReferences[0].Name,
	}
}
