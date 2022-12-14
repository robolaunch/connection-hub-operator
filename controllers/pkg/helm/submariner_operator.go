package helm

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/rest"
)

/*
helm install submariner-operator  ./submariner-operator \
--create-namespace --namespace "${SUBMARINER_NS}"  \
--set submariner.clusterCidr="${CLUSTER_CIDR}" \
--set submariner.serviceCidr="${SERVICE_CIDR}" \
--set ipsec.psk="${SUBMARINER_PSK}" \
--set broker.server="${SUBMARINER_BROKER_URL}" \
--set broker.token="${SUBMARINER_BROKER_TOKEN}" \
--set broker.namespace="${BROKER_NS}" \
--set broker.ca="${SUBMARINER_BROKER_CA}" \
--set submariner.serviceDiscovery=true \
--set submariner.cableDriver=wireguard \
--set submariner.clusterId="${CLUSTER_ID}" \
--set submariner.natEnabled="true" \
--set serviceAccounts.lighthouseAgent.create=true \
--set serviceAccounts.lighthouseCoreDns.create=true \
--set submariner.healthcheckEnabled=false \
--set ipsec.natPort=4500 \
--set ipsec.ikePort=500 \
--set ipsec.preferredServer="true" \
--set ipsec.natDiscovery=4490 \
--set gateway.image.repository="docker.io/robolaunchio/submariner-gateway" \
--set gateway.image.tag="dev-v11" \
--set operator.image.repository="docker.io/robolaunchio/submariner-operator" \
--set operator.image.tag="dev-v14" \
--set submariner.images.repository="docker.io/robolaunchio" \
--set submariner.images.tag="dev-v11"
*/

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
