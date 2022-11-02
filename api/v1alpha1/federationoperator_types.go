package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	FederationOperatorNamespace string = "kube-federation-system"
)

// FederationOperatorSpec defines the desired state of FederationOperator
type FederationOperatorSpec struct {
	// +kubebuilder:validation:Required
	Helm HelmChartProperties `json:"helm"`
}

type FederationOperatorPhase string

const (
	FederationOperatorPhaseCreatingNamespace FederationOperatorPhase = "CreatingNamespace"
	FederationOperatorPhaseDeployingChart    FederationOperatorPhase = "DeployingChart"
	FederationOperatorPhaseCheckingResources FederationOperatorPhase = "CheckingResources"
	FederationOperatorPhaseDeployed          FederationOperatorPhase = "Deployed"
	FederationOperatorPhaseMalfunctioned     FederationOperatorPhase = "Malfunctioned"

	FederationOperatorPhaseUninstallingChart    FederationOperatorPhase = "UninstallingChart"
	FederationOperatorPhaseTerminatingNamespace FederationOperatorPhase = "TerminatingNamespace"
)

// FederationOperatorStatus defines the observed state of FederationOperator
type FederationOperatorStatus struct {
	NamespaceStatus     NamespaceStatus         `json:"namespaceStatus,omitempty"`
	ChartStatus         ChartStatus             `json:"chartStatus,omitempty"`
	ChartResourceStatus ChartResourceStatus     `json:"chartResourceStatus,omitempty"`
	Phase               FederationOperatorPhase `json:"phase,omitempty"`
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

func (fo *FederationOperator) GetNamespaceMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: FederationOperatorNamespace,
	}
}

func (fo *FederationOperator) GetResourcesForCheck() []ResourceItem {
	return []ResourceItem{
		{
			ObjectKey: types.NamespacedName{
				Namespace: fo.GetNamespaceMetadata().Name,
				Name:      "kubefed-controller-manager",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
		{
			ObjectKey: types.NamespacedName{
				Namespace: fo.GetNamespaceMetadata().Name,
				Name:      "kubefed-admission-webhook",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
	}
}
