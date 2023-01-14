package helm

import connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"

// controllermanager.featureGates.RawResourceStatusCollection

type FeatureGates struct {
	RawResourceStatusCollection string `yaml:"RawResourceStatusCollection"`
}

type ControllerManager struct {
	CommonNodeSelector map[string]string `yaml:"commonNodeSelector"`
	FeatureGates       FeatureGates      `yaml:"featureGates"`
}

type FederationOperatorValues struct {
	ControllerManager ControllerManager `yaml:"controllermanager"`
}

func getFederationOperatorValuesDefault() FederationOperatorValues {
	return FederationOperatorValues{
		ControllerManager: ControllerManager{
			CommonNodeSelector: make(map[string]string),
		},
	}
}

func GetFederationOperatorValues(federationOperator connectionhubv1alpha1.FederationOperator) FederationOperatorValues {
	valuesObj := getFederationOperatorValuesDefault()
	valuesObj.ControllerManager.CommonNodeSelector = federationOperator.Labels
	valuesObj.ControllerManager.FeatureGates.RawResourceStatusCollection = "Enabled"

	return valuesObj
}
