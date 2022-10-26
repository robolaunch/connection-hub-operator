package controllers

import (
	"context"

	connectionhubv1alpha1 "github.com/robolaunch/connection-hub-operator/api/v1alpha1"
	helmops "github.com/robolaunch/connection-hub-operator/controllers/pkg/helm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SubmarinerBrokerReconciler) smbReconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {

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

func (r *SubmarinerOperatorReconciler) soReconcileCheckDeletion(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {

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

func (r *SubmarinerOperatorReconciler) soReconcileDeleteNamespace(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerOperator) error {
	operatorNamespace := &corev1.Namespace{}

	err := r.Get(ctx, types.NamespacedName{
		Name: connectionhubv1alpha1.SubmarinerOperatorNamespace,
	}, operatorNamespace)
	if err != nil {
		return err
	} else {
		err = r.Delete(ctx, operatorNamespace)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SubmarinerBrokerReconciler) smbReconcileDeleteNamespace(ctx context.Context, instance *connectionhubv1alpha1.SubmarinerBroker) error {
	brokerNamespace := &corev1.Namespace{}

	err := r.Get(ctx, types.NamespacedName{
		Name: connectionhubv1alpha1.SubmarinerBrokerNamespace,
	}, brokerNamespace)
	if err != nil {
		return err
	} else {
		err = r.Delete(ctx, brokerNamespace)
		if err != nil {
			return err
		}
	}

	return nil
}
