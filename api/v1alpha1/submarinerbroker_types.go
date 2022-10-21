package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HelmRepository struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

type HelmChartProperties struct {
	// +kubebuilder:validation:Required
	Repository HelmRepository `json:"repository"`
	// +kubebuilder:validation:Required
	ReleaseName string `json:"releaseName"`
	// +kubebuilder:validation:Required
	ChartName string `json:"chartName"`
	// +kubebuilder:validation:Required
	Version string `json:"version"`
	Values  string `json:"values,omitempty"`
}

// SubmarinerBrokerSpec defines the desired state of SubmarinerBroker
type SubmarinerBrokerSpec struct {
	// +kubebuilder:validation:Required
	Helm HelmChartProperties `json:"helm"`
}

type K8sNodeInfo struct {
	Name      string            `json:"name,omitempty"`
	Selectors map[string]string `json:"selectors,omitempty"`
}

type SubmarinerBrokerPhase string

const (
	SubmarinerBrokerPhaseNotExists     SubmarinerBrokerPhase = "NotExists"
	SubmarinerBrokerPhaseDeploying     SubmarinerBrokerPhase = "Deploying"
	SubmarinerBrokerPhaseRunning       SubmarinerBrokerPhase = "Running"
	SubmarinerBrokerPhaseMalfunctioned SubmarinerBrokerPhase = "Malfunctioned"
)

const (
	SubmarinerBrokerNamespace string = "submariner-k8s-broker"
)

// SubmarinerBrokerStatus defines the observed state of SubmarinerBroker
type SubmarinerBrokerStatus struct {
	// +kubebuilder:default="NotExists"
	Phase       SubmarinerBrokerPhase `json:"phase,omitempty"`
	NodeInfo    K8sNodeInfo           `json:"nodeInfo,omitempty"`
	BrokerURL   string                `json:"brokerURL,omitempty"`
	BrokerToken string                `json:"brokerToken,omitempty"`
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

func GetTenancySelectorsForSMB(submarinerBroker SubmarinerBroker) *Tenancy {

	tenancy := &Tenancy{}
	labels := submarinerBroker.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	return tenancy
}