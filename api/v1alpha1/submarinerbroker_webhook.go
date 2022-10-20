package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var submarinerbrokerlog = logf.Log.WithName("submarinerbroker-resource")

func (r *SubmarinerBroker) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-submarinerbroker,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submarinerbrokers,verbs=create;update,versions=v1alpha1,name=msubmarinerbroker.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &SubmarinerBroker{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SubmarinerBroker) Default() {
	submarinerbrokerlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-submarinerbroker,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submarinerbrokers,verbs=create;update,versions=v1alpha1,name=vsubmarinerbroker.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &SubmarinerBroker{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerBroker) ValidateCreate() error {
	submarinerbrokerlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerBroker) ValidateUpdate(old runtime.Object) error {
	submarinerbrokerlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerBroker) ValidateDelete() error {
	submarinerbrokerlog.Info("validate delete", "name", r.Name)

	return nil
}
