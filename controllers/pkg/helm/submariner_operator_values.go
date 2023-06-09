package helm

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

type CoreDNSCustomConfig struct{}

type Images struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Submariner struct {
	DeployCR           bool   `yaml:"deployCR"`
	ClusterID          string `yaml:"clusterId"`
	ClusterCIDR        string `yaml:"clusterCidr"`
	ServiceCIDR        string `yaml:"serviceCidr"`
	NatEnabled         bool   `yaml:"natEnabled"`
	ServiceDiscovery   bool   `yaml:"serviceDiscovery"`
	CableDriver        string `yaml:"cableDriver"`
	HealthCheckEnabled bool   `yaml:"healthcheckEnabled"`
}

func getSubmarinerDefault() Submariner {
	return Submariner{
		DeployCR:           false,
		ClusterID:          "",
		ClusterCIDR:        "",
		ServiceCIDR:        "",
		NatEnabled:         false,
		ServiceDiscovery:   true,
		CableDriver:        "libreswan",
		HealthCheckEnabled: true,
	}
}

type Broker struct {
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	// Insecure  bool   `yaml:"insecure"`
	Ca string `yaml:"ca"`
	// GlobalNet bool   `yaml:"globalnet"`
}

func getBrokerDefault() Broker {
	return Broker{
		Server:    "example.k8s.apiserver",
		Token:     "test",
		Namespace: "xyz",
		// Insecure:  false,
		Ca: "",
		// GlobalNet: false,
	}
}

type IPSEC struct {
	PSK string `yaml:"psk"`
}

func getIPSECDefault() IPSEC {
	return IPSEC{
		PSK: "",
	}
}

type OperatorImage struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	// PullPolicy string `yaml:"pullPolicy"`
}

type OperatorResources struct{}

type OperatorToleration struct{}

type OperatorAffinity struct{}

type Operator struct{}

func getOperatorDefault() Operator {
	return Operator{}
}

type ServiceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

type ServiceAccounts struct {
	LighthouseAgent   ServiceAccount `yaml:"lighthouseAgent"`
	LighthouseCoreDNS ServiceAccount `yaml:"lighthouseCoreDns"`
}

func getServiceAccountsDefault() ServiceAccounts {
	return ServiceAccounts{
		LighthouseAgent: ServiceAccount{
			Create: true,
			Name:   "",
		},
		LighthouseCoreDNS: ServiceAccount{
			Create: true,
			Name:   "",
		},
	}
}

type SubmarinerOperatorValues struct {
	Submariner      Submariner        `yaml:"submariner"`
	Broker          Broker            `yaml:"broker"`
	IPSEC           IPSEC             `yaml:"ipsec"`
	Operator        Operator          `yaml:"operator"`
	ServiceAccounts ServiceAccounts   `yaml:"serviceAccounts"`
	NodeSelector    map[string]string `yaml:"nodeSelector"`
}

func getSubmarinerOperatorValuesDefault() SubmarinerOperatorValues {
	return SubmarinerOperatorValues{
		Submariner:      getSubmarinerDefault(),
		Broker:          getBrokerDefault(),
		IPSEC:           getIPSECDefault(),
		Operator:        getOperatorDefault(),
		ServiceAccounts: getServiceAccountsDefault(),
	}
}

func GetSubmarinerOperatorValues(submarinerOperator connectionhubv1alpha1.SubmarinerOperator) SubmarinerOperatorValues {
	valuesObj := getSubmarinerOperatorValuesDefault()
	valuesObj.Submariner.ClusterCIDR = submarinerOperator.Spec.ClusterCIDR
	valuesObj.Submariner.ServiceCIDR = submarinerOperator.Spec.ServiceCIDR
	valuesObj.IPSEC.PSK = submarinerOperator.Spec.PresharedKey
	valuesObj.Broker.Namespace = connectionhubv1alpha1.SubmarinerBrokerNamespace
	valuesObj.Broker.Server = submarinerOperator.Spec.APIServerURL
	valuesObj.Broker.Token = submarinerOperator.Spec.BrokerCredentials.Token
	valuesObj.Broker.Ca = submarinerOperator.Spec.BrokerCredentials.CA
	valuesObj.Submariner.ServiceDiscovery = true
	valuesObj.Submariner.CableDriver = "libreswan"
	valuesObj.Submariner.ClusterID = submarinerOperator.Spec.ClusterID
	valuesObj.Submariner.NatEnabled = submarinerOperator.Spec.NetworkType == connectionhubv1alpha1.NetworkTypeExternal
	valuesObj.ServiceAccounts.LighthouseAgent.Create = true
	valuesObj.ServiceAccounts.LighthouseCoreDNS.Create = true
	valuesObj.Submariner.HealthCheckEnabled = true
	return valuesObj
}
