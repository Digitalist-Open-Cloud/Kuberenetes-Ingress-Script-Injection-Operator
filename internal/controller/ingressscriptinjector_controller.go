/*
Copyright 2024 by Digitalist Open Cloud.

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
)

// Constants for annotations
const (
	headEnd            = "digitalist.cloud/add-script-head-end"
	headStart          = "digitalist.cloud/add-script-head-start"
	bodyStart          = "digitalist.cloud/add-script-body-start"
	bodyEnd            = "digitalist.cloud/add-script-body-end"
	nginxConfigSnippet = "nginx.ingress.kubernetes.io/configuration-snippet"
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
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile handles the main reconciliation logic
func (r *IngressScriptInjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var ingress networkingv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		logger.Error(err, "Unable to fetch Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Collect configuration snippets based on annotations
	snippets := []string{}

	// Utility function to retrieve the ConfigMap script
	getScriptFromConfigMap := func(configMapName string) (string, error) {
		var configMap corev1.ConfigMap
		if err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: req.Namespace}, &configMap); err != nil {
			return "", err
		}
		script, exists := configMap.Data["script"]
		if !exists {
			return "", fmt.Errorf("ConfigMap %s does not contain 'script' key", configMapName)
		}
		return script, nil
	}

	// Check for each annotation and prepare snippets based on position
	if configMapName, ok := ingress.Annotations[headEnd]; ok {
		if script, err := getScriptFromConfigMap(configMapName); err == nil {
			snippets = append(snippets, fmt.Sprintf("sub_filter '</head>' '%s</head>';", script))
		} else {
			logger.Error(err, "Unable to fetch ConfigMap for head end script", "ConfigMap", configMapName)
		}
	}

	if configMapName, ok := ingress.Annotations[headStart]; ok {
		if script, err := getScriptFromConfigMap(configMapName); err == nil {
			snippets = append(snippets, fmt.Sprintf("sub_filter '<head>' '<head>%s';", script))
		} else {
			logger.Error(err, "Unable to fetch ConfigMap for head beginning script", "ConfigMap", configMapName)
		}
	}

	if configMapName, ok := ingress.Annotations[bodyStart]; ok {
		if script, err := getScriptFromConfigMap(configMapName); err == nil {
			snippets = append(snippets, fmt.Sprintf("sub_filter '<body>' '<body>%s';", script))
		} else {
			logger.Error(err, "Unable to fetch ConfigMap for body beginning script", "ConfigMap", configMapName)
		}
	}

	if configMapName, ok := ingress.Annotations[bodyEnd]; ok {
		if script, err := getScriptFromConfigMap(configMapName); err == nil {
			snippets = append(snippets, fmt.Sprintf("sub_filter '</body>' '%s</body>';", script))
		} else {
			logger.Error(err, "Unable to fetch ConfigMap for body end script", "ConfigMap", configMapName)
		}
	}

	// If no snippets to add, exit reconciliation early
	if len(snippets) == 0 {
		return ctrl.Result{}, nil
	}

	// Combine snippets and update the annotation
	combinedSnippet := ""
	if existingSnippet, hasSnippet := ingress.Annotations[nginxConfigSnippet]; hasSnippet {
		combinedSnippet = fmt.Sprintf("%s\n%s", existingSnippet, combinedSnippet)
	}
	combinedSnippet += fmt.Sprintf("%s\n", snippets)

	// Ensure Ingress annotations map is initialized
	if ingress.Annotations == nil {
		ingress.Annotations = make(map[string]string)
	}
	ingress.Annotations[nginxConfigSnippet] = combinedSnippet

	// Apply the update to the Ingress
	if err := r.Update(ctx, &ingress); err != nil {
		logger.Error(err, "Failed to update Ingress with configuration snippets")
		return ctrl.Result{}, err
	}

	logger.Info("Successfully updated Ingress with configuration snippets", "Ingress", req.NamespacedName)
	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager
func (r *IngressScriptInjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}). // Watch for changes to Ingress resources
		Complete(r)
}
