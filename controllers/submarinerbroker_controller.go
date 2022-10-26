package controllers

import (
	"context"
	"encoding/base64"
	basicErr "errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	err = r.reconcileCheckDeletion(ctx, instance)
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

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Second * 10,
	}, nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckStatus(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	switch instance.Status.NamespaceStatus.Created {
	case true:

		switch instance.Status.ChartStatus.Deployed {
		case true:

			switch instance.Status.ChartResourceStatus.Deployed {
			case true:

				instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseDeployed

			case false:

				instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseCheckingResources

			}

		case false:

			err := r.smbReconcileInstallChart(ctx, instance)
			if err != nil {
				return err
			}

		}

	case false:

		err := r.smbReconcileCreateNamespace(ctx, instance)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckResources(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	brokerNamespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{
		Name: instance.GetNamespaceMetadata().Name,
	}, brokerNamespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		instance.Status.NamespaceStatus.Created = false
	} else if err != nil {
		return err
	}

	if ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig); err != nil {
		return err
	} else {
		if ok {
			instance.Status.ChartStatus.Deployed = true
		} else {
			instance.Status.ChartStatus.Deployed = false
		}
	}

	// get token and ca
	err = r.smbReconcileUpdateBrokerInfo(ctx, instance)
	if err != nil {
		return err
	}

	if instance.Status.Broker.BrokerToken != "" && instance.Status.Broker.BrokerCA != "" {
		instance.Status.ChartResourceStatus.Deployed = true
	}

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCreateNamespace(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseCreatingNamespace

	brokerNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.GetNamespaceMetadata().Name,
		},
	}

	err := ctrl.SetControllerReference(instance, brokerNamespace, r.Scheme)
	if err != nil {
		return err
	}

	err = r.Create(ctx, brokerNamespace)
	if err != nil {
		return err
	}

	logger.Info("STATUS: Submariner Broker's namespace is created.")
	instance.Status.NamespaceStatus.Created = true

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileInstallChart(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseDeployingChart

	err := helmops.InstallSubmarinerBrokerChart(*instance, r.RESTConfig)
	if err != nil {
		return err
	}

	instance.Status.ChartStatus.Deployed = true

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileUpdateBrokerInfo(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	secretList := &corev1.SecretList{}
	err := r.List(ctx, secretList, &client.ListOptions{
		Namespace: instance.GetNamespaceMetadata().Name,
	})
	if err != nil {
		return err
	}

	for _, secret := range secretList.Items {
		if val, ok := secret.GetAnnotations()["kubernetes.io/service-account.name"]; ok {

			if val == instance.GetNamespaceMetadata().Name+"-client" {

				if token, ok := secret.Data["token"]; ok {
					instance.Status.Broker.BrokerToken = string(token[:])
				}

				if ca, ok := secret.Data["ca.crt"]; ok {
					instance.Status.Broker.BrokerCA = base64.StdEncoding.EncodeToString([]byte(string(ca[:])))
				}

				break

			}
		}
	}

	instance.Status.Broker.BrokerURL = instance.Spec.BrokerURL

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileCheckNode(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {
	tenancy := instance.GetTenancySelectors()

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
