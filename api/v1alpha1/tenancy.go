package v1alpha1

const (
	RobolaunchCloudInstanceLabelKey    = "robolaunch.io/cloud-instance"
	RobolaunchPhysicalInstanceLabelKey = "robolaunch.io/physical-instance"
)

// Not used in robot manifest, needed for internal use.
type Tenancy struct {
	RobolaunchCloudInstance    string `json:"cloudInstance,omitempty"`
	RobolaunchPhysicalInstance string `json:"physicalInstance,omitempty"`
}
