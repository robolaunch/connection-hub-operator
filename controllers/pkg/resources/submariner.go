package resources

import (
	submv1alpha1 "github.com/robolaunch/connection-hub-operator/api/external/submariner/v1alpha1"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSubmarinerBroker(cr *connectionhubv1alpha1.Submariner) *connectionhubv1alpha1.SubmarinerBroker {

	brokerSpec := connectionhubv1alpha1.SubmarinerBrokerSpec{
		Helm:         cr.Spec.BrokerHelmChart,
		APIServerURL: cr.Spec.APIServerURL,
	}

	broker := connectionhubv1alpha1.SubmarinerBroker{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.GetSubmarinerBrokerMetadata().Name,
			Labels: cr.GetLabels(),
		},
		Spec: brokerSpec,
	}

	return &broker
}

func GetSubmarinerOperator(cr *connectionhubv1alpha1.Submariner) *connectionhubv1alpha1.SubmarinerOperator {

	tenancy := cr.GetTenancySelectors()

	var clusterID connectionhubv1alpha1.InstanceType
	var token string
	var ca string

	if cr.Spec.InstanceType == connectionhubv1alpha1.InstanceTypeCloud {
		clusterID = connectionhubv1alpha1.InstanceType(tenancy.RobolaunchCloudInstance)
		token = cr.Status.BrokerStatus.Status.BrokerCredentials.Token
		ca = cr.Status.BrokerStatus.Status.BrokerCredentials.CA
	} else if cr.Spec.InstanceType == connectionhubv1alpha1.InstanceTypePhysical {
		clusterID = connectionhubv1alpha1.InstanceType(tenancy.RobolaunchPhysicalInstance)
		token = cr.Spec.BrokerCredentials.Token
		ca = cr.Spec.BrokerCredentials.CA
	}

	operatorSpec := connectionhubv1alpha1.SubmarinerOperatorSpec{
		ClusterCIDR:  cr.Spec.ClusterCIDR,
		ServiceCIDR:  cr.Spec.ServiceCIDR,
		PresharedKey: cr.Spec.PresharedKey,
		BrokerCredentials: connectionhubv1alpha1.BrokerCredentials{
			Token: token,
			CA:    ca,
		},
		ClusterID:    string(clusterID),
		APIServerURL: cr.Spec.APIServerURL,
		Helm:         cr.Spec.OperatorHelmChart,
	}

	operator := connectionhubv1alpha1.SubmarinerOperator{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.GetSubmarinerOperatorMetadata().Name,
			Labels: cr.GetLabels(),
		},
		Spec: operatorSpec,
	}

	return &operator
}

func GetSubmarinerCustomResource(cr *connectionhubv1alpha1.Submariner) *submv1alpha1.Submariner {

	submarinerOperator := GetSubmarinerOperator(cr)
	valuesObj := helm.GetSubmarinerOperatorValues(*submarinerOperator)

	submariner := submv1alpha1.Submariner{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetSubmarinerCustomResourceMetadata().Name,
			Namespace: cr.GetSubmarinerCustomResourceMetadata().Namespace,
		},
		Spec: submv1alpha1.SubmarinerSpec{
			Broker:                   "k8s",
			BrokerK8sApiServer:       valuesObj.Broker.Server,
			BrokerK8sApiServerToken:  valuesObj.Broker.Token,
			BrokerK8sCA:              valuesObj.Broker.Ca,
			BrokerK8sRemoteNamespace: connectionhubv1alpha1.SubmarinerBrokerNamespace,
			// BrokerK8sInsecure:        valuesObj.Broker.Insecure,
			// CeIPSecDebug:             valuesObj.IPSEC.Debug,
			// CeIPSecForceUDPEncaps:    valuesObj.IPSEC.ForceUDPEncaps,
			CeIPSecIKEPort:         valuesObj.IPSEC.IKEPort,
			CeIPSecNATTPort:        valuesObj.IPSEC.NATPort,
			CeNatDiscovery:         valuesObj.IPSEC.NATDiscovery,
			CeIPSecPreferredServer: valuesObj.IPSEC.PreferredServer,
			CeIPSecPSK:             valuesObj.IPSEC.PSK,
			ClusterCIDR:            valuesObj.Submariner.ClusterCIDR,
			ClusterID:              valuesObj.Submariner.ClusterID,
			// ColorCodes:               valuesObj.Submariner.ColorCodes,
			// Debug:                    valuesObj.Submariner.Debug,
			Namespace:   connectionhubv1alpha1.SubmarinerOperatorNamespace,
			NatEnabled:  valuesObj.Submariner.NatEnabled,
			Repository:  valuesObj.Submariner.Images.Repository,
			Version:     valuesObj.Submariner.Images.Tag,
			ServiceCIDR: valuesObj.Submariner.ServiceCIDR,
			// GlobalCIDR:               valuesObj.Submariner.GlobalCIDR,
			ServiceDiscoveryEnabled: valuesObj.Submariner.ServiceDiscovery,
			CableDriver:             valuesObj.Submariner.CableDriver,
			ConnectionHealthCheck: &submv1alpha1.HealthCheckSpec{
				Enabled:            valuesObj.Submariner.HealthCheckEnabled,
				IntervalSeconds:    1,
				MaxPacketLossCount: 5,
			},
		},
	}

	return &submariner

}
