package v1alpha1

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
