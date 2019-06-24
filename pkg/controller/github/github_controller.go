/*
Copyright 2019 TAKAISHI Ryo.

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

package github

import (
	"context"
	"fmt"
	showksv1beta1 "github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/apis/showks/v1beta1"
	"github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/gh"
	"github.com/google/go-github/github"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new GitHub Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	c := gh.NewClient()
	return add(mgr, newReconciler(mgr, c))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, ghClient gh.GitHubClientInterface) reconcile.Reconciler {
	return &ReconcileGitHub{Client: mgr.GetClient(), scheme: mgr.GetScheme(), ghClient: ghClient}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("github-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to GitHub
	err = c.Watch(&source.Kind{Type: &showksv1beta1.GitHub{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by GitHub - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &showksv1beta1.GitHub{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileGitHub{}

// ReconcileGitHub reconciles a GitHub object
type ReconcileGitHub struct {
	client.Client
	scheme   *runtime.Scheme
	ghClient gh.GitHubClientInterface
}

// Reconcile reads that state of the cluster for a GitHub object and makes changes based on the state read
// and what is in the GitHub.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=showks.cloudnativedays.jp,resources=githubs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=showks.cloudnativedays.jp,resources=githubs/status,verbs=get;update;patch
func (r *ReconcileGitHub) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	fmt.Println("Reconcile")

	instance := &showksv1beta1.GitHub{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	_, err = r.ghClient.GetRepository(instance.Spec.Repository.Org, instance.Spec.Repository.Name)
	if err != nil {
		if _, ok := err.(*gh.NotFoundError); ok {
			repoSpec := &github.Repository{Name: &instance.Spec.Repository.Name}
			repo, err := r.ghClient.CreateRepository(instance.Spec.Repository.Org, repoSpec)
			if err != nil {
				return reconcile.Result{}, err
			}

			instance.Status.ID = *repo.ID

			if err := r.Status().Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
