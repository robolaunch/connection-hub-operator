package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
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

type SubnetsInformation struct {
	List      []string `json:"list,omitempty"`
	ListInStr string   `json:"listInStr,omitempty"`
}

type RelayServerPodStatus struct {
	Created bool            `json:"created,omitempty"`
	Phase   corev1.PodPhase `json:"phase,omitempty"`
}

type RelayServerServiceStatus struct {
	Created bool `json:"created,omitempty"`
}

type PhysicalInstancePhase string
type PhysicalInstanceMulticastConnectionPhase string
type PhysicalInstanceFederationConnectionPhase string

const (
	PhysicalInstancePhaseLookingForDeployer  PhysicalInstancePhase = "LookingForDeployer"
	PhysicalInstancePhaseWaitingForDeployer  PhysicalInstancePhase = "WaitingForDeployer"
	PhysicalInstancePhaseRegistered          PhysicalInstancePhase = "Registered"
	PhysicalInstancePhaseCreatingRelayServer PhysicalInstancePhase = "CreatingRelayServer"
	PhysicalInstancePhaseConnected           PhysicalInstancePhase = "Connected"
	PhysicalInstancePhaseNotReady            PhysicalInstancePhase = "NotReady"
)

const (
	PhysicalInstanceMulticastConnectionPhaseWaitingForConnection PhysicalInstanceMulticastConnectionPhase = "WaitingForConnection"
	PhysicalInstanceMulticastConnectionPhaseConnecting           PhysicalInstanceMulticastConnectionPhase = "Connecting"
	PhysicalInstanceMulticastConnectionPhaseConnected            PhysicalInstanceMulticastConnectionPhase = "Connected"
	PhysicalInstanceMulticastConnectionPhaseFailed               PhysicalInstanceMulticastConnectionPhase = "Failed"
)

const (
	PhysicalInstanceFederationConnectionPhaseWaitingForMulticast   PhysicalInstanceFederationConnectionPhase = "WaitingForMulticast"
	PhysicalInstanceFederationConnectionPhaseWaitingForCredentials PhysicalInstanceFederationConnectionPhase = "WaitingForCredentials"
	PhysicalInstanceFederationConnectionPhaseConnecting            PhysicalInstanceFederationConnectionPhase = "Connecting"
	PhysicalInstanceFederationConnectionPhaseConnected             PhysicalInstanceFederationConnectionPhase = "Connected"
	PhysicalInstanceFederationConnectionPhaseFailed                PhysicalInstanceFederationConnectionPhase = "Failed"
)

// PhysicalInstanceStatus defines the observed state of PhysicalInstance
type PhysicalInstanceStatus struct {
	Submariner                SubmarinerResourceStates                  `json:"submariner,omitempty"`
	FederationMember          FederationMemberInstanceStatus            `json:"federation,omitempty"`
	Phase                     PhysicalInstancePhase                     `json:"phase,omitempty"`
	MulticastConnectionPhase  PhysicalInstanceMulticastConnectionPhase  `json:"multicastPhase,omitempty"`
	FederationConnectionPhase PhysicalInstanceFederationConnectionPhase `json:"federationPhase,omitempty"`
	Subnets                   SubnetsInformation                        `json:"subnets,omitempty"`
	RelayServerPodStatus      RelayServerPodStatus                      `json:"relayServerPodStatus,omitempty"`
	RelayServerServiceStatus  RelayServerServiceStatus                  `json:"relayServerServiceStatus,omitempty"`
	ConnectionURL             string                                    `json:"connectionURL,omitempty"`
	BootID                    string                                    `json:"bootID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Gateway",type=string,JSONPath=`.status.submariner.gatewayConnection.gatewayResource`
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.status.submariner.gatewayConnection.hostname`
//+kubebuilder:printcolumn:name="Cluster ID",type=string,JSONPath=`.status.submariner.gatewayConnection.clusterID`
//+kubebuilder:printcolumn:name="Subnets",type=string,JSONPath=`.status.subnets.list`
//+kubebuilder:printcolumn:name="Multicast",type=string,JSONPath=`.status.multicastPhase`
//+kubebuilder:printcolumn:name="Federation",type=string,JSONPath=`.status.federationPhase`
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

func (physicalinstance *PhysicalInstance) GetRelayServerPodMetadata() types.NamespacedName {
	return types.NamespacedName{
		Name:      physicalinstance.Name + "-relay",
		Namespace: "connection-hub-system",
	}
}

func (physicalinstance *PhysicalInstance) GetRelayServerServiceMetadata() types.NamespacedName {
	return types.NamespacedName{
		Name:      physicalinstance.Name + "-relay",
		Namespace: "connection-hub-system",
	}
}

func (physicalinstance *PhysicalInstance) GetConnectionHubMetadata() types.NamespacedName {
	return types.NamespacedName{
		Name: "connection-hub",
	}
}
