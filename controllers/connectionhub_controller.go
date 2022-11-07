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
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	"github.com/robolaunch/connection-hub-operator/controllers/pkg/resources"
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

	err = r.reconcileCheckNode(ctx, instance)
	if err != nil {
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

					instance.Status.Phase = connectionhubv1alpha1.ConnectionHubPhaseReadyForOperation

					// create federation host if cloud instance

					// create cloud instance if physical instance

				default:

					// wait for federation to be ready

				}

			case false:

				err := r.reconcileCreateFederation(ctx, instance)
				if err != nil {
					return err
				}

			}

		default:

			// wait for submariner to be ready

		}

	case false:

		err := r.reconcileCreateSubmariner(ctx, instance)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *ConnectionHubReconciler) reconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {

	// check submariner
	submariner := &connectionhubv1alpha1.Submariner{}
	err := r.Get(ctx, types.NamespacedName{Name: instance.GetSubmarinerMetadata().Name}, submariner)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.Submariner = connectionhubv1alpha1.SubmarinerInstanceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.Submariner.Created = true
		instance.Status.Submariner.Phase = submariner.Status.Phase
	}

	// check federation
	federation := &connectionhubv1alpha1.FederationOperator{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.GetFederationMetadata().Name}, federation)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.Federation = connectionhubv1alpha1.FederationInstanceStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.Federation.Created = true
		instance.Status.Federation.Phase = federation.Status.Phase
	}

	return nil
}

func (r *ConnectionHubReconciler) reconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {
	tenancy := instance.GetTenancySelectors()

	instance.Status.NodeInfo.Selectors = make(map[string]string)

	if tenancy.RobolaunchCloudInstance != "" {
		instance.Status.NodeInfo.Selectors[connectionhubv1alpha1.RobolaunchCloudInstanceLabelKey] = tenancy.RobolaunchCloudInstance
	}

	if tenancy.RobolaunchPhysicalInstance != "" {
		instance.Status.NodeInfo.Selectors[connectionhubv1alpha1.RobolaunchPhysicalInstanceLabelKey] = tenancy.RobolaunchPhysicalInstance
	}

	requirements := []labels.Requirement{}
	requirementsMap := instance.Status.NodeInfo.Selectors
	for k, v := range requirementsMap {
		newReq, err := labels.NewRequirement(k, selection.In, []string{v})
		if err != nil {
			return err
		}
		requirements = append(requirements, *newReq)
	}

	nodeSelector := labels.NewSelector().Add(requirements...)

	nodes := &corev1.NodeList{}
	err := r.List(ctx, nodes, &client.ListOptions{
		LabelSelector: nodeSelector,
	})
	if err != nil {
		return err
	}

	if len(nodes.Items) == 0 {
		return basicErr.New("no nodes are listed with node selector")
	} else if len(nodes.Items) > 1 {
		return basicErr.New("multiple nodes are listed, no specific target")
	}

	node := nodes.Items[0]
	instance.Status.NodeInfo.Name = node.Name

	return nil
}

func (r *ConnectionHubReconciler) reconcileCreateSubmariner(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {

	instance.Status.Phase = connectionhubv1alpha1.ConnectionHubPhaseSubmarinerSettingUp

	submariner := resources.GetSubmariner(instance)

	err := ctrl.SetControllerReference(instance, submariner, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, submariner)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner is created.")

	instance.Status.Submariner.Created = true

	return nil
}

func (r *ConnectionHubReconciler) reconcileCreateFederation(ctx context.Context, instance *connectionhubv1alpha1.ConnectionHub) error {

	instance.Status.Phase = connectionhubv1alpha1.ConnectionHubPhaseFederationSettingUp

	federation := resources.GetSubmariner(instance)

	err := ctrl.SetControllerReference(instance, federation, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, federation)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Federation is created.")

	instance.Status.Federation.Created = true

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
