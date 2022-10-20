package helm

import (
	helm "github.com/mittwald/go-helm-client"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"
)

func getClient(config *rest.Config, namespace string) (helm.Client, error) {
	return helm.NewClientFromRestConf(&helm.RestConfClientOptions{
		Options: &helm.Options{
			Namespace: namespace,
		},
		RestConfig: config,
	})
}

func addRepository(config *rest.Config, namespace string, submarinerBroker connectionhubv1alpha1.SubmarinerBroker) error {
	cli, err := getClient(config, namespace)
	if err != nil {
		return err
	}

	// TODO: Add SSL to repository
	err = cli.AddOrUpdateChartRepo(repo.Entry{
		Name:                  "repo-name",
		URL:                   "http://ip:port",
		InsecureSkipTLSverify: true,
	})

	return err
}
