package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var federationoperatorlog = logf.Log.WithName("federationoperator-resource")

func (r *FederationOperator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-federationoperator,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=federationoperators,verbs=create;update,versions=v1alpha1,name=mfederationoperator.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &FederationOperator{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FederationOperator) Default() {
	federationoperatorlog.Info("default", "name", r.Name)

	r.SetName(GlobalFederationOperatorResourceName)
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-federationoperator,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=federationoperators,verbs=create;update,versions=v1alpha1,name=vfederationoperator.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &FederationOperator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *FederationOperator) ValidateCreate() error {
	federationoperatorlog.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *FederationOperator) ValidateUpdate(old runtime.Object) error {
	federationoperatorlog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *FederationOperator) ValidateDelete() error {
	federationoperatorlog.Info("validate delete", "name", r.Name)
	return nil
}
