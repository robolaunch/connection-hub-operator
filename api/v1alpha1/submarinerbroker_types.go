package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// SubmarinerBrokerSpec defines the desired state of SubmarinerBroker
type SubmarinerBrokerSpec struct {
	// +kubebuilder:validation:Required
	HelmRepository HelmRepository `json:"helmRepository"`
	HelmChart      HelmChart      `json:"helmChart"`
	// +kubebuilder:validation:Required
	APIServerURL string `json:"apiServerURL"`
}

type K8sNodeInfo struct {
	Name      string            `json:"name,omitempty"`
	Selectors map[string]string `json:"selectors,omitempty"`
}

type SubmarinerBrokerPhase string

const (
	SubmarinerBrokerPhaseCreatingNamespace SubmarinerBrokerPhase = "CreatingNamespace"
	SubmarinerBrokerPhaseDeployingChart    SubmarinerBrokerPhase = "DeployingChart"
	SubmarinerBrokerPhaseCheckingResources SubmarinerBrokerPhase = "CheckingResources"
	SubmarinerBrokerPhaseDeployed          SubmarinerBrokerPhase = "Deployed"
	SubmarinerBrokerPhaseMalfunctioned     SubmarinerBrokerPhase = "Malfunctioned"

	SubmarinerBrokerPhaseUninstallingChart    SubmarinerBrokerPhase = "UninstallingChart"
	SubmarinerBrokerPhaseTerminatingNamespace SubmarinerBrokerPhase = "TerminatingNamespace"
)

const (
	SubmarinerBrokerNamespace string = "submariner-k8s-broker"
)

type BrokerCredentials struct {
	Token string `json:"token,omitempty"`
	CA    string `json:"ca,omitempty"`
}

// SubmarinerBrokerStatus defines the observed state of SubmarinerBroker
type SubmarinerBrokerStatus struct {
	NamespaceStatus     NamespaceStatus       `json:"namespaceStatus,omitempty"`
	ChartStatus         ChartStatus           `json:"chartStatus,omitempty"`
	ChartResourceStatus ChartResourceStatus   `json:"chartResourceStatus,omitempty"`
	Phase               SubmarinerBrokerPhase `json:"phase,omitempty"`
	NodeInfo            K8sNodeInfo           `json:"nodeInfo,omitempty"`
	BrokerCredentials   BrokerCredentials     `json:"brokerCredentials,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// SubmarinerBroker is the Schema for the submarinerbrokers API
type SubmarinerBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubmarinerBrokerSpec   `json:"spec,omitempty"`
	Status SubmarinerBrokerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubmarinerBrokerList contains a list of SubmarinerBroker
type SubmarinerBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubmarinerBroker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubmarinerBroker{}, &SubmarinerBrokerList{})
}

func (smb *SubmarinerBroker) GetTenancySelectors() *Tenancy {

	tenancy := &Tenancy{}
	labels := smb.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	if cloudInstanceAlias, ok := labels[tenancy.RobolaunchCloudInstanceAlias]; ok {
		tenancy.RobolaunchCloudInstanceAlias = cloudInstanceAlias
	}

	return tenancy
}

func (smb *SubmarinerBroker) GetNamespaceMetadata() *types.NamespacedName {
	return &types.NamespacedName{
		Name: SubmarinerBrokerNamespace,
	}
}
