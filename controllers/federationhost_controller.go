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
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"github.com/robolaunch/connection-hub-operator/controllers/pkg/resources"
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

		err := r.reconcileUpdateMemberObjects(ctx, instance)
		if err != nil {
			return err
		}

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

	err := r.reconcileUpdateMemberStatuses(ctx, instance)
	if err != nil {
		return err
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

func (r *FederationHostReconciler) reconcileUpdateMemberStatuses(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {

	if instance.Status.MemberStatuses == nil {
		instance.Status.MemberStatuses = make(map[string]connectionhubv1alpha1.MemberStatus)
	}

	// update member statuses
	for name := range instance.Spec.FederationMembers {
		if status, ok := instance.Status.MemberStatuses[name]; ok {
			// keep member active
			status.ResourcePhase = connectionhubv1alpha1.MemberResourcePhaseActive
			instance.Status.MemberStatuses[name] = status
		} else {
			// add member to status
			instance.Status.MemberStatuses[name] = connectionhubv1alpha1.MemberStatus{
				Created:       false,
				ResourcePhase: connectionhubv1alpha1.MemberResourcePhaseActive,
			}
		}

		// check member object
		federationMember := connectionhubv1alpha1.FederationMember{}
		err := r.Get(ctx, types.NamespacedName{Name: name}, &federationMember)
		if err != nil && errors.IsNotFound(err) {
			status := instance.Status.MemberStatuses[name]
			status.Created = false
			status.Status = connectionhubv1alpha1.FederationMemberStatus{}
			instance.Status.MemberStatuses[name] = status
		} else if err != nil {
			return err
		} else {
			status := instance.Status.MemberStatuses[name]
			status.Created = true
			status.Status = federationMember.Status
			instance.Status.MemberStatuses[name] = status
		}
	}

	// make member idle
	for name, status := range instance.Status.MemberStatuses {
		if _, ok := instance.Spec.FederationMembers[name]; !ok {
			status.ResourcePhase = connectionhubv1alpha1.MemberResourcePhaseIdle
			instance.Status.MemberStatuses[name] = status
		}
	}

	// delete idles and deleted ones from status
	for name, status := range instance.Status.MemberStatuses {
		if !status.Created && status.ResourcePhase == connectionhubv1alpha1.MemberResourcePhaseIdle {
			delete(instance.Status.MemberStatuses, name)
		}
	}

	return nil
}

func (r *FederationHostReconciler) reconcileUpdateMemberObjects(ctx context.Context, instance *connectionhubv1alpha1.FederationHost) error {

	// check objects

	// create members in spec
	for name, member := range instance.Spec.FederationMembers {
		if status, ok := instance.Status.MemberStatuses[name]; ok && !status.Created && status.ResourcePhase == connectionhubv1alpha1.MemberResourcePhaseActive {
			federationMember := resources.GetFederationMember(name, member)
			err := ctrl.SetControllerReference(instance, federationMember, r.Scheme)
			if err != nil {
				return err
			}

			err = r.Create(ctx, federationMember)
			if err != nil {
				return err
			}

			status.Created = true
			instance.Status.MemberStatuses[name] = status
		}
	}

	// delete idle member objects
	for name, status := range instance.Status.MemberStatuses {
		if status.Created && status.ResourcePhase == connectionhubv1alpha1.MemberResourcePhaseIdle {
			idleFederationMember := &connectionhubv1alpha1.FederationMember{}
			err := r.Get(ctx, types.NamespacedName{Name: name}, idleFederationMember)
			if err != nil {
				return err
			}

			err = r.Delete(ctx, idleFederationMember)
			if err != nil {
				return err
			}

			delete(instance.Status.MemberStatuses, name)
		}
	}

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
		Owns(&connectionhubv1alpha1.FederationMember{}).
		Complete(r)
}
