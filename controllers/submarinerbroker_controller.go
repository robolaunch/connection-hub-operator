package controllers

import (
	"context"
	basicErr "errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
)

// SubmarinerBrokerReconciler reconciles a SubmarinerBroker object
type SubmarinerBrokerReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	RESTConfig *rest.Config
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarinerbrokers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarinerbrokers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarinerbrokers/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices/restricted,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lighthouse.submariner.io,resources=*,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multicluster.x-k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=endpoints,verbs=get;list;watch;create;update;patch;delete

var logger logr.Logger

func (r *SubmarinerBrokerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.smbReconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.smbReconcileCheckNode(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.smbReconcileCheckDeletion(ctx, instance)
	if err != nil {

		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	err = r.smbReconcileCheckStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.smbReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.smbReconcileCheckResources(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.smbReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	// TODO: Check Submariner Broker namespace existance

	switch instance.Status.Phase {
	case connectionhubv1alpha1.SubmarinerBrokerPhaseNotExists:
		err := r.smbReconcileInstallChart(ctx, instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	if ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig); err != nil {
		return err
	} else {
		if ok {
			instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseRunning
		} else {
			instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseNotExists
		}
	}

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	submarinerBrokerFinalizer := "connection-hub.roboscale.io/finalizer"

	if instance.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {
			controllerutil.AddFinalizer(instance, submarinerBrokerFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

	} else {

		if controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {
			if ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig); err != nil {
				return err
			} else {
				if ok {
					err := helmops.UninstallSubmarinerBrokerChart(*instance, r.RESTConfig)
					if err != nil {
						return err
					}
				}
			}

			controllerutil.RemoveFinalizer(instance, submarinerBrokerFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

		return errors.NewNotFound(schema.GroupResource{
			Group:    instance.GetObjectKind().GroupVersionKind().Group,
			Resource: instance.GetObjectKind().GroupVersionKind().Kind,
		}, instance.Name)
	}

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileInstallChart(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	err := helmops.InstallSubmarinerBrokerChart(*instance, r.RESTConfig)
	if err != nil {
		return err
	}

	instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseRunning

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckIfChartExisted(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) (bool, error) {

	ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig)
	if err != nil {
		return false, err
	}

	if !ok {
		instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseRunning
		return false, nil
	}

	return true, nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {
	tenancy := connectionhubv1alpha1.GetTenancySelectorsForSMB(*instance)

	instance.Status.NodeInfo.Selectors = make(map[string]string)

	if tenancy.RobolaunchCloudInstance != "" {
		instance.Status.NodeInfo.Selectors[connectionhubv1alpha1.RobolaunchCloudInstanceLabelKey] = tenancy.RobolaunchCloudInstance
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

func (r *SubmarinerBrokerReconciler) smbReconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.SubmarinerBroker, error) {
	instance := &connectionhubv1alpha1.SubmarinerBroker{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.SubmarinerBroker{}, err
	}

	return instance, nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.SubmarinerBroker{}
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
func (r *SubmarinerBrokerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.SubmarinerBroker{}).
		Complete(r)
}
