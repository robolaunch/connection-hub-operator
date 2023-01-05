package controllers

import (
	"context"
	basicErr "errors"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubefed/pkg/apis/core/common"

	kubefedv1beta1 "github.com/robolaunch/connection-hub-operator/api/external/kubefed/v1beta1"
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

func (r *FederationMemberReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {
	switch instance.Status.Host.Exists {
	case true:

		switch instance.Status.JoinAttempted {
		case true:

			switch instance.Status.KubeFedClusterStatus.Created {
			case true:

				switch instance.Status.KubeFedClusterStatus.ConditionType {

				case common.ClusterReady:

					instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseReady

				case common.ClusterOffline:

					logger.Info("STATUS: Federation member is offline.")
					instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseOffline

				case common.ClusterConfigMalformed:

					logger.Info("STATUS: Federation member config is malfunctioned.")
					instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseMalfunctioned

				}

			case false:

				logger.Info("STATUS: Cluster " + instance.Name + " cannot join the federation.")
				instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseCannotJoinFederation

			}

		case false:

			switch instance.Spec.Server == "" {
			case true:

				instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseWaitingForCredentials

			case false:

				logger.Info("STATUS: Cluster " + instance.Name + " is joining the federation.")

				instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseJoiningFederation

				err := utils.JoinMember(instance, r.RESTConfig)
				if err != nil {
					return err
				}

				instance.Status.JoinAttempted = true

			}

		}

	case false:

		logger.Info("STATUS: Searching for federation host.")
		instance.Status.Phase = connectionhubv1alpha1.FederationMemberPhaseSearchingForHost

	}

	return nil
}

func (r *FederationMemberReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {

	err := r.reconcileCheckKubeFedClusterInstance(ctx, instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *FederationMemberReconciler) reconcileCheckKubeFedClusterInstance(ctx context.Context, instance *connectionhubv1alpha1.FederationMember) error {

	federationHostList := &connectionhubv1alpha1.FederationHostList{}
	err := r.List(ctx, federationHostList)
	if err != nil {
		return err
	} else {

		if len(federationHostList.Items) == 0 {
			instance.Status.Host = connectionhubv1alpha1.HostInstanceStatus{}
		} else if len(federationHostList.Items) > 1 {
			return basicErr.New("more than one federation hosts are listed")
		} else {
			federationHost := federationHostList.Items[0]
			instance.Status.Host.Exists = true
			instance.Status.Host.Name = federationHost.Name
		}
	}

	kubefedCluster := kubefedv1beta1.KubeFedCluster{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: connectionhubv1alpha1.FederationOperatorNamespace}, &kubefedCluster)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.KubeFedClusterStatus = connectionhubv1alpha1.KubeFedClusterStatus{}
	} else if err != nil {
		return err
	} else {

		instance.Status.KubeFedClusterStatus.Created = true

		if len(kubefedCluster.Status.Conditions) > 0 {
			instance.Status.KubeFedClusterStatus.ConditionType = kubefedCluster.Status.Conditions[0].Type
			instance.Status.KubeFedClusterStatus.Reason = *kubefedCluster.Status.Conditions[0].Reason
		} else {
			instance.Status.KubeFedClusterStatus.ConditionType = common.ClusterOffline
			instance.Status.KubeFedClusterStatus.Reason = "UnknownCondition"
		}
	}

	return nil
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
		Watches(
			&source.Kind{Type: &kubefedv1beta1.KubeFedCluster{}},
			handler.EnqueueRequestsFromMapFunc(r.watchKubeFedClusters)).
		Watches(
			&source.Kind{Type: &connectionhubv1alpha1.FederationHost{}},
			handler.EnqueueRequestsFromMapFunc(r.watchFederationHost),
			builder.WithPredicates(predicate.Funcs{
				CreateFunc: func(ce event.CreateEvent) bool {
					return true
				},
				DeleteFunc: func(de event.DeleteEvent) bool {
					return true
				},
			})).
		Complete(r)
}

func (r *FederationMemberReconciler) watchKubeFedClusters(o client.Object) []reconcile.Request {

	cluster := o.(*kubefedv1beta1.KubeFedCluster)

	federationMember := &connectionhubv1alpha1.FederationMember{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: cluster.Name,
	}, federationMember)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name: federationMember.Name,
			},
		},
	}
}

func (r *FederationMemberReconciler) watchFederationHost(o client.Object) []reconcile.Request {

	federationMembers := &connectionhubv1alpha1.FederationMemberList{}
	err := r.List(context.TODO(), federationMembers)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(federationMembers.Items))
	for i, item := range federationMembers.Items {

		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: item.Name,
			},
		}

	}

	return requests
}
