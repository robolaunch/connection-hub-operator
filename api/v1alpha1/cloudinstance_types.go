package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	SubmarinerDeployerLabelKey = "robolaunch.io/submariner-deployer"
)

// CloudInstanceSpec defines the desired state of CloudInstance
type CloudInstanceSpec struct {
}

type CloudInstancePhase string

const (
	CloudInstancePhaseLookingForDeployer CloudInstancePhase = "LookingForDeployer"
	CloudInstancePhaseOwningDeployer     CloudInstancePhase = "OwningDeployer"
	CloudInstancePhaseWaitingForDeployer CloudInstancePhase = "WaitingForDeployer"
	CloudInstancePhaseTryingToConnect    CloudInstancePhase = "TryingToConnect"
	CloudInstancePhaseConnected          CloudInstancePhase = "Connected"
)

type DeployerStatus struct {
	Name   string          `json:"name,omitempty"`
	Exists bool            `json:"exists,omitempty"`
	Phase  SubmarinerPhase `json:"phase,omitempty"`
}

// CloudInstanceStatus defines the observed state of CloudInstance
type CloudInstanceStatus struct {
	DeployerStatus DeployerStatus     `json:"deployerStatus,omitempty"`
	Phase          CloudInstancePhase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// CloudInstance is the Schema for the cloudinstances API
type CloudInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudInstanceSpec   `json:"spec,omitempty"`
	Status CloudInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CloudInstanceList contains a list of CloudInstance
type CloudInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CloudInstance{}, &CloudInstanceList{})
}

func (cloudinstance *CloudInstance) GetSubmarinerDeployerMetadata() types.NamespacedName {

	return types.NamespacedName{
		Name: GlobalSubmarinerResourceName,
	}
}
