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

type SubmarinerInstanceStatus struct {
	Created bool            `json:"created,omitempty"`
	Phase   SubmarinerPhase `json:"phase,omitempty"`
}

type FederationInstanceStatus struct {
	Created bool                    `json:"created,omitempty"`
	Phase   FederationOperatorPhase `json:"phase,omitempty"`
}

type ConnectionHubPhase string

const (
	ConnectionHubPhaseSubmarinerSettingUp ConnectionHubPhase = "SubmarinerSettingUp"
	ConnectionHubPhaseFederationSettingUp ConnectionHubPhase = "FederationSettingUp"
	ConnectionHubPhaseReadyForOperation   ConnectionHubPhase = "ReadyForOperation"

	ConnectionHubPhaseLabelsNotMatched ConnectionHubPhase = "LabelsNotMatched"
	ConnectionHubPhaseMalfunctioned    ConnectionHubPhase = "Malfunctioned"
)

// ConnectionHubSpec defines the desired state of ConnectionHub
type ConnectionHubSpec struct {
	// +kubebuilder:validation:Enum=CloudInstance;PhysicalInstance
	InstanceType `json:"instanceType,omitempty"`
	// +kubebuilder:validation:Required
	SubmarinerSpec SubmarinerSpec `json:"submarinerSpec"`
	// +kubebuilder:validation:Required
	FederationSpec FederationOperatorSpec `json:"federationSpec"`
}

// ConnectionHubStatus defines the observed state of ConnectionHub
type ConnectionHubStatus struct {
	NodeInfo   K8sNodeInfo              `json:"nodeInfo,omitempty"`
	Phase      ConnectionHubPhase       `json:"phase,omitempty"`
	Submariner SubmarinerInstanceStatus `json:"submariner,omitempty"`
	Federation FederationInstanceStatus `json:"federation,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// ConnectionHub is the Schema for the connectionhubs API
type ConnectionHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionHubSpec   `json:"spec,omitempty"`
	Status ConnectionHubStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectionHubList contains a list of ConnectionHub
type ConnectionHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectionHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConnectionHub{}, &ConnectionHubList{})
}

func (ch *ConnectionHub) GetTenancySelectors() *Tenancy {

	tenancy := &Tenancy{}
	labels := ch.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	if physicalInstance, ok := labels[RobolaunchPhysicalInstanceLabelKey]; ok {
		tenancy.RobolaunchPhysicalInstance = physicalInstance
	}

	return tenancy
}

func (ch *ConnectionHub) GetSubmarinerMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: GlobalSubmarinerResourceName,
	}
}

func (ch *ConnectionHub) GetFederationMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: GlobalFederationOperatorResourceName,
	}
}
