package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SubmarinerDeployerLabelKey = "robolaunch.io/submariner-deployer"
)

// CloudInstanceSpec defines the desired state of CloudInstance
type CloudInstanceSpec struct {
}

type CloudInstancePhase string

const (
	CloudInstancePhaseDeployerNotFound CloudInstancePhase = "DeployerNotFound"
	CloudInstancePhaseDeployerNotReady CloudInstancePhase = "DeployerNotReady"
	CloudInstancePhaseTryingToConnect  CloudInstancePhase = "TryingToConnect"
	CloudInstancePhaseConnected        CloudInstancePhase = "Connected"
)

// CloudInstanceStatus defines the observed state of CloudInstance
type CloudInstanceStatus struct {
	Phase CloudInstancePhase `json:"phase,omitempty"`
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
