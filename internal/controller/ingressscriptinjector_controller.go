/*
Copyright 2024.

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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ingressv1alpha1 "github.com/Digitalist-Open-Cloud/Kuberenetes-Ingress-Script-Injection-Operator/api/v1alpha1"
)

// Constants for annotations
const (
	headEndAnnotation            = "digitalist.cloud/add-script-head-end"
	nginxConfigSnippetAnnotation = "nginx.ingress.kubernetes.io/configuration-snippet"
)

// IngressScriptInjectorReconciler reconciles a IngressScriptInjector object
type IngressScriptInjectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ingress.digitalist.cloud,resources=ingressscriptinjectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ingress.digitalist.cloud,resources=ingressscriptinjectors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ingress.digitalist.cloud,resources=ingressscriptinjectors/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;update;patch

// Reconcile handles the main reconciliation logic
func (r *IngressScriptInjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Ingress resource
	var ingress networkingv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		logger.Error(err, "Unable to fetch Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the Ingress has the add-script-head-end annotation
	configMapName, ok := ingress.Annotations[headEndAnnotation]
	if !ok {
		// If the annotation is missing, no processing is needed
		return ctrl.Result{}, nil
	}

	// Fetch the ConfigMap referenced by the add-script-head-end annotation
	var configMap corev1.ConfigMap
	if err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: req.Namespace}, &configMap); err != nil {
		logger.Error(err, "Unable to fetch ConfigMap", "ConfigMap", configMapName)
		return ctrl.Result{}, err
	}

	// Get the script content from the ConfigMap data
	script, exists := configMap.Data["script"]
	if !exists {
		logger.Info("ConfigMap does not contain 'script' key", "ConfigMap", configMapName)
		return ctrl.Result{}, nil
	}

	// Construct the NGINX configuration snippet
	configurationSnippet := fmt.Sprintf("sub_filter '</head>' '%s</head>';", script)

	// Check if there is an existing configuration snippet to merge with
	existingSnippet, hasSnippet := ingress.Annotations[nginxConfigSnippetAnnotation]
	if hasSnippet {
		configurationSnippet = fmt.Sprintf("%s\n%s", existingSnippet, configurationSnippet)
	}

	// Update the Ingress with the new or merged configuration snippet
	if ingress.Annotations == nil {
		ingress.Annotations = make(map[string]string)
	}
	ingress.Annotations[nginxConfigSnippetAnnotation] = configurationSnippet

	// Apply the update to the Ingress
	if err := r.Update(ctx, &ingress); err != nil {
		logger.Error(err, "Failed to update Ingress with configuration snippet")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully updated Ingress with configuration snippet", "Ingress", req.NamespacedName)
	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager
func (r *IngressScriptInjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ingressv1alpha1.IngressScriptInjector{}).
		Named("ingressscriptinjector").
		Complete(r)
}
