package helm

type CoreDNSCustomConfig struct{}

type Images struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type Submariner struct {
	ClusterID           string              `yaml:"clusterId"`
	Token               string              `yaml:"token"`
	ClusterCIDR         string              `yaml:"clusterCidr"`
	ServiceCIDR         string              `yaml:"serviceCidr"`
	GlobalCIDR          string              `yaml:"globalCidr"`
	NotEnabled          bool                `yaml:"natEnabled"`
	ColorCodes          string              `yaml:"colorCodes"`
	Debug               bool                `yaml:"debug"`
	ServiceDiscovery    bool                `yaml:"serviceDiscovery"`
	CableDriver         string              `yaml:"cableDriver"`
	HealthCheckEnabled  bool                `yaml:"healthcheckEnabled"`
	CoreDNSCustomConfig CoreDNSCustomConfig `yaml:"coreDNSCustomConfig"`
	Images              Images              `yaml:"images"`
}

func getSubmarinerDefault() Submariner {
	return Submariner{
		ClusterID:           "",
		Token:               "",
		ClusterCIDR:         "",
		ServiceCIDR:         "",
		GlobalCIDR:          "",
		NotEnabled:          false,
		ColorCodes:          "blue",
		Debug:               false,
		ServiceDiscovery:    true,
		CableDriver:         "libreswan",
		HealthCheckEnabled:  true,
		CoreDNSCustomConfig: CoreDNSCustomConfig{},
		Images: Images{
			Repository: "quay.io/submariner",
			Tag:        "0.10.1",
		},
	}
}

type Broker struct {
	Server    string `yaml:"server"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	Insecure  bool   `yaml:"insecure"`
	Ca        string `yaml:"ca"`
	GlobalNet bool   `yaml:"globalnet"`
}

func getBrokerDefault() Broker {
	return Broker{
		Server:    "example.k8s.apiserver",
		Token:     "test",
		Namespace: "xyz",
		Insecure:  false,
		Ca:        "",
		GlobalNet: false,
	}
}

type RBAC struct {
	Create bool `yaml:"create"`
}

func getRBACDefault() RBAC {
	return RBAC{
		Create: true,
	}
}

type IPSEC struct {
	PSK             string `yaml:"psk"`
	Debug           bool   `yaml:"debug"`
	ForceUDPEncaps  bool   `yaml:"forceUDPEncaps"`
	IKEPort         int    `yaml:"ikePort"`
	NATPort         int    `yaml:"natPort"`
	NATDiscovery    int    `yaml:"natDiscovery"`
	PreferredServer bool   `yaml:"preferredServer"`
}

func getIPSECDefault() IPSEC {
	return IPSEC{
		PSK:             "",
		Debug:           false,
		ForceUDPEncaps:  false,
		IKEPort:         500,
		NATPort:         4500,
		NATDiscovery:    4490,
		PreferredServer: false,
	}
}

type Leadership struct {
	LeaseDuration int `yaml:"leaseDuration"`
	RenewDeadline int `yaml:"renewDeadline"`
	RetryPeriod   int `yaml:"retryPeriod"`
}

func getLeadershipDefault() Leadership {
	return Leadership{
		LeaseDuration: 10,
		RenewDeadline: 5,
		RetryPeriod:   2,
	}
}

type OperatorImage struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
	PullPolicy string `yaml:"pullPolicy"`
}

type OperatorResources struct{}

type OperatorToleration struct{}

type OperatorAffinity struct{}

type Operator struct {
	Image       OperatorImage        `yaml:"image"`
	Resources   OperatorResources    `yaml:"resources"`
	Tolerations []OperatorToleration `yaml:"tolerations"`
	Affinity    OperatorAffinity     `yaml:"affinity"`
}

func getOperatorDefault() Operator {
	return Operator{
		Image: OperatorImage{
			Repository: "quay.io/submariner/submariner-operator",
			Tag:        "0.10.1",
			PullPolicy: "IfNotPresent",
		},
		Resources:   OperatorResources{},
		Tolerations: []OperatorToleration{},
		Affinity:    OperatorAffinity{},
	}
}

type Gateway struct {
	Image Images `yaml:"image"`
}

func getOperatorGatewayDefault() Gateway {
	return Gateway{
		Image: Images{
			Repository: "quay.io/submariner/submariner-gateway",
			Tag:        "0.10.1",
		},
	}
}

type ServiceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

type ServiceAccounts struct {
	Operator          ServiceAccount `yaml:"operator"`
	Gateway           ServiceAccount `yaml:"gateway"`
	RouteAgent        ServiceAccount `yaml:"routeAgent"`
	GlobalNet         ServiceAccount `yaml:"globalnet"`
	LighthouseAgent   ServiceAccount `yaml:"lighthouseAgent"`
	LighthouseCoreDNS ServiceAccount `yaml:"lighthouseCoreDns"`
}

func getServiceAccountsDefault() ServiceAccounts {
	return ServiceAccounts{
		Operator: ServiceAccount{
			Create: true,
			Name:   "",
		},
		Gateway: ServiceAccount{
			Create: true,
			Name:   "",
		},
		RouteAgent: ServiceAccount{
			Create: true,
			Name:   "",
		},
		GlobalNet: ServiceAccount{
			Create: true,
			Name:   "",
		},
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
	Submariner      Submariner      `yaml:"submariner"`
	Broker          Broker          `yaml:"broker"`
	RBAC            RBAC            `yaml:"rbac"`
	IPSEC           IPSEC           `yaml:"ipsec"`
	Leadership      Leadership      `yaml:"leadership"`
	Operator        Operator        `yaml:"operator"`
	Gateway         Gateway         `yaml:"gateway"`
	ServiceAccounts ServiceAccounts `yaml:"serviceAccounts"`
}

func GetSubmarinerOperatorValuesDefault() SubmarinerOperatorValues {
	return SubmarinerOperatorValues{
		Submariner:      getSubmarinerDefault(),
		Broker:          getBrokerDefault(),
		RBAC:            getRBACDefault(),
		IPSEC:           getIPSECDefault(),
		Leadership:      getLeadershipDefault(),
		Operator:        getOperatorDefault(),
		Gateway:         getOperatorGatewayDefault(),
		ServiceAccounts: getServiceAccountsDefault(),
	}
}
