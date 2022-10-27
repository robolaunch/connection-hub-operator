package controllers

import (
	"context"

	submv1alpha1 "github.com/robolaunch/connection-hub-operator/api/external/submariner/v1alpha1"
	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
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
			if ok, err := helmops.CheckIfSubmarinerBrokerExists(*instance, r.RESTConfig); err != nil {
				return err
			} else {
				if ok {
					err := helmops.UninstallSubmarinerBrokerChart(*instance, r.RESTConfig)
					if err != nil {
						return err
					}

					logger.Info("FINALIZER: Broker chart is uninstalled.")

					// err = r.smbReconcileDeleteNamespace(ctx, instance)
					// if err != nil {
					// 	return err
					// }

					// logger.Info("FINALIZER: Broker namespace is deleted.")
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
			if ok, err := helmops.CheckIfSubmarinerOperatorExists(*instance, r.RESTConfig); err != nil {
				return err
			} else {
				if ok {

					err = helmops.UninstallSubmarinerOperatorChart(*instance, r.RESTConfig)
					if err != nil {
						return err
					}

					logger.Info("FINALIZER: Operator chart is uninstalled.")

					// err = r.soReconcileDeleteNamespace(ctx, instance)
					// if err != nil {
					// 	return err
					// }

					// logger.Info("FINALIZER: Namespace is deleted.")

				}
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

			submarinerCRQuery := &submv1alpha1.Submariner{}
			err := r.Get(ctx, *instance.GetSubmarinerCustomResourceMetadata(), submarinerCRQuery)
			if err != nil && errors.IsNotFound(err) {
				// do nothing
			} else if err != nil {
				return err
			} else {
				logger.Info("FINALIZER: Submariner CR is being deleted.")
				err := r.Delete(ctx, submarinerCRQuery)
				if err != nil {
					return err
				}

				instance.Status.Phase = connectionhubv1alpha1.SubmarinerPhaseTerminatingSubmariner
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

func (r *SubmarinerReconciler) reconcileDeleteCustomResource(ctx context.Context, instance *connectionhubv1alpha1.Submariner) error {

	return nil
}
