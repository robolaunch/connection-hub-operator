package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SubmarinerOperatorNamespace string = "submariner-operator"
)

type SubmarinerOperatorResourceItem struct {
	client.ObjectKey
	GroupVersionKind metav1.GroupVersionKind
}

func (so *SubmarinerOperator) GetResourcesForCheck() []SubmarinerOperatorResourceItem {
	return []SubmarinerOperatorResourceItem{
		{
			ObjectKey: types.NamespacedName{
				Namespace: SubmarinerOperatorNamespace,
				Name:      "submariner-operator",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
	}
}

// SubmarinerOperatorSpec defines the desired state of SubmarinerOperator
type SubmarinerOperatorSpec struct {
	// +kubebuilder:validation:Required
	ClusterCIDR string `json:"clusterCIDR"`
	// +kubebuilder:validation:Required
	ServiceCIDR string `json:"serviceCIDR"`
	// +kubebuilder:validation:Required
	PresharedKey string `json:"presharedKey"`
	// +kubebuilder:validation:Required
	Broker BrokerInfo `json:"broker"`
	// +kubebuilder:validation:Required
	ClusterID string `json:"clusterID"`
	// +kubebuilder:validation:Required
	Helm HelmChartProperties `json:"helm"`
}

type SubmarinerOperatorPhase string

const (
	SubmarinerOperatorPhaseCreatingNamespace SubmarinerOperatorPhase = "CreatingNamespace"
	SubmarinerOperatorPhaseDeployingChart    SubmarinerOperatorPhase = "DeployingChart"
	SubmarinerOperatorPhaseCheckingResources SubmarinerOperatorPhase = "CheckingResources"
	SubmarinerOperatorPhaseDeployed          SubmarinerOperatorPhase = "Deployed"
	SubmarinerOperatorPhaseMalfunctioned     SubmarinerOperatorPhase = "Malfunctioned"
)

type OperatorNamespaceStatus struct {
	Created bool `json:"created,omitempty"`
}

type ChartStatus struct {
	Deployed bool `json:"deployed,omitempty"`
}

type ChartResourceStatus struct {
	Deployed bool `json:"deployed,omitempty"`
}

// SubmarinerOperatorStatus defines the observed state of SubmarinerOperator
type SubmarinerOperatorStatus struct {
	OperatorNamespaceStatus OperatorNamespaceStatus `json:"operatorNamespaceStatus,omitempty"`
	ChartStatus             ChartStatus             `json:"chartStatus,omitempty"`
	ChartResourceStatus     ChartResourceStatus     `json:"chartResourceStatus,omitempty"`
	Phase                   SubmarinerOperatorPhase `json:"phase,omitempty"`
	NodeInfo                K8sNodeInfo             `json:"nodeInfo,omitempty"`
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

func GetTenancySelectorsForSO(submarinerOperator SubmarinerOperator) *Tenancy {

	tenancy := &Tenancy{}
	labels := submarinerOperator.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	if physicalInstance, ok := labels[RobolaunchPhysicalInstanceLabelKey]; ok {
		tenancy.RobolaunchPhysicalInstance = physicalInstance
	}

	return tenancy
}
