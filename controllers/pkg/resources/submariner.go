package resources

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSubmarinerBroker(cr *connectionhubv1alpha1.Submariner) *connectionhubv1alpha1.SubmarinerBroker {

	brokerSpec := connectionhubv1alpha1.SubmarinerBrokerSpec{
		Helm:      cr.Spec.BrokerHelmChart,
		BrokerURL: cr.Spec.APIServerURL,
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

	token := cr.Status.BrokerStatus.Status.Broker.BrokerToken
	ca := cr.Status.BrokerStatus.Status.Broker.BrokerCA

	// TODO: generate some of the fields in cr

	operatorSpec := connectionhubv1alpha1.SubmarinerOperatorSpec{
		ClusterCIDR:  cr.Spec.ClusterCIDR,
		ServiceCIDR:  cr.Spec.ServiceCIDR,
		PresharedKey: cr.Spec.PresharedKey,
		Broker: connectionhubv1alpha1.BrokerInfo{
			BrokerURL:   cr.Spec.APIServerURL,
			BrokerToken: token,
			BrokerCA:    ca,
		},
		ClusterID: cr.Spec.ClusterID,
		Helm:      cr.Spec.OperatorHelmChart,
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
