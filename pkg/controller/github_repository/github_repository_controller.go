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

package github_repository

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
var finalizerName = "finalizer.github.showks.cloudnativedays.jp"

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

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.setFinalizer(instance); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		return r.runFinalizer(instance)
	}

	_, err = r.ghClient.GetRepository(instance.Spec.Org, instance.Spec.Name)
	if err != nil {
		if _, ok := err.(*gh.NotFoundError); ok {
			repoSpec := &github.Repository{Name: &instance.Spec.Name}
			repo, err := r.ghClient.CreateRepository(instance.Spec.Org, repoSpec)
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

	err = r.ReconcileCollaborators(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.ReconcileBranchProtection(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.ReconcileWebHook(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGitHub) setFinalizer(instance *showksv1beta1.GitHub) error {
	fmt.Println("setFinalizer")
	if !containsString(instance.ObjectMeta.Finalizers, finalizerName) {
		instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(context.Background(), instance); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileGitHub) runFinalizer(instannce *showksv1beta1.GitHub) (reconcile.Result, error) {
	fmt.Println("runFinalizer")
	if containsString(instannce.ObjectMeta.Finalizers, finalizerName) {
		if err := r.deleteExternalDependency(instannce); err != nil {
			return reconcile.Result{}, err
		}

		instannce.ObjectMeta.Finalizers = removeString(instannce.ObjectMeta.Finalizers, finalizerName)
		if err := r.Update(context.Background(), instannce); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileGitHub) deleteExternalDependency(instance *showksv1beta1.GitHub) error {
	fmt.Println("deleteExternalDependency")
	return r.ghClient.DeleteRepository(instance.Spec.Org, instance.Spec.Name)
}

func (r *ReconcileGitHub) ReconcileCollaborators(instance *showksv1beta1.GitHub) error {
	for _, collaborator := range instance.Spec.Collaborators {
		_, err := r.ghClient.GetPermissionLevel(instance.Spec.Org, instance.Spec.Name, collaborator.Name)
		if _, ok := err.(*gh.NotFoundError); ok {
			err = r.ghClient.AddCollaborator(instance.Spec.Org, instance.Spec.Name, collaborator.Name, collaborator.Permission)
			if err != nil {
				return err
			}
		} else {
			return err

		}
	}

	return nil
}

func (r *ReconcileGitHub) ReconcileBranchProtection(instance *showksv1beta1.GitHub) error {
	log.Info("Reconcile: BranchProtection")
	owner := instance.Spec.Org
	repo := instance.Spec.Name
	for _, bpSpec := range instance.Spec.BranchProtections {
		bp := &github.ProtectionRequest{
			RequiredStatusChecks: &github.RequiredStatusChecks{
				Strict:   bpSpec.RequiredStatusChecks.Strict,
				Contexts: bpSpec.RequiredStatusChecks.Contexts,
			},
			EnforceAdmins: bpSpec.EnforceAdmin,
			Restrictions: &github.BranchRestrictionsRequest{
				Users: bpSpec.Restrictions.Users,
				Teams: bpSpec.Restrictions.Teams,
			},
		}

		log.Info("Updating BranchProtection...", "branch", bpSpec.BranchName)
		_, err := r.ghClient.UpdateBranchProtection(owner, repo, bpSpec.BranchName, bp)
		if err != nil {
			return err
		}
		log.Info("Updated BranchProtection.", "branch", bpSpec.BranchName)
	}

	return nil
}

// リポジトリに同じURLのWebHookがなければ作成、あれば更新する
// リポジトリにあってSpecにないWebHookは削除する
func (r *ReconcileGitHub) ReconcileWebHook(instance *showksv1beta1.GitHub) error {
	log.Info("Reconcile: WebHook")
	owner := instance.Spec.Org
	repo := instance.Spec.Name

	existsHooks, err := r.ghClient.ListHook(owner, repo)
	if err != nil {
		return err
	}

	for _, whSpec := range instance.Spec.Webhooks {
		existsHook := findHookByUrl(existsHooks, whSpec.Config.Url)

		if existsHook != nil {
			log.Info("WebHook already exists on this repository", "url", whSpec.Config.Url)
			err := r.updateWebHook(owner, repo, *existsHook.ID, whSpec)
			if err != nil {
				return err
			}
		} else {
			log.Info("WebHook does'nt exists on this repository", "url", whSpec.Config.Url)
			err := r.createWebHook(owner, repo, whSpec)
			if err != nil {
				return err
			}
		}
	}

	for _, existsHook := range existsHooks {
		exist := false
		for _, whSpec := range instance.Spec.Webhooks {
			if whSpec.Config.Url == existsHook.Config["url"].(string) {
				exist = true
			}
		}

		if !exist {
			err := r.deleteWebHook(owner, repo, *existsHook.ID, existsHook)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ReconcileGitHub) createWebHook(owner string, repo string, hook showksv1beta1.WebhookSpec) error {
	log.Info("Creating WebHook...", "url", hook.Config.Url)

	wh := &github.Hook{
		Config: map[string]interface{}{
			"url":          hook.Config.Url,
			"content_type": hook.Config.ContentType,
		},
		Events: hook.Events,
		Active: &hook.Active,
	}

	_, err := r.ghClient.CreateHook(owner, repo, wh)
	if err != nil {
		return err
	}

	log.Info("Created WebHook.", "url", hook.Config.Url)

	return nil
}

func (r *ReconcileGitHub) updateWebHook(owner string, repo string, id int64, hook showksv1beta1.WebhookSpec) error {
	log.Info("Updating WebHook...", "url", hook.Config.Url)

	wh := &github.Hook{
		Config: map[string]interface{}{
			"url":          hook.Config.Url,
			"content_type": hook.Config.ContentType,
		},
		Events: hook.Events,
		Active: &hook.Active,
	}

	_, err := r.ghClient.UpdateHook(owner, repo, id, wh)
	if err != nil {
		return err
	}

	log.Info("Updated WebHook.", "url", hook.Config.Url)

	return nil
}

func (r *ReconcileGitHub) deleteWebHook(owner string, repo string, id int64, hook *github.Hook) error {
	log.Info("Deleting WebHook...", "url", hook.Config["url"].(string))

	err := r.ghClient.DeleteHook(owner, repo, *hook.ID)
	if err != nil {
		return err
	}

	log.Info("Deleted WebHook...", "url", hook.Config["url"].(string))

	return nil
}

func findHookByUrl(hooks []*github.Hook, url string) *github.Hook {
	for _, hook := range hooks {
		if url == hook.Config["url"].(string) {
			return hook
		}
	}

	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
