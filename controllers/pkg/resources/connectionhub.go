package resources

import (
	"bytes"
	"encoding/base64"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
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

// for cloud instance
func GetFederationHost(cr *connectionhubv1alpha1.ConnectionHub) *connectionhubv1alpha1.FederationHost {

	federationHost := &connectionhubv1alpha1.FederationHost{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.GetFederationHostMetadata().Name,
		},
	}

	return federationHost
}

// for physical instances
func GetCloudInstance(cr *connectionhubv1alpha1.ConnectionHub) *connectionhubv1alpha1.CloudInstance {

	cloudInstance := &connectionhubv1alpha1.CloudInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.GetCloudInstanceMetadata().Name,
		},
	}

	return cloudInstance
}

// for cloud instances --dry-run
func GetConnectionHubTemplateKey(cr *connectionhubv1alpha1.ConnectionHub) (string, error) {

	cr.ObjectMeta.Labels["robolaunch.io/physical-instance"] = "<PHYSICAL-INSTANCE>"

	chTemplate := connectionhubv1alpha1.ConnectionHub{
		TypeMeta: metav1.TypeMeta{
			APIVersion: cr.APIVersion,
			Kind:       cr.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.Name,
			Labels: cr.ObjectMeta.Labels,
		},
		Spec: cr.Status.ConnectionInterfaces.ForPhysicalInstance,
	}

	var chTemplateYAMLBytes bytes.Buffer

	p := printers.YAMLPrinter{}
	p.PrintObj(&chTemplate, &chTemplateYAMLBytes)

	chTemplateEncoded := base64.StdEncoding.EncodeToString(chTemplateYAMLBytes.Bytes())

	return chTemplateEncoded, nil
}
