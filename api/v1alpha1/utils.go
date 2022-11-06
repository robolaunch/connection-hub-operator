package v1alpha1

import (
	"encoding/base64"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func (member *FederationMember) GetMemberConfig() (*rest.Config, error) {
	ca, err := decode(member.Spec.Credentials.CertificateAuthority)
	if err != nil {
		return nil, err
	}

	cert, err := decode(member.Spec.Credentials.ClientCertificate)
	if err != nil {
		return nil, err
	}

	key, err := decode(member.Spec.Credentials.ClientKey)
	if err != nil {
		return nil, err
	}

	return &rest.Config{
		Host: member.Spec.Server,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			CertData: cert,
			KeyData:  key,
		},
	}, nil
}

func (member *FederationMember) GetMemberClientset() (*kubernetes.Clientset, error) {
	memberConfig, err := member.GetMemberConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(memberConfig)
}

func (member *FederationMember) GetHostClientset(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

func decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
