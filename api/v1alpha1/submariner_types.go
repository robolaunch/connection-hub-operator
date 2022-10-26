package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// SubmarinerSpec defines the desired state of Submariner
type SubmarinerSpec struct {
	// +kubebuilder:validation:Required
	APIServerURL      string              `json:"apiServerURL"`
	BrokerHelmChart   HelmChartProperties `json:"brokerHelmChart,omitempty"`
	OperatorHelmChart HelmChartProperties `json:"operatorHelmChart,omitempty"`
}

type BrokerStatus struct {
	Created bool                  `json:"created,omitempty"`
	Phase   SubmarinerBrokerPhase `json:"phase,omitempty"`
}

type OperatorStatus struct {
	Created bool                    `json:"created,omitempty"`
	Phase   SubmarinerOperatorPhase `json:"phase,omitempty"`
}

type CustomResourceStatus struct {
	Created bool `json:"created,omitempty"`
}

type SubmarinerPhase string

const (
	SubmarinerPhaseCreatingBroker         SubmarinerPhase = "CreatingBroker"
	SubmarinerPhaseCreatingOperator       SubmarinerPhase = "CreatingOperator"
	SubmarinerPhaseCreatingCustomResource SubmarinerPhase = "CreatingCustomResource"
	SubmarinerPhaseReadyToConnect         SubmarinerPhase = "ReadyToConnect"
	SubmarinerPhaseMalfunctioned          SubmarinerPhase = "Malfunctioned"
)

// SubmarinerStatus defines the observed state of Submariner
type SubmarinerStatus struct {
	Phase    SubmarinerPhase `json:"phase,omitempty"`
	NodeInfo K8sNodeInfo     `json:"nodeInfo,omitempty"`

	BrokerStatus         BrokerStatus         `json:"brokerStatus,omitempty"`
	OperatorStatus       OperatorStatus       `json:"operatorStatus,omitempty"`
	CustomResourceStatus CustomResourceStatus `json:"customResourceStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// Submariner is the Schema for the submariners API
type Submariner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubmarinerSpec   `json:"spec,omitempty"`
	Status SubmarinerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubmarinerList contains a list of Submariner
type SubmarinerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Submariner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Submariner{}, &SubmarinerList{})
}

func GetTenancySelectorsForSubmariner(submariner Submariner) *Tenancy {

	tenancy := &Tenancy{}
	labels := submariner.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	if physicalInstance, ok := labels[RobolaunchPhysicalInstanceLabelKey]; ok {
		tenancy.RobolaunchPhysicalInstance = physicalInstance
	}

	return tenancy
}

func (submariner *Submariner) GetSubmarinerBrokerMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: submariner.Name + "-broker",
	}
}

func (submariner *Submariner) GetSubmarinerOperatorMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: submariner.Name + "-operator",
	}
}
