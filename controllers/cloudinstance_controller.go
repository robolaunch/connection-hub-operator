package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

// CloudInstanceReconciler reconciles a CloudInstance object
type CloudInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances/finalizers,verbs=update

func (r *CloudInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *CloudInstanceReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) error {
	return nil
}

func (r *CloudInstanceReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) error {
	return nil
}

func (r *CloudInstanceReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.CloudInstance, error) {
	instance := &connectionhubv1alpha1.CloudInstance{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.CloudInstance{}, err
	}

	return instance, nil
}

func (r *CloudInstanceReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.CloudInstance{}
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
func (r *CloudInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.CloudInstance{}).
		Complete(r)
}
