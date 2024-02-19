package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	RobolaunchCloudInstanceLabelKey      = "robolaunch.io/cloud-instance"
	RobolaunchCloudInstanceAliasLabelKey = "robolaunch.io/cloud-instance-alias"
	RobolaunchPhysicalInstanceLabelKey   = "robolaunch.io/physical-instance"
)

// Not used in robot manifest, needed for internal use.
type Tenancy struct {
	RobolaunchCloudInstance      string `json:"cloudInstance,omitempty"`
	RobolaunchCloudInstanceAlias string `json:"cloudInstanceAlias,omitempty"`
	RobolaunchPhysicalInstance   string `json:"physicalInstance,omitempty"`
}

func GetTenancyMap(obj metav1.Object) map[string]string {
	labels := obj.GetLabels()
	tenancyMap := make(map[string]string)

	if cloudInstance, ok := labels[RobolaunchCloudInstanceLabelKey]; ok {
		tenancyMap[RobolaunchCloudInstanceLabelKey] = cloudInstance
	}

	if cloudInstanceAlias, ok := labels[RobolaunchCloudInstanceAliasLabelKey]; ok {
		tenancyMap[RobolaunchCloudInstanceAliasLabelKey] = cloudInstanceAlias
	}

	if physicalInstance, ok := labels[RobolaunchPhysicalInstanceLabelKey]; ok {
		tenancyMap[RobolaunchPhysicalInstanceLabelKey] = physicalInstance
	}

	return tenancyMap
}
