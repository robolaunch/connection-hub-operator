package resources

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSubmarinerBroker(cr *connectionhubv1alpha1.Submariner) *connectionhubv1alpha1.SubmarinerBroker {

	brokerSpec := connectionhubv1alpha1.SubmarinerBrokerSpec{
		HelmRepository: cr.Spec.HelmRepository,
		HelmChart:      cr.Spec.BrokerHelmChart,
		APIServerURL:   cr.Spec.APIServerURL,
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
		NetworkType:  cr.Spec.NetworkType,
		InstanceType: cr.Spec.InstanceType,
		ClusterCIDR:  cr.Spec.ClusterCIDR,
		ServiceCIDR:  cr.Spec.ServiceCIDR,
		PresharedKey: cr.Spec.PresharedKey,
		BrokerCredentials: connectionhubv1alpha1.BrokerCredentials{
			Token: token,
			CA:    ca,
		},
		ClusterID:      string(clusterID),
		APIServerURL:   cr.Spec.APIServerURL,
		HelmRepository: cr.Spec.HelmRepository,
		HelmChart:      cr.Spec.OperatorHelmChart,
		CableDriver:    cr.Spec.CableDriver,
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
