package v1alpha1

import (
	brokerv1 "github.com/robolaunch/connection-hub-operator/api/external/submariner/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	EndpointClusterIDLabelKey = "submariner-io/clusterID"
)

// CloudInstanceSpec defines the desired state of CloudInstance
type CloudInstanceSpec struct {
}

type CloudInstancePhase string

const (
	CloudInstancePhaseLookingForDeployer  CloudInstancePhase = "LookingForDeployer"
	CloudInstancePhaseOwningDeployer      CloudInstancePhase = "OwningDeployer"
	CloudInstancePhaseWaitingForDeployer  CloudInstancePhase = "WaitingForDeployer"
	CloudInstancePhaseWaitingForResources CloudInstancePhase = "WaitingForResources"
	CloudInstancePhaseConnecting          CloudInstancePhase = "Connecting"
	CloudInstancePhaseConnected           CloudInstancePhase = "Connected"
	CloudInstancePhaseNotConnected        CloudInstancePhase = "NotConnected"
)

type DeployerStatus struct {
	Name   string          `json:"name,omitempty"`
	Exists bool            `json:"exists,omitempty"`
	Phase  SubmarinerPhase `json:"phase,omitempty"`
}

type ConnectionResourceStatus struct {
	Name   string `json:"name,omitempty"`
	Exists bool   `json:"exists,omitempty"`
}

type GatewayConnection struct {
	GatewayResource  string                    `json:"gatewayResource,omitempty"`
	ClusterID        string                    `json:"clusterID,omitempty"`
	Hostname         string                    `json:"hostname,omitempty"`
	ConnectionStatus brokerv1.ConnectionStatus `json:"connectionStatus,omitempty"`
}

type ConnectionResourceStatuses struct {
	ClusterStatus  ConnectionResourceStatus `json:"clusterStatus,omitempty"`
	EndpointStatus ConnectionResourceStatus `json:"endpointStatus,omitempty"`
}

// CloudInstanceStatus defines the observed state of CloudInstance
type CloudInstanceStatus struct {
	DeployerStatus      DeployerStatus             `json:"deployerStatus,omitempty"`
	ConnectionResources ConnectionResourceStatuses `json:"connectionResources,omitempty"`
	GatewayConnection   GatewayConnection          `json:"gatewayConnection,omitempty"`
	Phase               CloudInstancePhase         `json:"phase,omitempty"`
	BootID              string                     `json:"bootID,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Gateway",type=string,JSONPath=`.status.gatewayConnection.gatewayResource`
//+kubebuilder:printcolumn:name="Hostname",type=string,JSONPath=`.status.gatewayConnection.hostname`
//+kubebuilder:printcolumn:name="Cluster ID",type=string,JSONPath=`.status.gatewayConnection.clusterID`
//+kubebuilder:printcolumn:name="Multicast",type=string,JSONPath=`.status.phase`

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

func (cloudinstance *CloudInstance) GetSubmarinerClusterMetadata() types.NamespacedName {

	return types.NamespacedName{
		Name:      cloudinstance.Name,
		Namespace: SubmarinerOperatorNamespace,
	}
}
