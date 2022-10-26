package controllers

import (
	"context"
	basicErr "errors"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
)

// SubmarinerOperatorReconciler reconciles a SubmarinerOperator object
type SubmarinerOperatorReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	RESTConfig *rest.Config
}

//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarineroperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarineroperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=connection-hub.roboscale.io,resources=submarineroperators/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=*
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=*
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=*
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=*
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=*
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=*
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=*
//+kubebuilder:rbac:groups=core,resources=endpoints,verbs=*
//+kubebuilder:rbac:groups=core,resources=events,verbs=*
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=*
//+kubebuilder:rbac:groups=core,resources=pods,verbs=*
//+kubebuilder:rbac:groups=core,resources=services,verbs=*
//+kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=*
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=*,verbs=*
//+kubebuilder:rbac:groups=submariner.io,resources=servicediscoveries,verbs=*
//+kubebuilder:rbac:groups=operator.openshift.io,resources=dnses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.openshift.io,resources=networks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices/restricted,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=submariner.io,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=submariner.io,resources=globalingressips,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multicluster.x-k8s.io,resources=*,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lighthouse.submariner.io,resources=*,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lighthouse.submariner.io,resources=serviceexports,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.openshift.io,resources=service/externalips,verbs=*

func (r *SubmarinerOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger = log.FromContext(ctx)

	instance, err := r.soReconcileGetInstance(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	err = r.soReconcileCheckNode(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.soReconcileCheckDeletion(ctx, instance)
	if err != nil {

		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	err = r.soReconcileCheckStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.soReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.soReconcileCheckResources(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.soReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *SubmarinerOperatorReconciler) soReconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

	switch instance.Status.NamespaceStatus.Created {
	case true:

		switch instance.Status.ChartStatus.Deployed {
		case true:

			switch instance.Status.ChartResourceStatus.Deployed {
			case true:

				instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseDeployed

			case false:

				logger.Info("STATUS: Checking for Submariner Operator resources.")
				instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseCheckingResources

			}

		case false:
			err := r.soReconcileInstallChart(ctx, instance)
			if err != nil {
				return err
			}
		}

	case false:

		err := r.soReconcileCreateNamespace(ctx, instance)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *SubmarinerOperatorReconciler) soReconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {
	operatorNamespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{
		Name: connectionhubv1alpha1.SubmarinerOperatorNamespace,
	}, operatorNamespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.NamespaceStatus.Created = false
	} else if err != nil {
		return err
	}

	ok, err := helmops.CheckIfSubmarinerOperatorExists(*instance, r.RESTConfig)
	if err != nil {
		return err
	} else {
		if ok {
			instance.Status.ChartStatus.Deployed = true
		} else {
			instance.Status.ChartStatus.Deployed = false
		}
	}

	instance.Status.ChartResourceStatus.Deployed = true
	resources := instance.GetResourcesForCheck()
	for _, resource := range resources {
		var obj client.Object

		if resource.GroupVersionKind.Kind == "Deployment" {
			obj = &appsv1.Deployment{}
		} else {
			return basicErr.New("RESOURCE: Operator resource's kind cannot be detected.")
		}

		objKey := resource.ObjectKey
		err := r.Get(ctx, objKey, obj)
		if err != nil {
			instance.Status.ChartResourceStatus.Deployed = false
		}
	}

	return nil
}

func (r *SubmarinerOperatorReconciler) soReconcileCreateNamespace(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

	instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseCreatingNamespace

	operatorNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: connectionhubv1alpha1.SubmarinerOperatorNamespace,
		},
	}

	err := ctrl.SetControllerReference(instance, operatorNamespace, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, operatorNamespace)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner Operator's namespace is created.")
	instance.Status.NamespaceStatus.Created = true

	return nil
}

func (r *SubmarinerOperatorReconciler) soReconcileInstallChart(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

	instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseDeployingChart

	err := helmops.InstallSubmarinerOperatorChart(*instance, r.RESTConfig)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner Operator Helm chart is deployed.")
	instance.Status.ChartResourceStatus.Deployed = true

	return nil
}

func (r *SubmarinerOperatorReconciler) soReconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {
	tenancy := instance.GetTenancySelectors()

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

func (r *SubmarinerOperatorReconciler) soReconcileGetInstance(ctx context.Context, meta types.NamespacedName) (*connectionhubv1alpha1.SubmarinerOperator, error) {
	instance := &connectionhubv1alpha1.SubmarinerOperator{}
	err := r.Get(ctx, meta, instance)
	if err != nil {
		return &connectionhubv1alpha1.SubmarinerOperator{}, err
	}

	return instance, nil
}

func (r *SubmarinerOperatorReconciler) soReconcileUpdateInstanceStatus(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		instanceLV := &connectionhubv1alpha1.SubmarinerOperator{}
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
func (r *SubmarinerOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&connectionhubv1alpha1.SubmarinerOperator{}).
		Owns(&corev1.Namespace{}).
		Complete(r)
}
