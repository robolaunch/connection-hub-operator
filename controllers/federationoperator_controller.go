package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
)

// FederationOperatorReconciler reconciles a FederationOperator object
type FederationOperatorReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient dynamic.Interface
	RESTConfig    *rest.Config
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationoperators/finalizers,verbs=update

//+kubebuilder:rbac:groups=multiclusterdns.kubefed.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=scheduling.kubefed.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=core.kubefed.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=types.kubefed.io,resources=*,verbs=*

func (r *FederationOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.reconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.reconcileCheckStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileCheckResources(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FederationOperatorReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	switch instance.Status.NamespaceStatus.Created {
	case true:

		switch instance.Status.ChartStatus.Deployed {
		case true:

			switch instance.Status.ChartResourceStatus.Deployed {
			case true:

				instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseDeployed

			case false:

				logger.Info("STATUS: Checking for Federation Operator resources.")
				instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseCheckingResources

			}

		case false:

			err := r.reconcileInstallChart(ctx, instance)
			if err != nil {
				return err
			}

		}

	case false:

		err := r.reconcileCreateNamespace(ctx, instance)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *FederationOperatorReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {
	return nil
}

func (r *FederationOperatorReconciler) reconcileCreateNamespace(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseCreatingNamespace

	operatorNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.GetNamespaceMetadata().Name,
		},
	}

	err := ctrl.SetControllerReference(instance, operatorNamespace, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, operatorNamespace)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Federation Operator's namespace is created.")
	instance.Status.NamespaceStatus.Created = true

	return nil
}

func (r *FederationOperatorReconciler) reconcileInstallChart(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseDeployingChart

	err := helmops.InstallFederationOperatorChart(*instance, r.RESTConfig)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Federation Operator Helm chart is deployed.")
	instance.Status.ChartStatus.Deployed = true

	return nil
}

func (r *FederationOperatorReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.FederationOperator, error) {
	instance := &connectionhubv1alpha1.FederationOperator{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.FederationOperator{}, err
	}

	return instance, nil
}

func (r *FederationOperatorReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.FederationOperator{}
		err := r.Get(ctx, types.NamespacedName{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		}, instanceLV)

		if err == nil {
			instance.ResourceVersion = instanceLV.ResourceVersion
		}

		err1 := r.Status().Update(ctx, instance)
		return err1
	})
}

// SetupWithManager sets up the controller with the Manager.
func (r *FederationOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.FederationOperator{}).
		Complete(r)
}
