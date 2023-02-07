package helm

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

// controllermanager.featureGates.RawResourceStatusCollection

type FeatureGates struct {
	RawResourceStatusCollection string `yaml:"RawResourceStatusCollection"`
}

type Controller struct {
	Image string `yaml:"image"`
}

type Webhook struct {
	Image string `yaml:"image"`
}

type ControllerManager struct {
	CommonNodeSelector map[string]string `yaml:"commonNodeSelector"`
	FeatureGates       FeatureGates      `yaml:"featureGates"`
	Controller         Controller        `yaml:"controller"`
	Webhook            Webhook           `yaml:"webhook"`
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
	valuesObj.ControllerManager.Controller.Image = federationOperator.Spec.ControllerImage.Repository + federationOperator.Spec.ControllerImage.Tag
	valuesObj.ControllerManager.Webhook.Image = federationOperator.Spec.WebhookImage.Repository + federationOperator.Spec.WebhookImage.Tag

	return valuesObj
}
