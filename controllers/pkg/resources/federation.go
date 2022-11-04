package resources

import (
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetFederationMember(member connectionhubv1alpha1.MemberInfo) *connectionhubv1alpha1.FederationMember {

	memberObj := connectionhubv1alpha1.FederationMember{
		ObjectMeta: metav1.ObjectMeta{
			Name: member.Name,
		},
		Spec: member.MemberSpec,
	}

	return &memberObj
}
