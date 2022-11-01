package controllers

import (
	"context"
	basicErr "errors"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	brokerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
)

// PhysicalInstanceReconciler reconciles a PhysicalInstance object
type PhysicalInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=physicalinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=physicalinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=physicalinstances/finalizers,verbs=update

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submariners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete

func (r *PhysicalInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

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

func (r *PhysicalInstanceReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.PhysicalInstance) error {

	switch instance.Status.DeployerStatus.Exists {
	case true:

		switch instance.Status.DeployerStatus.Phase {
		case connectionhubv1alpha1.SubmarinerPhaseReadyToConnect:

			switch instance.Status.ConnectionResources.ClusterStatus.Exists && instance.Status.ConnectionResources.EndpointStatus.Exists {
			case true:

				instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseConnected

			case false:

				instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseRegisteredAndTryingToConnect

			}

		default:

			instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseWaitingForDeployer

		}

	case false:
		instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseLookingForDeployer
	}

	return nil
}

func (r *PhysicalInstanceReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.PhysicalInstance) error {

	// check submariners.connection-hub.roboscale.io

	instance.Status.DeployerStatus.Name = instance.GetSubmarinerDeployerMetadata().Name

	submarinerDeployer := &connectionhubv1alpha1.Submariner{}
	err := r.Get(ctx, instance.GetSubmarinerDeployerMetadata(), submarinerDeployer)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.DeployerStatus = connectionhubv1alpha1.DeployerStatus{}
	} else if err != nil {
		return err
	} else {

		instance.Status.DeployerStatus.Exists = true
		instance.Status.DeployerStatus.Phase = submarinerDeployer.Status.Phase

	}

	// check clusters.submariner.io

	instance.Status.ConnectionResources.ClusterStatus.Name = instance.GetSubmarinerClusterMetadata().Name

	submarinerCluster := &brokerv1.Cluster{}
	err = r.Get(ctx, instance.GetSubmarinerClusterMetadata(), submarinerCluster)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.ConnectionResources.ClusterStatus = connectionhubv1alpha1.ConnectionResourceStatus{}
	} else if err != nil {
		return err
	} else {

		instance.Status.ConnectionResources.ClusterStatus.Exists = true

	}

	// check endpoints.submariner.io

	req, err := labels.NewRequirement(connectionhubv1alpha1.EndpointClusterIDLabelKey, selection.Equals, []string{instance.Name})
	if err != nil {
		return err
	}

	labelSelector := labels.NewSelector().Add(*req)

	submarinerEndpoints := &brokerv1.EndpointList{}
	err = r.List(ctx, submarinerEndpoints, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     connectionhubv1alpha1.SubmarinerOperatorNamespace,
	})
	if err != nil {
		return err
	} else {

		if len(submarinerEndpoints.Items) < 1 {
			instance.Status.ConnectionResources.EndpointStatus = connectionhubv1alpha1.ConnectionResourceStatus{}
		} else if len(submarinerEndpoints.Items) == 1 {
			instance.Status.ConnectionResources.EndpointStatus.Name = submarinerEndpoints.Items[0].Name
			instance.Status.ConnectionResources.EndpointStatus.Exists = true
		} else {
			return basicErr.New("more than one endpoints is listed with same clusterID")
		}

	}

	return nil
}

func (r *PhysicalInstanceReconciler) reconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.PhysicalInstance, error) {
	instance := &connectionhubv1alpha1.PhysicalInstance{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.PhysicalInstance{}, err
	}

	return instance, nil
}

func (r *PhysicalInstanceReconciler) reconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.PhysicalInstance) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.PhysicalInstance{}
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
func (r *PhysicalInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.PhysicalInstance{}).
		Watches(
			&source.Kind{Type: &connectionhubv1alpha1.Submariner{}},
			handler.EnqueueRequestsFromMapFunc(r.watchSubmarinerDeployer),
		).
		Watches(
			&source.Kind{Type: &brokerv1.Cluster{}},
			handler.EnqueueRequestsFromMapFunc(r.watchClusters)).
		Watches(
			&source.Kind{Type: &brokerv1.Endpoint{}},
			handler.EnqueueRequestsFromMapFunc(r.watchEndpoints)).
		Complete(r)
}

func (r *PhysicalInstanceReconciler) watchSubmarinerDeployer(o client.Object) []reconcile.Request {
	physicalInstances := &connectionhubv1alpha1.PhysicalInstanceList{}
	err := r.List(context.TODO(), physicalInstances)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(physicalInstances.Items))
	for i, item := range physicalInstances.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: item.Name,
			},
		}
	}

	return requests
}

func (r *PhysicalInstanceReconciler) watchClusters(o client.Object) []reconcile.Request {

	cluster := o.(*brokerv1.Cluster)

	physicalInstance := &connectionhubv1alpha1.PhysicalInstance{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: cluster.Name,
	}, physicalInstance)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name: physicalInstance.Name,
			},
		},
	}
}

func (r *PhysicalInstanceReconciler) watchEndpoints(o client.Object) []reconcile.Request {

	endpoint := o.(*brokerv1.Endpoint)

	physicalInstance := &connectionhubv1alpha1.PhysicalInstance{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: endpoint.Labels[connectionhubv1alpha1.EndpointClusterIDLabelKey],
	}, physicalInstance)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name: physicalInstance.Name,
			},
		},
	}

}
