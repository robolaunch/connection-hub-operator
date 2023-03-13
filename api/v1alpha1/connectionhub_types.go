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

type HelmRepository struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

type HelmChart struct {
	// +kubebuilder:validation:Required
	ReleaseName string `json:"releaseName"`
	// +kubebuilder:validation:Required
	ChartName string `json:"chartName"`
	// +kubebuilder:validation:Required
	Version string `json:"version"`
}

type SubmarinerInstanceStatus struct {
	Created bool            `json:"created,omitempty"`
	Phase   SubmarinerPhase `json:"phase,omitempty"`
}

type FederationInstanceStatus struct {
	Created bool                    `json:"created,omitempty"`
	Phase   FederationOperatorPhase `json:"phase,omitempty"`
}

type FederationHostInstanceStatus struct {
	Created bool                `json:"created,omitempty"`
	Phase   FederationHostPhase `json:"phase,omitempty"`
}

type CloudInstanceInstanceStatus struct {
	Created bool               `json:"created,omitempty"`
	Phase   CloudInstancePhase `json:"phase,omitempty"`
}

type ConnectionHubPhase string

const (
	ConnectionHubPhaseSubmarinerSettingUp    ConnectionHubPhase = "SubmarinerSettingUp"
	ConnectionHubPhaseFederationSettingUp    ConnectionHubPhase = "FederationSettingUp"
	ConnectionHubPhaseCreatingFederationHost ConnectionHubPhase = "CreatingFederationHost"
	ConnectionHubPhaseCreatingCloudInstance  ConnectionHubPhase = "CreatingCloudInstance"
	ConnectionHubPhaseReadyForOperation      ConnectionHubPhase = "ReadyForOperation"

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

type ConnectionInterfaces struct {
	ForPhysicalInstance ConnectionHubSpec               `json:"forPhysicalInstance,omitempty"`
	ForCloudInstance    map[string]FederationMemberSpec `json:"forCloudInstance,omitempty"`
}

// ConnectionHubStatus defines the observed state of ConnectionHub
type ConnectionHubStatus struct {
	NodeInfo   K8sNodeInfo              `json:"nodeInfo,omitempty"`
	Phase      ConnectionHubPhase       `json:"phase,omitempty"`
	Submariner SubmarinerInstanceStatus `json:"submariner,omitempty"`
	Federation FederationInstanceStatus `json:"federation,omitempty"`

	FederationHost FederationHostInstanceStatus `json:"federationHost,omitempty"`
	CloudInstance  CloudInstanceInstanceStatus  `json:"cloudInstance,omitempty"`

	ConnectionInterfaces ConnectionInterfaces `json:"connectionInterfaces,omitempty"`
	Key                  string               `json:"key,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`

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

	if cloudInstanceAlias, ok := labels[tenancy.RobolaunchCloudInstanceAlias]; ok {
		tenancy.RobolaunchCloudInstanceAlias = cloudInstanceAlias
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

func (ch *ConnectionHub) GetFederationHostMetadata() *types.NamespacedName {

	tenancy := ch.GetTenancySelectors()

	return &types.NamespacedName{
		Name: tenancy.RobolaunchCloudInstance,
	}
}

func (ch *ConnectionHub) GetCloudInstanceMetadata() *types.NamespacedName {

	tenancy := ch.GetTenancySelectors()

	return &types.NamespacedName{
		Name: tenancy.RobolaunchCloudInstance,
	}
}
