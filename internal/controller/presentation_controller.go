/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	presentationsv1alpha1 "github.com/operator-sdk-demo/slide-controller/api/v1alpha1"
	"github.com/operator-sdk-demo/slide-controller/pkg/mdparser"
	"github.com/operator-sdk-demo/slide-controller/pkg/mdrender"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// PresentationReconciler reconciles a Presentation object
type PresentationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=presentations.haavard.dev,resources=presentations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=presentations.haavard.dev,resources=presentations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=presentations.haavard.dev,resources=presentations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Presentation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *PresentationReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	presentation := &presentationsv1alpha1.Presentation{}
	err := r.Get(ctx, req.NamespacedName, presentation)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("not found, ignoring")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if presentation.GetDeletionTimestamp() != nil {
		// The resource is being deleted. Skip the reconciliation.
		logger.Info("Skipping reconcile as the resource is being deleted")
		return reconcile.Result{}, nil
	}

	rendered := mdrender.RenderMarkdown(&presentation.Spec)

	if err = r.SetupPresentation(ctx, req, presentation, rendered); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PresentationReconciler) SetupPresentation(
	ctx context.Context,
	req ctrl.Request,
	presentation *presentationsv1alpha1.Presentation,
	rendered string,
) error {
	logger := log.FromContext(ctx)

	configMap, deployment, service := mdparser.CreateMarkdownParser(req.Name, rendered)

	if err := ctrl.SetControllerReference(presentation, configMap, r.Scheme); err != nil {
		logger.Error(err, "unable to set controller reference for configmap")
		return err
	}
	if err := ctrl.SetControllerReference(presentation, deployment, r.Scheme); err != nil {
		logger.Error(err, "unable to set controller reference for deployment")
		return err
	}
	if err := ctrl.SetControllerReference(presentation, service, r.Scheme); err != nil {
		logger.Error(err, "unable to set controller reference for service")
		return err
	}

	// Apply ConfigMap
	existingConfigMap := &corev1.ConfigMap{}
	err := r.Get(
		ctx,
		client.ObjectKey{Name: configMap.Name, Namespace: configMap.Namespace},
		existingConfigMap,
	)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new ConfigMap
			if err := r.Create(ctx, configMap); err != nil {
				logger.Error(err, "unable to create ConfigMap")
				return err
			}
		} else {
			logger.Error(err, "unable to get ConfigMap")
			return err
		}
	} else {
		// Update existing ConfigMap
		existingConfigMap.Data = configMap.Data
		if err := r.Update(ctx, existingConfigMap); err != nil {
			logger.Error(err, "unable to update ConfigMap")
			return err
		}
	}

	// Apply Deployment
	existingDeployment := &appsv1.Deployment{}
	err = r.Get(
		ctx,
		client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace},
		existingDeployment,
	)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new Deployment
			if err := r.Create(ctx, deployment); err != nil {
				logger.Error(err, "unable to create Deployment")
				return err
			}
		} else {
			logger.Error(err, "unable to get Deployment")
			return err
		}
	} else {
		// Update existing Deployment
		existingDeployment.Spec = deployment.Spec
		if err := r.Update(ctx, existingDeployment); err != nil {
			logger.Error(err, "unable to update Deployment")
			return err
		}
	}

	// Apply Service
	existingService := &corev1.Service{}
	err = r.Get(
		ctx,
		client.ObjectKey{Name: service.Name, Namespace: service.Namespace},
		existingService,
	)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new Service
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "unable to create Service")
				return err
			}
		} else {
			logger.Error(err, "unable to get Service")
			return err
		}
	} else {
		// Update existing Service
		existingService.Spec = service.Spec
		if err := r.Update(ctx, existingService); err != nil {
			logger.Error(err, "unable to update Service")
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PresentationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&presentationsv1alpha1.Presentation{}).
		Complete(r)
}
