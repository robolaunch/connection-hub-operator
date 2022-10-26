package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubmarinerSpec defines the desired state of Submariner
type SubmarinerSpec struct {
}

// SubmarinerStatus defines the observed state of Submariner
type SubmarinerStatus struct {
	NodeInfo K8sNodeInfo `json:"nodeInfo,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// Submariner is the Schema for the submariners API
type Submariner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubmarinerSpec   `json:"spec,omitempty"`
	Status SubmarinerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubmarinerList contains a list of Submariner
type SubmarinerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Submariner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Submariner{}, &SubmarinerList{})
}

func GetTenancySelectorsForSubmariner(submariner Submariner) *Tenancy {

	tenancy := &Tenancy{}
	labels := submariner.GetLabels()

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancy.RobolaunchCloudInstance = cloudInstance
	}

	if physicalInstance, ok := labels[RobolaunchPhysicalInstanceLabelKey]; ok {
		tenancy.RobolaunchPhysicalInstance = physicalInstance
	}

	return tenancy
}
