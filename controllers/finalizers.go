package controllers

import (
	"context"
	"time"

	submv1alpha1 "github.com/robolaunch/connection-hub-operator/api/external/submariner/v1alpha1"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
	corev1 "k8s.io/api/core/v1"
	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SubmarinerBrokerReconciler) reconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	submarinerBrokerFinalizer := "submarinerbroker.connection-hub.roboscale.io/finalizer"

	if instance.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {
			controllerutil.AddFinalizer(instance, submarinerBrokerFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

	} else {

		if controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {

			err := r.waitForChartDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForNamespaceDeletion(ctx, instance)
			if err != nil {
				return err
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

// TODO: confirm chart deletions by checking chart's resources
func (r *SubmarinerBrokerReconciler) waitForChartDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {
	instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseUninstallingChart
	err := r.smbReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	if ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig); err != nil {
		return err
	} else {
		if ok {
			err := helmops.UninstallSubmarinerBrokerChart(*instance, r.RESTConfig)
			if err != nil {
				return err
			}

			logger.Info("FINALIZER: Broker chart is uninstalled.")

		}
	}

	return nil
}

func (r *SubmarinerBrokerReconciler) waitForNamespaceDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

	submarinerBrokerNamespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, *instance.GetNamespaceMetadata(), submarinerBrokerNamespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Submariner Broker namespace is being deleted.")
		err := r.Delete(ctx, submarinerBrokerNamespaceQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.SubmarinerBrokerPhaseTerminatingNamespace
		err = r.smbReconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    submarinerBrokerNamespaceQuery.GroupVersionKind().Group,
			Version:  submarinerBrokerNamespaceQuery.GroupVersionKind().Version,
			Resource: "namespaces",
		})
		submarinerBrokerNamespaceWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetNamespaceMetadata().Name,
		})

		defer submarinerBrokerNamespaceWatcher.Stop()

		submarinerBrokerNamespaceDeleted := false
		for {
			if !submarinerBrokerNamespaceDeleted {
				select {
				case event := <-submarinerBrokerNamespaceWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Submariner Broker namespace is deleted gracefully.")
						submarinerBrokerNamespaceDeleted = true
					}
				}
			} else {
				break
			}

		}
	}
	return nil

}

func (r *SubmarinerOperatorReconciler) reconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

	submarinerOperatorFinalizer := "submarineroperator.connection-hub.roboscale.io/finalizer"

	if instance.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(instance, submarinerOperatorFinalizer) {
			controllerutil.AddFinalizer(instance, submarinerOperatorFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

	} else {

		if controllerutil.ContainsFinalizer(instance, submarinerOperatorFinalizer) {

			err := r.waitForChartDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForNamespaceDeletion(ctx, instance)
			if err != nil {
				return err
			}

			controllerutil.RemoveFinalizer(instance, submarinerOperatorFinalizer)
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

// TODO: confirm chart deletions by checking chart's resources
func (r *SubmarinerOperatorReconciler) waitForChartDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {
	instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseUninstallingChart
	err := r.soReconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	if ok, err := helmops.CheckIfSubmarinerOperatorExists(*instance, r.RESTConfig); err != nil {
		return err
	} else {
		if ok {

			err = helmops.UninstallSubmarinerOperatorChart(*instance, r.RESTConfig)
			if err != nil {
				return err
			}

			logger.Info("FINALIZER: Operator chart is uninstalled.")
		}
	}

	return nil
}

func (r *SubmarinerOperatorReconciler) waitForNamespaceDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

	submarinerOperatorNamespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, *instance.GetNamespaceMetadata(), submarinerOperatorNamespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Submariner Operator namespace is being deleted.")
		err := r.Delete(ctx, submarinerOperatorNamespaceQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.SubmarinerOperatorPhaseTerminatingNamespace
		err = r.soReconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    submarinerOperatorNamespaceQuery.GroupVersionKind().Group,
			Version:  submarinerOperatorNamespaceQuery.GroupVersionKind().Version,
			Resource: "namespaces",
		})
		submarinerOperatorNamespaceWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetNamespaceMetadata().Name,
		})

		defer submarinerOperatorNamespaceWatcher.Stop()

		submarinerOperatorNamespaceDeleted := false
		for {
			if !submarinerOperatorNamespaceDeleted {
				select {
				case event := <-submarinerOperatorNamespaceWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Submariner Operator namespace is deleted gracefully.")
						submarinerOperatorNamespaceDeleted = true
					}
				}
			} else {
				break
			}

		}
	}
	return nil

}

func (r *SubmarinerReconciler) reconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	submarinerFinalizer := "submariner.connection-hub.roboscale.io/finalizer"

	if instance.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(instance, submarinerFinalizer) {
			controllerutil.AddFinalizer(instance, submarinerFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

	} else {

		if controllerutil.ContainsFinalizer(instance, submarinerFinalizer) {

			err := r.waitForSubmarinerCRDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForSubmarinerOperatorDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForSubmarinerBrokerDeletion(ctx, instance)
			if err != nil {
				return err
			}

			controllerutil.RemoveFinalizer(instance, submarinerFinalizer)
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

func (r *SubmarinerReconciler) waitForSubmarinerCRDeletion(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	submarinerCRQuery := &submv1alpha1.Submariner{}
	err := r.Get(ctx, *instance.GetSubmarinerCustomResourceMetadata(), submarinerCRQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Submariner CR is being deleted.")
		err := r.Delete(ctx, submarinerCRQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseTerminatingSubmarinerCR
		err = r.submarinerReconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    submarinerCRQuery.GroupVersionKind().Group,
			Version:  submarinerCRQuery.GroupVersionKind().Version,
			Resource: "submariners",
		})
		submarinerWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetSubmarinerCustomResourceMetadata().Name,
		})

		defer submarinerWatcher.Stop()

		submarinerCRDeleted := false
		for {
			if !submarinerCRDeleted {
				select {
				case event := <-submarinerWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Submariner CR is deleted gracefully.")
						submarinerCRDeleted = true
					}
				}
			} else {
				break
			}

		}
	}
	return nil

}

func (r *SubmarinerReconciler) waitForSubmarinerOperatorDeletion(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	submarinerOperatorQuery := &connectionhubv1alpha1.SubmarinerOperator{}
	err := r.Get(ctx, *instance.GetSubmarinerOperatorMetadata(), submarinerOperatorQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Submariner Operator is being deleted.")
		err := r.Delete(ctx, submarinerOperatorQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseTerminatingSubmarinerOperator
		err = r.submarinerReconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    submarinerOperatorQuery.GroupVersionKind().Group,
			Version:  submarinerOperatorQuery.GroupVersionKind().Version,
			Resource: "submarineroperators",
		})
		submarinerOperatorWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetSubmarinerOperatorMetadata().Name,
		})

		defer submarinerOperatorWatcher.Stop()

		submarinerOperatorDeleted := false
		for {
			if !submarinerOperatorDeleted {
				select {
				case event := <-submarinerOperatorWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Submariner Operator is deleted gracefully.")
						submarinerOperatorDeleted = true
					}
				}
			} else {
				break
			}

		}
	}
	return nil

}

func (r *SubmarinerReconciler) waitForSubmarinerBrokerDeletion(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	submarinerBrokerQuery := &connectionhubv1alpha1.SubmarinerBroker{}
	err := r.Get(ctx, *instance.GetSubmarinerBrokerMetadata(), submarinerBrokerQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Submariner Broker is being deleted.")
		err := r.Delete(ctx, submarinerBrokerQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseTerminatingSubmarinerBroker
		err = r.submarinerReconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    submarinerBrokerQuery.GroupVersionKind().Group,
			Version:  submarinerBrokerQuery.GroupVersionKind().Version,
			Resource: "submarinerbrokers",
		})
		submarinerBrokerWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetSubmarinerBrokerMetadata().Name,
		})

		defer submarinerBrokerWatcher.Stop()

		submarinerBrokerDeleted := false
		for {
			if !submarinerBrokerDeleted {
				select {
				case event := <-submarinerBrokerWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Submariner Broker is deleted gracefully.")
						submarinerBrokerDeleted = true
					}
				}
			} else {
				break
			}

		}
	}
	return nil

}

func (r *FederationOperatorReconciler) reconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	submarinerBrokerFinalizer := "federationoperator.connection-hub.roboscale.io/finalizer"

	if instance.DeletionTimestamp.IsZero() {

		if !controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {
			controllerutil.AddFinalizer(instance, submarinerBrokerFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				return err
			}
		}

	} else {

		if controllerutil.ContainsFinalizer(instance, submarinerBrokerFinalizer) {

			err := r.waitForFederatedTypeCRDDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForFederatedTypeConfigsDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForChartDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForFederatedCoreCRDDeletion(ctx, instance)
			if err != nil {
				return err
			}

			err = r.waitForNamespaceDeletion(ctx, instance)
			if err != nil {
				return err
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

func (r *FederationOperatorReconciler) waitForFederatedTypeCRDDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	logger.Info("FINALIZER: Federated Type CRDs are being deleted.")

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseDeletingFederatedTypeCRDs
	err := r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	apiGroup := "types.kubefed.io"

	doesItemExists := func() bool {
		crds := &extensionsv1.CustomResourceDefinitionList{}
		_ = r.List(ctx, crds)
		for _, crd := range crds.Items {
			if crd.Spec.Group == apiGroup {
				return true
			}
		}
		return false
	}

	for doesItemExists() {
		crdList := &extensionsv1.CustomResourceDefinitionList{}
		_ = r.List(ctx, crdList)

		for _, crd := range crdList.Items {
			if crd.Spec.Group == apiGroup {
				_ = r.Delete(ctx, &crd)
			}
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func (r *FederationOperatorReconciler) waitForFederatedTypeConfigsDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	logger.Info("FINALIZER: FederatedTypeConfigs are being deleted.")

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseDeletingFederatedTypeConfigs
	err := r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
		Group:    "core.kubefed.io",
		Version:  "v1beta1",
		Resource: "federatedtypeconfigs",
	})

	doesItemExists := func() bool {
		tcList, _ := resourceInterface.List(ctx, metav1.ListOptions{})
		return len(tcList.Items) != 0
	}

	for doesItemExists() {
		typeConfigList, _ := resourceInterface.List(ctx, metav1.ListOptions{})
		for _, tc := range typeConfigList.Items {
			_ = resourceInterface.Delete(ctx, tc.GetName(), metav1.DeleteOptions{})
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

func (r *FederationOperatorReconciler) waitForChartDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	logger.Info("FINALIZER: Federation Operator helm chart is being uninstalled.")

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseUninstallingChart
	err := r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	if ok, err := helmops.CheckIfFederationOperatorExists(*instance, r.RESTConfig); err != nil {
		return err
	} else {
		if ok {

			err = helmops.UninstallFederationOperatorChart(*instance, r.RESTConfig)
			if err != nil {
				return err
			}

			logger.Info("FINALIZER: Federation Operator chart is uninstalled.")
		}
	}

	return nil
}

func (r *FederationOperatorReconciler) waitForFederatedCoreCRDDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	logger.Info("FINALIZER: Federated Core CRDs are being deleted.")

	instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseDeletingFederatedCoreCRDs
	err := r.reconcileUpdateInstanceStatus(ctx, instance)
	if err != nil {
		return err
	}

	apiGroups := []string{"core.kubefed.io", "scheduling.kubefed.io"}

	doesItemExists := func() bool {
		crds := &extensionsv1.CustomResourceDefinitionList{}
		_ = r.List(ctx, crds)
		for _, crd := range crds.Items {
			for _, ag := range apiGroups {
				if crd.Spec.Group == ag {
					return true
				}
			}
		}
		return false
	}

	for doesItemExists() {
		crdList := &extensionsv1.CustomResourceDefinitionList{}
		_ = r.List(ctx, crdList)

		for _, crd := range crdList.Items {
			for _, apiGroup := range apiGroups {
				if crd.Spec.Group == apiGroup {
					_ = r.Delete(ctx, &crd)
				}
			}
		}

		time.Sleep(2 * time.Second)

	}

	return nil

}

func (r *FederationOperatorReconciler) waitForNamespaceDeletion(ctx context.Context, instance *connectionhubv1alpha1.FederationOperator) error {

	federationOperatorNamespaceQuery := &corev1.Namespace{}
	err := r.Get(ctx, *instance.GetNamespaceMetadata(), federationOperatorNamespaceQuery)
	if err != nil && errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	} else {
		logger.Info("FINALIZER: Federation Operator namespace is being deleted.")
		err := r.Delete(ctx, federationOperatorNamespaceQuery)
		if err != nil {
			return err
		}

		instance.Status.Phase = connectionhubv1alpha1.FederationOperatorPhaseTerminatingNamespace
		err = r.reconcileUpdateInstanceStatus(ctx, instance)
		if err != nil {
			return err
		}

		resourceInterface := r.DynamicClient.Resource(schema.GroupVersionResource{
			Group:    federationOperatorNamespaceQuery.GroupVersionKind().Group,
			Version:  federationOperatorNamespaceQuery.GroupVersionKind().Version,
			Resource: "namespaces",
		})
		federationOperatorNamespaceWatcher, err := resourceInterface.Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + instance.GetNamespaceMetadata().Name,
		})

		defer federationOperatorNamespaceWatcher.Stop()

		federationOperatorNamespaceDeleted := false
		for {
			if !federationOperatorNamespaceDeleted {
				select {
				case event := <-federationOperatorNamespaceWatcher.ResultChan():

					if event.Type == watch.Deleted {
						logger.Info("FINALIZER: Federation Operator namespace is deleted gracefully.")
						federationOperatorNamespaceDeleted = true
					}
				}
			} else {
				break
			}

		}
	}

	return nil
}
