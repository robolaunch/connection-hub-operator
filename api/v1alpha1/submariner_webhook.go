package v1alpha1

import (
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
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-submariner,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submariners,verbs=create;update,versions=v1alpha1,name=vsubmariner.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Submariner{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateCreate() error {
	submarinerlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateUpdate(old runtime.Object) error {
	submarinerlog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Submariner) ValidateDelete() error {
	submarinerlog.Info("validate delete", "name", r.Name)

	return nil
}
