package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var connectionhublog = logf.Log.WithName("connectionhub-resource")

func (r *ConnectionHub) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-connectionhub,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=connectionhubs,verbs=create;update,versions=v1alpha1,name=mconnectionhub.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ConnectionHub{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ConnectionHub) Default() {
	connectionhublog.Info("default", "name", r.Name)

	r.setInstanceType()
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-connectionhub,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=connectionhubs,verbs=create;update,versions=v1alpha1,name=vconnectionhub.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ConnectionHub{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ConnectionHub) ValidateCreate() error {
	connectionhublog.Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ConnectionHub) ValidateUpdate(old runtime.Object) error {
	connectionhublog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ConnectionHub) ValidateDelete() error {
	connectionhublog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *ConnectionHub) setInstanceType() {
	r.Spec.InstanceType = r.getInstanceType()
}

func (r *ConnectionHub) getInstanceType() InstanceType {
	tenancy := r.GetTenancySelectors()
	if tenancy.RobolaunchPhysicalInstance != "" {
		return InstanceTypePhysical
	}

	return InstanceTypeCloud
}
