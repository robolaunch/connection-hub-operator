package helm

import (
	"context"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
)

func InstallSubmarinerBrokerChart(submarinerBroker connectionhubv1alpha1.SubmarinerBroker, config *rest.Config) error {
	cli, err := getClient(config, connectionhubv1alpha1.SubmarinerBrokerNamespace)
	if err != nil {
		return err
	}

	err = addRepository(config, connectionhubv1alpha1.SubmarinerBrokerNamespace, submarinerBroker)
	if err != nil {
		// TODO: Check Helm client error types
		if !errors.IsAlreadyExists(err) {
			return err
		}
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
