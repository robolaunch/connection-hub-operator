package helm

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

// controllermanager.featureGates.RawResourceStatusCollection

type FeatureGates struct {
	RawResourceStatusCollection string `yaml:"RawResourceStatusCollection"`
}

type Controller struct {
	Repository string `yaml:"repository"`
	Image      string `yaml:"image"`
	Tag        string `yaml:"tag"`
}

type Webhook struct {
	Repository string `yaml:"repository"`
	Image      string `yaml:"image"`
	Tag        string `yaml:"tag"`
}

type PostInstallJob struct {
	Repository string `yaml:"repository"`
	Image      string `yaml:"image"`
	Tag        string `yaml:"tag"`
}

type ControllerManager struct {
	CommonNodeSelector map[string]string `yaml:"commonNodeSelector"`
	FeatureGates       FeatureGates      `yaml:"featureGates"`
	Controller         Controller        `yaml:"controller"`
	Webhook            Webhook           `yaml:"webhook"`
	PostInstallJob     PostInstallJob    `yaml:"postInstallJob"`
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

	valuesObj.ControllerManager.Controller.Repository = federationOperator.Spec.ControllerImage.Repository
	valuesObj.ControllerManager.Webhook.Repository = federationOperator.Spec.WebhookImage.Repository
	valuesObj.ControllerManager.PostInstallJob.Repository = federationOperator.Spec.PostInstallJobImage.Repository

	valuesObj.ControllerManager.Controller.Image = federationOperator.Spec.ControllerImage.Image
	valuesObj.ControllerManager.Webhook.Image = federationOperator.Spec.WebhookImage.Image
	valuesObj.ControllerManager.PostInstallJob.Image = federationOperator.Spec.PostInstallJobImage.Image

	valuesObj.ControllerManager.Controller.Tag = federationOperator.Spec.ControllerImage.Tag
	valuesObj.ControllerManager.Webhook.Tag = federationOperator.Spec.WebhookImage.Tag
	valuesObj.ControllerManager.PostInstallJob.Tag = federationOperator.Spec.PostInstallJobImage.Tag

	return valuesObj
}
