package helm

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/rest"
)

func CheckIfSubmarinerOperatorExists(submarinerOperator connectionhubv1alpha1.SubmarinerOperator, config *rest.Config) (bool, error) {
	cli, err := getClient(config, submarinerOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return false, err
	}

	_, err = cli.GetRelease(submarinerOperator.Spec.HelmChart.ReleaseName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func InstallSubmarinerOperatorChart(submarinerOperator connectionhubv1alpha1.SubmarinerOperator, config *rest.Config) error {
	cli, err := getClient(config, submarinerOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := submarinerOperator.Spec.HelmRepository.Name
	repoURL := submarinerOperator.Spec.HelmRepository.URL

	err = addRepository(config, submarinerOperator.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	valuesObj := GetSubmarinerOperatorValues(submarinerOperator)

	valuesBytes, err := yaml.Marshal(&valuesObj)
	if err != nil {
		return err
	}

	_, err = cli.InstallChart(
		context.Background(),
		&helmclient.ChartSpec{
			Namespace:   submarinerOperator.GetNamespaceMetadata().Name,
			ReleaseName: submarinerOperator.Spec.HelmChart.ReleaseName,
			ChartName:   submarinerOperator.Spec.HelmChart.ChartName,
			Version:     submarinerOperator.Spec.HelmChart.Version,
			ValuesYaml:  string(valuesBytes),
			Wait:        true,
			Timeout:     time.Minute * 2,
		},
		&helmclient.GenericHelmOptions{},
	)

	return err
}

func UninstallSubmarinerOperatorChart(submarinerOperator connectionhubv1alpha1.SubmarinerOperator, config *rest.Config) error {
	cli, err := getClient(config, submarinerOperator.GetNamespaceMetadata().Name)
	if err != nil {
		return err
	}

	repoName := submarinerOperator.Spec.HelmRepository.Name
	repoURL := submarinerOperator.Spec.HelmRepository.URL

	err = addRepository(config, submarinerOperator.GetNamespaceMetadata().Name, repoName, repoURL)
	if err != nil {
		return err
	}

	err = cli.UninstallRelease(&helmclient.ChartSpec{
		ReleaseName: submarinerOperator.Spec.HelmChart.ReleaseName,
		ChartName:   submarinerOperator.Spec.HelmChart.ChartName,
		Version:     submarinerOperator.Spec.HelmChart.Version,
		Wait:        true,
		Timeout:     time.Minute * 2,
	})

	return err
}
