package helm

import connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"

type ControllerManager struct {
	CommonNodeSelector map[string]string `yaml:"commonNodeSelector"`
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

	return valuesObj
}
