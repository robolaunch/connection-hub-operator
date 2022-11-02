package helm

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"k8s.io/client-go/rest"
)

func CheckIfFederationOperatorExists(federationOperator connectionhubv1alpha1.FederationOperator, config *rest.Config) (bool, error) {
	cli, err := getClient(config, federationOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return false, err
	}

	_, err = cli.GetRelease(federationOperator.Spec.Helm.ReleaseName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func InstallFederationOperatorChart(federationOperator connectionhubv1alpha1.FederationOperator, config *rest.Config) error {
	cli, err := getClient(config, federationOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := federationOperator.Spec.Helm.Repository.Name
	repoURL := federationOperator.Spec.Helm.Repository.URL

	err = addRepository(config, federationOperator.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	_, err = cli.InstallChart(
		context.Background(),
		&helmclient.ChartSpec{
			Namespace:   federationOperator.GetNamespaceMetadata().Name,
			ReleaseName: federationOperator.Spec.Helm.ReleaseName,
			ChartName:   federationOperator.Spec.Helm.ChartName,
			Version:     federationOperator.Spec.Helm.Version,
			Wait:        true,
			Timeout:     time.Minute * 2,
		},
		&helmclient.GenericHelmOptions{},
	)

	return err
}

func UninstallFederationOperatorChart(federationOperator connectionhubv1alpha1.FederationOperator, config *rest.Config) error {
	cli, err := getClient(config, federationOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := federationOperator.Spec.Helm.Repository.Name
	repoURL := federationOperator.Spec.Helm.Repository.URL

	err = addRepository(config, federationOperator.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	err = cli.UninstallRelease(&helmclient.ChartSpec{
		ReleaseName: federationOperator.Spec.Helm.ReleaseName,
		ChartName:   federationOperator.Spec.Helm.ChartName,
		Version:     federationOperator.Spec.Helm.Version,
		Wait:        true,
		Timeout:     time.Minute * 2,
	})

	return err
}
