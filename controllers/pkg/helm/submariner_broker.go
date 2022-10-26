package helm

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"k8s.io/client-go/rest"
)

func CheckIfSubmarinerBrokerExists(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) (bool, error) {
	cli, err := getClient(config, submarinerBroker.GetNamespaceMetadata().Name)
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
	cli, err := getClient(config, submarinerBroker.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := submarinerBroker.Spec.Helm.Repository.Name
	repoURL := submarinerBroker.Spec.Helm.Repository.URL

	err = addRepository(config, submarinerBroker.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	_, err = cli.InstallChart(
		context.Background(),
		&helmclient.ChartSpec{
			ReleaseName: submarinerBroker.Spec.Helm.ReleaseName,
			ChartName:   submarinerBroker.Spec.Helm.ChartName,
			Version:     submarinerBroker.Spec.Helm.Version,
			Wait:        true,
			Timeout:     time.Minute * 2,
		},
		&helmclient.GenericHelmOptions{},
	)

	time.Sleep(2 * time.Second)

	return err
}

func UninstallSubmarinerBrokerChart(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) error {
	cli, err := getClient(config, submarinerBroker.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := submarinerBroker.Spec.Helm.Repository.Name
	repoURL := submarinerBroker.Spec.Helm.Repository.URL

	err = addRepository(config, submarinerBroker.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	err = cli.UninstallRelease(&helmclient.ChartSpec{
		ReleaseName: submarinerBroker.Spec.Helm.ReleaseName,
		ChartName:   submarinerBroker.Spec.Helm.ChartName,
		Version:     submarinerBroker.Spec.Helm.Version,
		Wait:        true,
		Timeout:     time.Minute * 2,
	})

	return err
}
