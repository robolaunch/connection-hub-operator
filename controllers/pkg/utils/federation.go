package utils

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/kubefed/pkg/kubefedctl"
)

func JoinMember(member *connectionhubv1alpha1.FederationMember, hostConfig *rest.Config) error {

	var memberConfig *rest.Config
	var err error

	if member.Status.Host.Name == member.Name {
		memberConfig = hostConfig
	} else {
		memberConfig, err = member.GetMemberConfig()
		if err != nil {
			return err
		}
	}

	_, err = kubefedctl.JoinCluster(
		hostConfig,
		memberConfig,
		connectionhubv1alpha1.FederationOperatorNamespace,
		member.Status.Host.Name,
		member.Name,
		"",
		v1.ClusterScoped,
		false,
		false,
	)

	if err != nil {
		return err
	}

	return nil
}
