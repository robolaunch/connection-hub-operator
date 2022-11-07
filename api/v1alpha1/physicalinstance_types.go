package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PhysicalInstanceSpec defines the desired state of PhysicalInstance
type PhysicalInstanceSpec struct {
}

type SubmarinerResourceStates struct {
	DeployerStatus      DeployerStatus             `json:"deployerStatus,omitempty"`
	ConnectionResources ConnectionResourceStatuses `json:"connectionResources,omitempty"`
	GatewayConnection   GatewayConnection          `json:"gatewayConnection,omitempty"`
}

type PhysicalInstancePhase string

const (
	PhysicalInstancePhaseLookingForDeployer PhysicalInstancePhase = "LookingForDeployer"
	PhysicalInstancePhaseWaitingForDeployer PhysicalInstancePhase = "WaitingForDeployer"
	PhysicalInstancePhaseRegistered         PhysicalInstancePhase = "Registered"
	PhysicalInstancePhaseConnecting         PhysicalInstancePhase = "Connecting"
	PhysicalInstancePhaseConnected          PhysicalInstancePhase = "Connected"
	PhysicalInstancePhaseNotConnected       PhysicalInstancePhase = "NotConnected"
)

// PhysicalInstanceStatus defines the observed state of PhysicalInstance
type PhysicalInstanceStatus struct {
	Submariner SubmarinerResourceStates `json:"submariner,omitempty"`
	Phase      PhysicalInstancePhase    `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Gateway",type=string,JSONPath=`.status.gatewayConnection.gatewayResource`
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.status.gatewayConnection.hostname`
//+kubebuilder:printcolumn:name="Cluster ID",type=string,JSONPath=`.status.gatewayConnection.clusterID`
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

// PhysicalInstance is the Schema for the physicalinstances API
type PhysicalInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PhysicalInstanceSpec   `json:"spec,omitempty"`
	Status PhysicalInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PhysicalInstanceList contains a list of PhysicalInstance
type PhysicalInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PhysicalInstance{}, &PhysicalInstanceList{})
}

func (physicalinstance *PhysicalInstance) GetSubmarinerDeployerMetadata() types.NamespacedName {

	return types.NamespacedName{
		Name: GlobalSubmarinerResourceName,
	}
}

func (physicalinstance *PhysicalInstance) GetSubmarinerClusterMetadata() types.NamespacedName {

	return types.NamespacedName{
		Name:      physicalinstance.Name,
		Namespace: SubmarinerOperatorNamespace,
	}
}
