package controllers

import (
	"context"
	basicErr "errors"

	corev1 "k8s.io/api/core/v1"
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

	brokerv1 "github.com/robolaunch/connection-hub-operator/api/external/submariner/v1"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	robotErr "github.com/robolaunch/connection-hub-operator/controllers/pkg/error"
	mcsv1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

// CloudInstanceReconciler reconciles a CloudInstance object
type CloudInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=cloudinstances/finalizers,verbs=update

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submariners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete

func (r *CloudInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	instance, err := r.reconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if instance.Status.BootID == "" {
		activeNode, err := r.reconcileCheckNode(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		instance.Status.BootID = activeNode.Status.NodeInfo.BootID
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

func (r *CloudInstanceReconciler) reconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) error {

	switch instance.Status.DeployerStatus.Exists {
	case true:

		switch instance.Status.DeployerStatus.Phase {
		case connectionhubv1alpha1.SubmarinerPhaseReadyToConnect:

			switch instance.Status.ConnectionResources.ClusterStatus.Exists && instance.Status.ConnectionResources.EndpointStatus.Exists {
			case true:

				switch instance.Status.GatewayConnection.ConnectionStatus {
				case brokerv1.Connected:

					activeNode, err := r.reconcileCheckNode(ctx, instance)
					if err != nil {
						return err
					}

					if instance.Status.Phase != connectionhubv1alpha1.CloudInstancePhaseConnected || instance.Status.BootID != activeNode.Status.NodeInfo.BootID {
						logger.Info("INFO: Deleting all of the ServiceExport objects.")
						serviceExportList := mcsv1alpha1.ServiceExportList{}
						err := r.List(ctx, &serviceExportList)
						if err != nil {
							return err
						}

						for _, v := range serviceExportList.Items {
							err := r.Delete(ctx, &v)
							if err != nil {
								return err
							}
						}

						instance.Status.BootID = activeNode.Status.NodeInfo.BootID
					}

					instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseConnected

				case brokerv1.Connecting:

					instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseConnecting

				case brokerv1.ConnectionError:

					instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseNotConnected

				}

			case false:

				instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseWaitingForResources

			}

		default:

			instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseWaitingForDeployer

		}

	case false:
		instance.Status.Phase = connectionhubv1alpha1.CloudInstancePhaseLookingForDeployer
	}

	return nil
}

func (r *CloudInstanceReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) error {

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

	req, err := labels.NewRequirement(connectionhubv1alpha1.EndpointClusterIDLabelKey, selection.In, []string{instance.Name})
	if err != nil {
		return err
	}

	labelSelector := labels.NewSelector().Add(*req)

	submarinerEndpoints := &brokerv1.EndpointList{}
	err = r.List(ctx, submarinerEndpoints, &client.ListOptions{
		LabelSelector: labelSelector,
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

	// check gateways.submariner.io

	submarinerGateways := &brokerv1.GatewayList{}
	err = r.List(ctx, submarinerGateways, &client.ListOptions{
		Namespace: connectionhubv1alpha1.SubmarinerOperatorNamespace,
	})
	if err != nil {
		return err
	} else {

		if len(submarinerGateways.Items) < 1 {
			instance.Status.GatewayConnection = connectionhubv1alpha1.GatewayConnection{}
		} else if len(submarinerGateways.Items) == 1 {
			gw := submarinerGateways.Items[0]
			instance.Status.GatewayConnection.GatewayResource = gw.Name
			for _, connection := range gw.Status.Connections {
				if instance.Name == connection.Endpoint.ClusterID {
					instance.Status.GatewayConnection.ClusterID = connection.Endpoint.ClusterID
					instance.Status.GatewayConnection.Hostname = connection.Endpoint.Hostname
					instance.Status.GatewayConnection.ConnectionStatus = connection.Status
				}
			}
		} else {
			return basicErr.New("more than one gateways is listed")
		}

	}

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

func (r *CloudInstanceReconciler) reconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.CloudInstance) (*corev1.Node, error) {

	tenancyMap := connectionhubv1alpha1.GetTenancyMap(instance)

	requirements := []labels.Requirement{}
	for k, v := range tenancyMap {
		newReq, err := labels.NewRequirement(k, selection.In, []string{v})
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, *newReq)
	}

	nodeSelector := labels.NewSelector().Add(requirements...)

	nodes := &corev1.NodeList{}
	err := r.List(ctx, nodes, &client.ListOptions{
		LabelSelector: nodeSelector,
	})
	if err != nil {
		return nil, err
	}

	if len(nodes.Items) == 0 {
		return nil, &robotErr.NodeNotFoundError{
			ResourceKind:      instance.Kind,
			ResourceName:      instance.Name,
			ResourceNamespace: instance.Namespace,
		}
	} else if len(nodes.Items) > 1 {
		return nil, &robotErr.MultipleNodeFoundError{
			ResourceKind:      instance.Kind,
			ResourceName:      instance.Name,
			ResourceNamespace: instance.Namespace,
		}
	}

	return &nodes.Items[0], nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CloudInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.CloudInstance{}).
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
		Watches(
			&source.Kind{Type: &brokerv1.Gateway{}},
			handler.EnqueueRequestsFromMapFunc(r.watchGateways)).
		Complete(r)
}

func (r *CloudInstanceReconciler) watchSubmarinerDeployer(o client.Object) []reconcile.Request {
	cloudInstances := &connectionhubv1alpha1.CloudInstanceList{}
	err := r.List(context.TODO(), cloudInstances)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(cloudInstances.Items))
	for i, item := range cloudInstances.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name: item.Name,
			},
		}
	}

	return requests
}

func (r *CloudInstanceReconciler) watchClusters(o client.Object) []reconcile.Request {

	cluster := o.(*brokerv1.Cluster)

	cloudInstance := &connectionhubv1alpha1.CloudInstance{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: cluster.Name,
	}, cloudInstance)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name: cloudInstance.Name,
			},
		},
	}
}

func (r *CloudInstanceReconciler) watchEndpoints(o client.Object) []reconcile.Request {

	endpoint := o.(*brokerv1.Endpoint)

	cloudInstance := &connectionhubv1alpha1.CloudInstance{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: endpoint.Labels[connectionhubv1alpha1.EndpointClusterIDLabelKey],
	}, cloudInstance)
	if err != nil {
		return []reconcile.Request{}
	}

	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name: cloudInstance.Name,
			},
		},
	}

}

func (r *CloudInstanceReconciler) watchGateways(o client.Object) []reconcile.Request {

	gateway := o.(*brokerv1.Gateway)

	cloudInstances := &connectionhubv1alpha1.CloudInstanceList{}
	err := r.List(context.TODO(), cloudInstances)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(cloudInstances.Items))
	for i, item := range cloudInstances.Items {

		for _, conn := range gateway.Status.Connections {

			if conn.Endpoint.ClusterID == item.Name {

				requests[i] = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name: item.Name,
					},
				}

			}

		}

	}

	return requests
}
