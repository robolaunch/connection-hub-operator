package helm

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

type CoreDNSCustomConfig struct{}

type Submariner struct {
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
		ClusterID:          "",
		ClusterCIDR:        "",
		ServiceCIDR:        "",
		NatEnabled:         true,
		ServiceDiscovery:   true,
		CableDriver:        "wireguard",
		HealthCheckEnabled: true,
	}
}

type Broker struct {
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	Insecure  bool   `yaml:"insecure"`
	Ca        string `yaml:"ca"`
}

func getBrokerDefault() Broker {
	return Broker{
		Server:    "example.k8s.apiserver",
		Token:     "test",
		Namespace: "xyz",
		Insecure:  false,
		Ca:        "",
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
	valuesObj.Submariner.CableDriver = "wireguard"
	valuesObj.Submariner.ClusterID = submarinerOperator.Spec.ClusterID
	valuesObj.Submariner.NatEnabled = submarinerOperator.Spec.NetworkType == connectionhubv1alpha1.NetworkTypeExternal
	valuesObj.ServiceAccounts.LighthouseAgent.Create = true
	valuesObj.ServiceAccounts.LighthouseCoreDNS.Create = true
	valuesObj.Submariner.HealthCheckEnabled = true
	return valuesObj
}
