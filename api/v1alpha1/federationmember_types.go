package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kubefed/pkg/apis/core/common"
)

// FederationMemberSpec defines the desired state of FederationMember
type FederationMemberSpec struct {
	// +kubebuilder:validation:Required
	Server string `json:"server"`
	// +kubebuilder:validation:Required
	Credentials PhysicalInstanceCredentials `json:"credentials"`
}

type FederationMemberPhase string

const (
	FederationMemberPhaseJoiningFederation FederationMemberPhase = "JoiningFederation"
	FederationMemberPhaseOffline           FederationMemberPhase = "Offline"
	FederationMemberPhaseReady             FederationMemberPhase = "Ready"
	FederationMemberPhaseMalfunctioned     FederationMemberPhase = "Malfunctioned"

	FederationMemberPhaseCannotJoinFederation FederationMemberPhase = "CannotJoinFederation"
	FederationMemberPhaseUnjoiningFederation  FederationMemberPhase = "UnjoiningFederation"
)

type KubeFedClusterStatus struct {
	Created       bool                        `json:"created,omitempty"`
	ConditionType common.ClusterConditionType `json:"conditionType,omitempty"`
	Reason        string                      `json:"reason,omitempty"`
}

// FederationMemberStatus defines the observed state of FederationMember
type FederationMemberStatus struct {
	JoinAttempted        bool                  `json:"joinAttempted,omitempty"`
	KubeFedClusterStatus KubeFedClusterStatus  `json:"kubefedClusterStatus,omitempty"`
	Phase                FederationMemberPhase `json:"phase,omitempty"`
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

type TargetCluster string

const (
	TargetClusterHost   TargetCluster = "Host"
	TargetClusterMember TargetCluster = "Member"
)

type MulticlusterResourceItem struct {
	TargetCluster
	client.ObjectKey
	GroupVersionKind metav1.GroupVersionKind
}

func (fm *FederationMember) GetResourcesForCheck() []MulticlusterResourceItem {
	return []MulticlusterResourceItem{
		{
			TargetCluster: TargetClusterMember,
			ObjectKey: types.NamespacedName{
				Namespace: FederationOperatorNamespace,
				Name:      fm.Name + "-" + fm.GetOwnerMetadata().Name,
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "ServiceAccount",
			},
		},
		{
			TargetCluster: TargetClusterMember,
			ObjectKey: types.NamespacedName{
				Namespace: FederationOperatorNamespace,
				Name:      "kubefed-controller-manager:" + fm.Name + "-" + fm.GetOwnerMetadata().Name,
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "rbac.authorization.k8s.io",
				Version: "v1",
				Kind:    "ClusterRole",
			},
		},
		{
			TargetCluster: TargetClusterMember,
			ObjectKey: types.NamespacedName{
				Namespace: FederationOperatorNamespace,
				Name:      "kubefed-controller-manager:" + fm.Name + "-" + fm.GetOwnerMetadata().Name,
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "rbac.authorization.k8s.io",
				Version: "v1",
				Kind:    "ClusterRoleBinding",
			},
		},
		// {
		// 	TargetCluster: TargetClusterMember,
		// 	ObjectKey: types.NamespacedName{
		// 		Namespace: FederationOperatorNamespace,
		// 		Name:      fm.Name + "-" + fm.GetOwnerMetadata().Name + "-token-xv87e",
		// 	},
		// 	GroupVersionKind: metav1.GroupVersionKind{
		// 		Group:   "",
		// 		Version: "v1",
		// 		Kind:    "Secret",
		// 	},
		// },
		// {
		// 	TargetCluster: TargetClusterHost,
		// 	ObjectKey: types.NamespacedName{
		// 		Namespace: FederationOperatorNamespace,
		// 		Name:      fm.Name + "-token-xv87e",
		// 	},
		// 	GroupVersionKind: metav1.GroupVersionKind{
		// 		Group:   "",
		// 		Version: "v1",
		// 		Kind:    "Secret",
		// 	},
		// },
		{
			TargetCluster: TargetClusterHost,
			ObjectKey: types.NamespacedName{
				Namespace: FederationOperatorNamespace,
				Name:      fm.Name,
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "core.kubefed.io",
				Version: "v1beta1",
				Kind:    "KubeFedCluster",
			},
		},
	}
}
