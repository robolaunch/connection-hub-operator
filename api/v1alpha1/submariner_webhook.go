package v1alpha1

import (
	"errors"

	"github.com/lucasjones/reggen"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var submarinerlog = logf.Log.WithName("submariner-resource")

func (r *Submariner) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-submariner,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submariners,verbs=create;update,versions=v1alpha1,name=msubmariner.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Submariner{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Submariner) Default() {
	submarinerlog.Info("default", "name", r.Name)

	r.SetName(GlobalSubmarinerResourceName)
	r.setInstanceType()

	if instanceType := r.getInstanceType(); instanceType == InstanceTypeCloud {
		r.generatePresharedKey()
	}
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-submariner,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submariners,verbs=create;update,versions=v1alpha1,name=vsubmariner.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Submariner{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateCreate() error {
	submarinerlog.Info("validate create", "name", r.Name)

	err := r.checkTenancyLabelsForSubmariner()
	if err != nil {
		return err
	}

	instanceType := r.getInstanceType()

	switch instanceType {
	case InstanceTypeCloud:

		err := r.checkBrokerCredentialsForCloudInstance()
		if err != nil {
			return err
		}

	case InstanceTypePhysical:

		err := r.checkBrokerCredentialsForPhysicalInstance()
		if err != nil {
			return err
		}

		err = r.checkPresharedKeyForPhysicalInstance()
		if err != nil {
			return err
		}

	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateUpdate(old runtime.Object) error {
	submarinerlog.Info("validate update", "name", r.Name)

	err := r.checkTenancyLabelsForSubmariner()
	if err != nil {
		return err
	}

	instanceType := r.getInstanceType()

	switch instanceType {
	case InstanceTypeCloud:

		err := r.checkBrokerCredentialsForCloudInstance()
		if err != nil {
			return err
		}

	case InstanceTypePhysical:

		err := r.checkBrokerCredentialsForPhysicalInstance()
		if err != nil {
			return err
		}

		err = r.checkPresharedKeyForPhysicalInstance()
		if err != nil {
			return err
		}

	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateDelete() error {
	submarinerlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *Submariner) checkTenancyLabelsForSubmariner() error {
	labels := r.GetLabels()

	if _, ok := labels[RobolaunchCloudInstanceLabelKey]; !ok {
		return errors.New("cloud instance label should be added with key " + RobolaunchCloudInstanceLabelKey)
	}

	return nil
}

func (r *Submariner) setInstanceType() {
	r.Spec.InstanceType = r.getInstanceType()
}

func (r *Submariner) getInstanceType() InstanceType {
	tenancy := r.GetTenancySelectors()
	if tenancy.RobolaunchPhysicalInstance != "" {
		return InstanceTypePhysical
	}

	return InstanceTypeCloud
}

func (r *Submariner) generatePresharedKey() {

	psk, err := reggen.Generate("^([A-Za-z0-9]){64}$", 0)
	if err != nil {
		r.Spec.PresharedKey = "cfZ7CsGkN88z6eW3Z0A6Pj5W0G46GJuNKQu6onvscD19FIbOYpqe9OrNmL1R303Q"
	}

	r.Spec.PresharedKey = psk

}

func (r *Submariner) checkPresharedKeyForPhysicalInstance() error {
	if r.Spec.PresharedKey == "" {
		return errors.New("field `spec.presharedKey` cannot be empty in physical instances")
	}

	return nil
}

func (r *Submariner) checkBrokerCredentialsForPhysicalInstance() error {
	if r.Spec.BrokerCredentials.Token == "" {
		return errors.New("field `spec.brokerCredentials.token` cannot be empty in physical instances")
	}

	if r.Spec.BrokerCredentials.CA == "" {
		return errors.New("field `spec.brokerCredentials.ca` cannot be empty in physical instances")
	}

	return nil
}

func (r *Submariner) checkBrokerCredentialsForCloudInstance() error {
	if r.Spec.BrokerCredentials.Token != "" {
		return errors.New("field `spec.brokerCredentials.token` should be empty in physical instances")
	}

	if r.Spec.BrokerCredentials.CA != "" {
		return errors.New("field `spec.brokerCredentials.ca` should be empty in physical instances")
	}

	return nil
}
