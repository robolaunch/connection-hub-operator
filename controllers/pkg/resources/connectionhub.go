package resources

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSubmariner(cr *connectionhubv1alpha1.ConnectionHub) *connectionhubv1alpha1.Submariner {

	labels := cr.GetLabels()

	submariner := &connectionhubv1alpha1.Submariner{
		ObjectMeta: metav1.ObjectMeta{
			Name:   connectionhubv1alpha1.GlobalSubmarinerResourceName,
			Labels: labels,
		},
		Spec: cr.Spec.SubmarinerSpec,
	}

	return submariner
}

func GetFederation(cr *connectionhubv1alpha1.ConnectionHub) *connectionhubv1alpha1.FederationOperator {

	labels := cr.GetLabels()

	federation := &connectionhubv1alpha1.FederationOperator{
		ObjectMeta: metav1.ObjectMeta{
			Name:   connectionhubv1alpha1.GlobalFederationOperatorResourceName,
			Labels: labels,
		},
		Spec: cr.Spec.FederationSpec,
	}

	return federation
}
