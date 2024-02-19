package helm

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

type Submariner struct {
	ClusterID           string              `yaml:"clusterId"`
	ClusterCIDR         string              `yaml:"clusterCidr"`
	ServiceCIDR         string              `yaml:"serviceCidr"`
	NatEnabled          bool                `yaml:"natEnabled"`
	ServiceDiscovery    bool                `yaml:"serviceDiscovery"`
	CableDriver         string              `yaml:"cableDriver"`
	HealthCheckEnabled  bool                `yaml:"healthcheckEnabled"`
	GlobalCIDR          string              `yaml:"globalCidr"`
	CoreDNSCustomConfig CoreDNSCustomConfig `yaml:"coreDNSCustomConfig"`
}

type CoreDNSCustomConfig struct {
	ConfigMapName string `yaml:"configMapName"`
	Namespace     string `yaml:"namespace"`
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
		GlobalCIDR:         "",
		// not a safe way to indicate CoreDNS
		// CoreDNSCustomConfig: CoreDNSCustomConfig{
		// 	ConfigMapName: "coredns-coredns",
		// 	Namespace:     "coredns",
		// },
	}
}

type Broker struct {
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	Insecure  bool   `yaml:"insecure"`
	Ca        string `yaml:"ca"`
	GlobalNet string `yaml:"globalnet"`
}

func getBrokerDefault() Broker {
	return Broker{
		Server:    "example.k8s.apiserver",
		Token:     "test",
		Namespace: "xyz",
		Insecure:  false,
		Ca:        "",
		GlobalNet: "",
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

type ServiceAccount struct {
	Create bool `yaml:"create"`
}

type ServiceAccountGlobalNet struct {
	Create string `yaml:"create"`
}

type ServiceAccounts struct {
	GlobalNet         ServiceAccountGlobalNet `yaml:"globalnet"`
	LighthouseAgent   ServiceAccount          `yaml:"lighthouseAgent"`
	LighthouseCoreDNS ServiceAccount          `yaml:"lighthouseCoreDns"`
}

func getServiceAccountsDefault() ServiceAccounts {
	return ServiceAccounts{
		GlobalNet: ServiceAccountGlobalNet{
			Create: "",
		},
		LighthouseAgent: ServiceAccount{
			Create: true,
		},
		LighthouseCoreDNS: ServiceAccount{
			Create: true,
		},
	}
}

type SubmarinerOperatorValues struct {
	Submariner      Submariner      `yaml:"submariner"`
	Broker          Broker          `yaml:"broker"`
	IPSEC           IPSEC           `yaml:"ipsec"`
	ServiceAccounts ServiceAccounts `yaml:"serviceAccounts"`
	// NodeSelector    map[string]string `yaml:"nodeSelector"`
}

func getSubmarinerOperatorValuesDefault() SubmarinerOperatorValues {
	return SubmarinerOperatorValues{
		Submariner:      getSubmarinerDefault(),
		Broker:          getBrokerDefault(),
		IPSEC:           getIPSECDefault(),
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
	valuesObj.Broker.Insecure = true
	valuesObj.Submariner.ServiceDiscovery = true
	valuesObj.Submariner.CableDriver = "wireguard"
	valuesObj.Submariner.ClusterID = submarinerOperator.Spec.ClusterID
	valuesObj.Submariner.NatEnabled = submarinerOperator.Spec.NetworkType == connectionhubv1alpha1.NetworkTypeExternal
	valuesObj.ServiceAccounts.LighthouseAgent.Create = true
	valuesObj.ServiceAccounts.LighthouseCoreDNS.Create = true
	valuesObj.Submariner.HealthCheckEnabled = true
	return valuesObj
}
