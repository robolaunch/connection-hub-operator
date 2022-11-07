package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PhysicalInstanceCredentials struct {
	// +kubebuilder:validation:Required
	CertificateAuthority string `json:"certificateAuthority"`
	// +kubebuilder:validation:Required
	ClientCertificate string `json:"clientCertificate"`
	// +kubebuilder:validation:Required
	ClientKey string `json:"clientKey"`
}

// PhysicalInstanceSpec defines the desired state of PhysicalInstance
type PhysicalInstanceSpec struct {
	// +kubebuilder:validation:Required
	Server string `json:"server"`
	// +kubebuilder:validation:Required
	Credentials PhysicalInstanceCredentials `json:"credentials"`
}

type SubmarinerResourceStates struct {
	DeployerStatus      DeployerStatus             `json:"deployerStatus,omitempty"`
	ConnectionResources ConnectionResourceStatuses `json:"connectionResources,omitempty"`
	GatewayConnection   GatewayConnection          `json:"gatewayConnection,omitempty"`
}

type FederationMemberInstanceStatus struct {
	Created bool                   `json:"created,omitempty"`
	Status  FederationMemberStatus `json:"status,omitempty"`
}

type PhysicalInstancePhase string

const (
	PhysicalInstancePhaseLookingForDeployer        PhysicalInstancePhase = "LookingForDeployer"
	PhysicalInstancePhaseWaitingForDeployer        PhysicalInstancePhase = "WaitingForDeployer"
	PhysicalInstancePhaseRegistered                PhysicalInstancePhase = "Registered"
	PhysicalInstancePhaseConnectingOverMulticast   PhysicalInstancePhase = "ConnectingOverMulticast"
	PhysicalInstancePhaseConnectingOverKubernetes  PhysicalInstancePhase = "ConnectingOverKubernetes"
	PhysicalInstancePhaseConnected                 PhysicalInstancePhase = "Connected"
	PhysicalInstancePhaseNotConnectedOverMulticast PhysicalInstancePhase = "NotConnectedOverMulticast"
)

// PhysicalInstanceStatus defines the observed state of PhysicalInstance
type PhysicalInstanceStatus struct {
	Submariner       SubmarinerResourceStates       `json:"submariner,omitempty"`
	FederationMember FederationMemberInstanceStatus `json:"federation,omitempty"`
	Phase            PhysicalInstancePhase          `json:"phase,omitempty"`
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
