package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var cloudinstancelog = logf.Log.WithName("cloudinstance-resource")

func (r *CloudInstance) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-cloudinstance,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=cloudinstances,verbs=create;update,versions=v1alpha1,name=mcloudinstance.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CloudInstance{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CloudInstance) Default() {
	cloudinstancelog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-cloudinstance,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=cloudinstances,verbs=create;update,versions=v1alpha1,name=vcloudinstance.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CloudInstance{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CloudInstance) ValidateCreate() error {
	cloudinstancelog.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CloudInstance) ValidateUpdate(old runtime.Object) error {
	cloudinstancelog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CloudInstance) ValidateDelete() error {
	cloudinstancelog.Info("validate delete", "name", r.Name)
	return nil
}
