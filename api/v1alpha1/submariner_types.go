package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type InstanceType string

const (
	InstanceTypeCloud    InstanceType = "CloudInstance"
	InstanceTypePhysical InstanceType = "PhysicalInstance"
)

// SubmarinerSpec defines the desired state of Submariner
type SubmarinerSpec struct {
	// +kubebuilder:validation:Enum=CloudInstance;PhysicalInstance
	// +kubebuilder:validation:Required
	InstanceType InstanceType `json:"instanceType"`
	// +kubebuilder:validation:Required
	ClusterCIDR string `json:"clusterCIDR"`
	// +kubebuilder:default="10.32.0.0/16"
	ServiceCIDR string `json:"serviceCIDR,omitempty"`

	PresharedKey string `json:"presharedKey,omitempty"`

	BrokerCredentials BrokerCredentials `json:"broker,omitempty"`
	// +kubebuilder:validation:Required
	APIServerURL string `json:"apiServerURL"`
	// +kubebuilder:validation:Required
	BrokerHelmChart HelmChartProperties `json:"brokerHelmChart"`
	// +kubebuilder:validation:Required
	OperatorHelmChart HelmChartProperties `json:"operatorHelmChart"`
}

type BrokerStatus struct {
	Created bool                   `json:"created,omitempty"`
	Phase   SubmarinerBrokerPhase  `json:"phase,omitempty"`
	Status  SubmarinerBrokerStatus `json:"status,omitempty"`
}

type OperatorStatus struct {
	Created bool                    `json:"created,omitempty"`
	Phase   SubmarinerOperatorPhase `json:"phase,omitempty"`
}

type CustomResourceStatus struct {
	Created             bool                `json:"created,omitempty"`
	OwnedResourceStatus OwnedResourceStatus `json:"ownedResourceStatus,omitempty"`
}

type OwnedResourceStatus struct {
	Deployed bool `json:"deployed,omitempty"`
}

type SubmarinerPhase string

const (
	SubmarinerPhaseCreatingBroker         SubmarinerPhase = "CreatingBroker"
	SubmarinerPhaseCreatingOperator       SubmarinerPhase = "CreatingOperator"
	SubmarinerPhaseCreatingCustomResource SubmarinerPhase = "CreatingCustomResource"
	SubmarinerPhaseCheckingResources      SubmarinerPhase = "CheckingResources"
	SubmarinerPhaseReadyToConnect         SubmarinerPhase = "ReadyToConnect"
	SubmarinerPhaseMalfunctioned          SubmarinerPhase = "Malfunctioned"

	SubmarinerPhaseTerminatingSubmarinerCR       SubmarinerPhase = "TerminatingSubmarinerCR"
	SubmarinerPhaseTerminatingSubmarinerOperator SubmarinerPhase = "TerminatingSubmarinerOperator"
	SubmarinerPhaseTerminatingSubmarinerBroker   SubmarinerPhase = "TerminatingSubmarinerBroker"
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
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

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

func (submariner *Submariner) GetTenancySelectors() *Tenancy {

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

func (submariner *Submariner) GetResourcesForCheck() []ResourceItem {
	return []ResourceItem{
		{
			ObjectKey: types.NamespacedName{
				Namespace: SubmarinerOperatorNamespace,
				Name:      "submariner-lighthouse-agent",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
		{
			ObjectKey: types.NamespacedName{
				Namespace: SubmarinerOperatorNamespace,
				Name:      "submariner-lighthouse-coredns",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
		{
			ObjectKey: types.NamespacedName{
				Namespace: SubmarinerOperatorNamespace,
				Name:      "submariner-gateway",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "DaemonSet",
			},
		},
		{
			ObjectKey: types.NamespacedName{
				Namespace: SubmarinerOperatorNamespace,
				Name:      "submariner-routeagent",
			},
			GroupVersionKind: metav1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "DaemonSet",
			},
		},
	}
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

func (submariner *Submariner) GetSubmarinerCustomResourceMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name:      submariner.Name + "-cr",
		Namespace: SubmarinerOperatorNamespace,
	}
}
