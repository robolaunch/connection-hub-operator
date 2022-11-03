package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	utils "github.com/robolaunch/connection-hub-operator/controllers/pkg/utils"
)

// FederationMemberReconciler reconciles a FederationMember object
type FederationMemberReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient dynamic.Interface
	RESTConfig    *rest.Config
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationmembers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationmembers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationmembers/finalizers,verbs=update

//+kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *FederationMemberReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

func (r *FederationMemberReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {
	switch instance.Status.Joined {
	case true:

	case false:

		host, err := r.reconcileGetOwner(ctx, instance)
		if err != nil {
			return err
		}

		err = utils.Join(host, instance, r.RESTConfig)
		if err != nil {
			return err
		}

		instance.Status.Joined = true
	}

	return nil
}

func (r *FederationMemberReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {

	// check joined resources

	return nil
}

func (r *FederationMemberReconciler) reconcileGetOwner(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) (*connectionhubv1alpha1.FederationHost, error) {
	host := &connectionhubv1alpha1.FederationHost{}
	err := r.Get(ctx, *instance.GetOwnerMetadata(), host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (r *FederationMemberReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.FederationMember, error) {
	instance := &connectionhubv1alpha1.FederationMember{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.FederationMember{}, err
	}

	return instance, nil
}

func (r *FederationMemberReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.FederationMember{}
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
func (r *FederationMemberReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.FederationMember{}).
		Complete(r)
}
