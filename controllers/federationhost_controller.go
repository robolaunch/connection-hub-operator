package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
)

// FederationHostReconciler reconciles a FederationHost object
type FederationHostReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationhosts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationhosts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationhosts/finalizers,verbs=update

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=federationmembers,verbs=get;list;watch;create;update;patch;delete

func (r *FederationHostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.reconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.reconcileCheckDeletion(ctx, instance)
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

func (r *FederationHostReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {

	switch instance.Status.SelfJoined {
	case true:

		instance.Status.Phase = connectionhubv1alpha1.FederationHostPhaseReady

	case false:

		instance.Status.Phase = connectionhubv1alpha1.FederationHostPhaseJoiningSelf

		err := r.reconcileCreateHostMember(ctx, instance)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *FederationHostReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {

	instance.Status.Members = make(map[string]connectionhubv1alpha1.FederationMemberStatus)

	federationMemberList := &connectionhubv1alpha1.FederationMemberList{}
	err := r.List(ctx, federationMemberList)
	if err != nil {
		return err
	} else {
		for _, member := range federationMemberList.Items {
			instance.Status.Members[member.Name] = member.Status
		}
	}

	return nil
}

func (r *FederationHostReconciler) reconcileCreateHostMember(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {

	member := &connectionhubv1alpha1.FederationMember{
		ObjectMeta: v1.ObjectMeta{
			Name: instance.Name,
		},
		Spec: connectionhubv1alpha1.FederationMemberSpec{
			Server: "",
			Credentials: connectionhubv1alpha1.PhysicalInstanceCredentials{
				CertificateAuthority: "",
				ClientKey:            "",
				ClientCertificate:    "",
			},
		},
	}

	err := ctrl.SetControllerReference(instance, member, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, member)
	if err != nil {
		return err
	}

	instance.Status.SelfJoined = true

	return nil
}

func (r *FederationHostReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.FederationHost, error) {
	instance := &connectionhubv1alpha1.FederationHost{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.FederationHost{}, err
	}

	return instance, nil
}

func (r *FederationHostReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.FederationHost{}
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
func (r *FederationHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.FederationHost{}).
		Watches(
			&source.Kind{Type: &connectionhubv1alpha1.FederationMember{}},
			handler.EnqueueRequestsFromMapFunc(r.watchMembers)).
		Complete(r)
}

func (r *FederationHostReconciler) watchMembers(o client.Object) []reconcile.Request {

	federationHosts := &connectionhubv1alpha1.FederationHostList{}
	err := r.List(context.TODO(), federationHosts)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(federationHosts.Items))
	for i, item := range federationHosts.Items {

		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: item.Name,
			},
		}

	}

	return requests
}
