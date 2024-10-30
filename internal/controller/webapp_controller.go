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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	samplecrdv1 "github.com/kushal2705/sample-crd/api/v1"
)

// WebAppReconciler reconciles a WebApp object
type WebAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sample-crd.example.com,resources=webapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sample-crd.example.com,resources=webapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sample-crd.example.com,resources=webapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *WebAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the WebApp instance
	webapp := &samplecrdv1.WebApp{}
	err := r.Get(ctx, req.NamespacedName, webapp)
	if err != nil {
		log.Error(err, "unable to fetch WebApp")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Update environment variables
	updatedEnv := updateEnvVariables(webapp.Spec.Env)
	if len(updatedEnv) != len(webapp.Spec.Env) {
		webapp.Spec.Env = updatedEnv
		err = r.Update(ctx, webapp)
		if err != nil {
			log.Error(err, "unable to update WebApp")
			return ctrl.Result{}, err
		}
		log.Info("Updated WebApp environment variables")
	}
	return ctrl.Result{}, nil
}

func updateEnvVariables(env []samplecrdv1.EnvVar) []samplecrdv1.EnvVar {
	updatedEnv := make([]samplecrdv1.EnvVar, 0)
	for _, e := range env {
		if e.Name == "UPDATE_ME" {
			e.Value = fmt.Sprintf("updated-%s", e.Value)
		}
		updatedEnv = append(updatedEnv, e)
	}
	return updatedEnv
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplecrdv1.WebApp{}).
		Named("webapp").
		Complete(r)
}
