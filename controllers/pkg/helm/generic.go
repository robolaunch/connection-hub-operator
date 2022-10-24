package helm

import (
	helm "github.com/mittwald/go-helm-client"
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

func addRepository(config *rest.Config, namespace string, repoName string, repoURL string) error {
	cli, err := getClient(config, namespace)
	if err != nil {
		return err
	}

	// TODO: Add SSL to repository
	err = cli.AddOrUpdateChartRepo(repo.Entry{
		Name:                  repoName,
		URL:                   repoURL,
		InsecureSkipTLSverify: true,
	})

	return err
}
