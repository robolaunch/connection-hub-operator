package controllers

import (
	"context"
	basicErr "errors"
	"reflect"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"github.com/robolaunch/connection-hub-operator/controllers/pkg/resources"
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

	switch instance.Status.Submariner.DeployerStatus.Exists {
	case true:

		switch instance.Status.Submariner.DeployerStatus.Phase {
		case connectionhubv1alpha1.SubmarinerPhaseReadyToConnect:

			switch instance.Status.Submariner.ConnectionResources.ClusterStatus.Exists && instance.Status.Submariner.ConnectionResources.EndpointStatus.Exists {
			case true:

				switch instance.Status.Submariner.GatewayConnection.ConnectionStatus {
				case brokerv1.Connected:

					instance.Status.MulticastConnectionPhase = connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseConnected

				case brokerv1.Connecting:

					instance.Status.MulticastConnectionPhase = connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseConnecting

				case brokerv1.ConnectionError:

					instance.Status.MulticastConnectionPhase = connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseFailed

				}

			case false:

				instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseRegistered
				instance.Status.MulticastConnectionPhase = connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseWaitingForConnection

			}

		default:

			instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseWaitingForDeployer

		}

	case false:
		instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseLookingForDeployer
	}

	switch instance.Status.MulticastConnectionPhase {
	case connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseConnected:

		switch instance.Status.FederationMember.Created {
		case true:

			switch instance.Status.FederationMember.Status.Phase {
			case connectionhubv1alpha1.FederationMemberPhaseReady:

				instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseConnected

			case connectionhubv1alpha1.FederationMemberPhaseWaitingForCredentials:

				instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseWaitingForCredentials

			default:

				instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseConnecting

			}

		case false:

			instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseConnecting

			err := r.reconcileCreateFederationMember(ctx, instance)
			if err != nil {
				return err
			}

		}

	default:

		instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseWaitingForMulticast

	}

	if instance.Status.MulticastConnectionPhase == connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseConnected &&
		instance.Status.FederationConnectionPhase == connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseConnected {

		switch instance.Status.RelayServerPodStatus.Created {
		case true:

			switch instance.Status.RelayServerPodStatus.Phase {
			case corev1.PodRunning:
				switch instance.Status.RelayServerServiceStatus.Created {
				case true:
					instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseConnected
				case false:
					instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseCreatingRelayServer
					svc := resources.GetRelayServerService(instance)
					err := r.Create(ctx, svc)
					if err != nil {
						return err
					}
					instance.Status.RelayServerServiceStatus.Created = true
				}
			}

		case false:
			instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseCreatingRelayServer
			pod := resources.GetRelayServerPod(instance)
			err := r.Create(ctx, pod)
			if err != nil {
				return err
			}
			instance.Status.RelayServerPodStatus.Created = true
		}

	} else {
		if instance.Status.MulticastConnectionPhase == connectionhubv1alpha1.PhysicalInstanceMulticastConnectionPhaseFailed ||
			instance.Status.FederationConnectionPhase == connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseFailed {
			instance.Status.Phase = connectionhubv1alpha1.PhysicalInstancePhaseNotReady
		}
	}

	return nil
}

func (r *PhysicalInstanceReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.PhysicalInstance) error {

	// check submariners.connection-hub.roboscale.io

	instance.Status.Submariner.DeployerStatus.Name = instance.GetSubmarinerDeployerMetadata().Name

	submarinerDeployer := &connectionhubv1alpha1.Submariner{}
	err := r.Get(ctx, instance.GetSubmarinerDeployerMetadata(), submarinerDeployer)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.Submariner.DeployerStatus = connectionhubv1alpha1.DeployerStatus{}
	} else if err != nil {
		return err
	} else {

		instance.Status.Submariner.DeployerStatus.Exists = true
		instance.Status.Submariner.DeployerStatus.Phase = submarinerDeployer.Status.Phase

	}

	// check clusters.submariner.io

	instance.Status.Submariner.ConnectionResources.ClusterStatus.Name = instance.GetSubmarinerClusterMetadata().Name

	submarinerCluster := &brokerv1.Cluster{}
	err = r.Get(ctx, instance.GetSubmarinerClusterMetadata(), submarinerCluster)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.Submariner.ConnectionResources.ClusterStatus = connectionhubv1alpha1.ConnectionResourceStatus{}
	} else if err != nil {
		return err
	} else {

		instance.Status.Submariner.ConnectionResources.ClusterStatus.Exists = true

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
			instance.Status.Submariner.ConnectionResources.EndpointStatus = connectionhubv1alpha1.ConnectionResourceStatus{}
		} else if len(submarinerEndpoints.Items) == 1 {
			instance.Status.Submariner.ConnectionResources.EndpointStatus.Name = submarinerEndpoints.Items[0].Name
			instance.Status.Submariner.ConnectionResources.EndpointStatus.Exists = true
			instance.Status.Subnets.List = submarinerEndpoints.Items[0].Spec.Subnets
			for _, sn := range instance.Status.Subnets.List {
				instance.Status.Subnets.ListInStr = sn + "\n"
			}
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
			instance.Status.Submariner.GatewayConnection = connectionhubv1alpha1.GatewayConnection{}
		} else if len(submarinerGateways.Items) == 1 {
			gw := submarinerGateways.Items[0]
			instance.Status.Submariner.GatewayConnection.GatewayResource = gw.Name
			for _, connection := range gw.Status.Connections {
				if instance.Name == connection.Endpoint.ClusterID {
					instance.Status.Submariner.GatewayConnection.ClusterID = connection.Endpoint.ClusterID
					instance.Status.Submariner.GatewayConnection.Hostname = connection.Endpoint.Hostname
					instance.Status.Submariner.GatewayConnection.ConnectionStatus = connection.Status
				}
			}
		} else {
			return basicErr.New("more than one gateways is listed")
		}

	}

	// check federation member

	federationMember := &connectionhubv1alpha1.FederationMember{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name}, federationMember)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.FederationMember = connectionhubv1alpha1.FederationMemberInstanceStatus{}
	} else if err != nil {
		return err
	} else {

		if !reflect.DeepEqual(instance.Spec.Server, federationMember.Spec.Server) || !reflect.DeepEqual(instance.Spec.Credentials, federationMember.Spec.Credentials) {
			federationMember.Spec.Server = instance.Spec.Server
			federationMember.Spec.Credentials = instance.Spec.Credentials

			err := r.Update(ctx, federationMember)
			if err != nil {
				return err
			}
		}

		instance.Status.FederationMember.Created = true
		instance.Status.FederationMember.Status = federationMember.Status

	}

	// check relay server resources

	relayServerPod := &corev1.Pod{}
	err = r.Get(ctx, instance.GetRelayServerPodMetadata(), relayServerPod)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.RelayServerPodStatus = connectionhubv1alpha1.RelayServerPodStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.RelayServerPodStatus.Created = true
		instance.Status.RelayServerPodStatus.Phase = relayServerPod.Status.Phase
	}

	relayServerService := &corev1.Service{}
	err = r.Get(ctx, instance.GetRelayServerServiceMetadata(), relayServerService)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.RelayServerServiceStatus = connectionhubv1alpha1.RelayServerServiceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.RelayServerServiceStatus.Created = true
		instance.Status.ConnectionURL = "IP:" + strconv.Itoa(int(relayServerService.Spec.Ports[0].NodePort))
	}

	return nil
}

func (r *PhysicalInstanceReconciler) reconcileCreateFederationMember(ctx context.Context, instance *connectionhubv1alpha1.PhysicalInstance) error {

	instance.Status.FederationConnectionPhase = connectionhubv1alpha1.PhysicalInstanceFederationConnectionPhaseConnecting

	member := &connectionhubv1alpha1.FederationMember{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Name,
		},
		Spec: connectionhubv1alpha1.FederationMemberSpec{
			Server: instance.Spec.Server,
			Credentials: connectionhubv1alpha1.PhysicalInstanceCredentials{
				CertificateAuthority: instance.Spec.Credentials.CertificateAuthority,
				ClientKey:            instance.Spec.Credentials.ClientKey,
				ClientCertificate:    instance.Spec.Credentials.ClientCertificate,
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

	instance.Status.FederationMember.Created = true

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
		Watches(
			&source.Kind{Type: &brokerv1.Gateway{}},
			handler.EnqueueRequestsFromMapFunc(r.watchGateways)).
		Owns(&connectionhubv1alpha1.FederationMember{}).
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

func (r *PhysicalInstanceReconciler) watchGateways(o client.Object) []reconcile.Request {

	gateway := o.(*brokerv1.Gateway)

	physicalInstances := &connectionhubv1alpha1.PhysicalInstanceList{}
	err := r.List(context.TODO(), physicalInstances)
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(physicalInstances.Items))
	for i, item := range physicalInstances.Items {

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
