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
	submv1alpha1 "github.com/submariner-io/submariner-operator/apis/submariner/v1alpha1"
)

// SubmarinerReconciler reconciles a Submariner object
type SubmarinerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submariners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submariners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submariners/finalizers,verbs=update

// +kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarinerbrokers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarineroperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=submariner.io,resources=submariners,verbs=get;list;watch;create;update;patch;delete

func (r *SubmarinerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.submarinerReconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.submarinerReconcileCheckNode(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.submarinerReconcileCheckStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.submarinerReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.submarinerReconcileCheckResources(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.submarinerReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SubmarinerReconciler) submarinerReconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {
	switch instance.Status.BrokerStatus.Created {
	case true:

		switch instance.Status.BrokerStatus.Phase {
		case connectionhubv1alpha1.SubmarinerBrokerPhaseDeployed:

			switch instance.Status.OperatorStatus.Created {
			case true:

				switch instance.Status.OperatorStatus.Phase {
				case connectionhubv1alpha1.SubmarinerOperatorPhaseDeployed:

					// create submariner custom resource

					instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseReadyToConnect

				}

			case false:
				err := r.submarinerReconcileCreateOperator(ctx, instance)
				if err != nil {
					return err
				}
			}

		}

	case false:
		err := r.submarinerReconcileCreateBroker(ctx, instance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SubmarinerReconciler) submarinerReconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	submarinerBrokerQuery := &connectionhubv1alpha1.SubmarinerBroker{}
	err := r.Get(ctx, *instance.GetSubmarinerBrokerMetadata(), submarinerBrokerQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.BrokerStatus = connectionhubv1alpha1.BrokerStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.BrokerStatus.Created = true
		instance.Status.BrokerStatus.Phase = submarinerBrokerQuery.Status.Phase
		instance.Status.BrokerStatus.Status = submarinerBrokerQuery.Status
	}

	submarinerOperatorQuery := &connectionhubv1alpha1.SubmarinerOperator{}
	err = r.Get(ctx, *instance.GetSubmarinerOperatorMetadata(), submarinerOperatorQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.OperatorStatus = connectionhubv1alpha1.OperatorStatus{}
	} else if err != nil {
		return err
	} else {
		instance.Status.OperatorStatus.Created = true
		instance.Status.OperatorStatus.Phase = submarinerOperatorQuery.Status.Phase
	}

	return nil
}

func (r *SubmarinerReconciler) submarinerReconcileCreateBroker(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {
	instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseCreatingBroker

	submarinerBroker := resources.GetSubmarinerBroker(instance)

	err := ctrl.SetControllerReference(instance, submarinerBroker, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, submarinerBroker)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner broker is created.")

	instance.Status.BrokerStatus.Created = true

	return nil
}

func (r *SubmarinerReconciler) submarinerReconcileCreateOperator(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {
	instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseCreatingOperator

	submarinerOperator := resources.GetSubmarinerOperator(instance)

	err := ctrl.SetControllerReference(instance, submarinerOperator, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, submarinerOperator)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner operator is created.")

	instance.Status.OperatorStatus.Created = true

	return nil
}

func (r *SubmarinerReconciler) submarinerReconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {
	tenancy := connectionhubv1alpha1.GetTenancySelectorsForSubmariner(*instance)

	instance.Status.NodeInfo.Selectors = make(map[string]string)

	if tenancy.RobolaunchCloudInstance != "" {
		instance.Status.NodeInfo.Selectors[connectionhubv1alpha1.RobolaunchCloudInstanceLabelKey] = tenancy.RobolaunchCloudInstance
	}

	if tenancy.RobolaunchPhysicalInstance != "" {
		instance.Status.NodeInfo.Selectors[connectionhubv1alpha1.RobolaunchPhysicalInstanceLabelKey] = tenancy.RobolaunchCloudInstance
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

func (r *SubmarinerReconciler) submarinerReconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.Submariner, error) {
	instance := &connectionhubv1alpha1.Submariner{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.Submariner{}, err
	}

	return instance, nil
}

func (r *SubmarinerReconciler) submarinerReconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.Submariner{}
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
func (r *SubmarinerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.Submariner{}).
		Owns(&connectionhubv1alpha1.SubmarinerBroker{}).
		Owns(&connectionhubv1alpha1.SubmarinerOperator{}).
		Owns(&submv1alpha1.Submariner{}).
		Complete(r)
}
