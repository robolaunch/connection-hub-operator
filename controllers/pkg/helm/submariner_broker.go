package helm

import (
	"context"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"k8s.io/client-go/rest"
)

func CheckIfSubmarinerBrokerExists(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) (bool, error) {
	cli, err := getClient(config, connectionhubv1alpha1.SubmarinerBrokerNamespace)
	if err != nil {
		return false, err
	}

	_, err = cli.GetRelease(submarinerBroker.Spec.Helm.ReleaseName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func InstallSubmarinerBrokerChart(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) error {
	cli, err := getClient(config, connectionhubv1alpha1.SubmarinerBrokerNamespace)
	if err != nil {
		return err
	}

	repoName := submarinerBroker.Spec.Helm.Repository.Name
	repoURL := submarinerBroker.Spec.Helm.Repository.URL

	err = addRepository(config, connectionhubv1alpha1.SubmarinerBrokerNamespace, repoName, repoURL)
	if err != nil {
		return err
	}

	_, err = cli.InstallChart(
		context.Background(),
		&helmclient.ChartSpec{
			ReleaseName: submarinerBroker.Spec.Helm.ReleaseName,
			ChartName:   submarinerBroker.Spec.Helm.ChartName,
			Version:     submarinerBroker.Spec.Helm.Version,
		},
		&helmclient.GenericHelmOptions{},
	)

	return err
}

func UninstallSubmarinerBrokerChart(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) error {
	cli, err := getClient(config, connectionhubv1alpha1.SubmarinerBrokerNamespace)
	if err != nil {
		return err
	}

	repoName := submarinerBroker.Spec.Helm.Repository.Name
	repoURL := submarinerBroker.Spec.Helm.Repository.URL

	err = addRepository(config, connectionhubv1alpha1.SubmarinerBrokerNamespace, repoName, repoURL)
	if err != nil {
		return err
	}

	err = cli.UninstallRelease(&helmclient.ChartSpec{
		ReleaseName: submarinerBroker.Spec.Helm.ReleaseName,
		ChartName:   submarinerBroker.Spec.Helm.ChartName,
		Version:     submarinerBroker.Spec.Helm.Version,
	})

	return err
}
