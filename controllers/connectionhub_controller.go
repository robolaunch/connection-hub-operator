package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

// ConnectionHubReconciler reconciles a ConnectionHub object
type ConnectionHubReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=connectionhubs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=connectionhubs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=connectionhubs/finalizers,verbs=update

func (r *ConnectionHubReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

func (r *ConnectionHubReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {

	switch instance.Status.Submariner.Created {
	case true:

		switch instance.Status.Submariner.Phase {
		case connectionhubv1alpha1.SubmarinerPhaseReadyToConnect:

			switch instance.Status.Federation.Created {
			case true:

				switch instance.Status.Federation.Phase {
				case connectionhubv1alpha1.FederationOperatorPhaseDeployed:

				default:

					// wait for federation to be ready

				}

			case false:

				// create federation

			}

		default:

			// wait for submariner to be ready

		}

	case false:

		// create submariner

	}

	return nil
}

func (r *ConnectionHubReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {

	// check submariner

	// check federation

	return nil
}

func (r *ConnectionHubReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.ConnectionHub, error) {
	instance := &connectionhubv1alpha1.ConnectionHub{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.ConnectionHub{}, err
	}

	return instance, nil
}

func (r *ConnectionHubReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.ConnectionHub{}
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
func (r *ConnectionHubReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.ConnectionHub{}).
		Complete(r)
}
