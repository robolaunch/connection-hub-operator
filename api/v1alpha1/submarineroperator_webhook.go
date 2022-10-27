package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var submarineroperatorlog = logf.Log.WithName("submarineroperator-resource")

func (r *SubmarinerOperator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-connection-hub-roboscale-io-v1alpha1-submarineroperator,mutating=true,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submarineroperators,verbs=create;update,versions=v1alpha1,name=msubmarineroperator.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &SubmarinerOperator{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SubmarinerOperator) Default() {
	submarineroperatorlog.Info("default", "name", r.Name)

	r.setInstanceType()
}

//+kubebuilder:webhook:path=/validate-connection-hub-roboscale-io-v1alpha1-submarineroperator,mutating=false,failurePolicy=fail,sideEffects=None,groups=connection-hub.roboscale.io,resources=submarineroperators,verbs=create;update,versions=v1alpha1,name=vsubmarineroperator.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &SubmarinerOperator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerOperator) ValidateCreate() error {
	submarineroperatorlog.Info("validate create", "name", r.Name)

	err := r.checkTenancyLabelsForSO()
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerOperator) ValidateUpdate(old runtime.Object) error {
	submarineroperatorlog.Info("validate update", "name", r.Name)

	err := r.checkTenancyLabelsForSO()
	if err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SubmarinerOperator) ValidateDelete() error {
	submarineroperatorlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *SubmarinerOperator) checkTenancyLabelsForSO() error {
	labels := r.GetLabels()

	if _, ok := labels[RobolaunchCloudInstanceLabelKey]; !ok {
		return errors.New("cloud instance label should be added with key " + RobolaunchCloudInstanceLabelKey)
	}

	return nil
}

func (r *SubmarinerOperator) setInstanceType() {

	tenancy := r.GetTenancySelectors()
	if tenancy.RobolaunchPhysicalInstance != "" {
		r.Spec.InstanceType = InstanceTypePhysical
	}

	r.Spec.InstanceType = InstanceTypeCloud

}
