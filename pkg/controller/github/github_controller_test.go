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
	"github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/gh"
	"github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	"testing"
	"time"

	showksv1beta1 "github.com/cloudnativedaysjp/showks-github-repository-operator/pkg/apis/showks/v1beta1"
	"github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
var repoKey = types.NamespacedName{Name: "foo", Namespace: "default"}

const timeout = time.Second * 5

var org = "sample"
var repoName = "aaa"
var repoID int64 = 1296269

func newGitHubClientMock(controller *gomock.Controller) gh.GitHubClientInterface {
	c := mock_gh.NewMockGitHubClientInterface(controller)
	repoSpec := &github.Repository{Name: &repoName}
	repoResp := &github.Repository{ID: &repoID, Name: &repoName}
	c.EXPECT().CreateRepository(org, repoSpec).Return(repoResp, nil).Times(1)
	c.EXPECT().DeleteRepository(org, repoName).Return(nil).Times(1)

	firstGetRepo := c.EXPECT().GetRepository(org, repoName).Return(nil, &gh.NotFoundError{}).Times(1)
	c.EXPECT().GetRepository(org, repoName).Return(repoResp, nil).After(firstGetRepo).Times(1)

	c.EXPECT().AddCollaborator(org, repoName, "alice", "admin").Return(nil).Times(1)
	firstGetPermission := c.EXPECT().GetPermissionLevel(org, repoName, "alice").Return("", &gh.NotFoundError{}).Times(1)
	c.EXPECT().GetPermissionLevel(org, repoName, "alice").Return("admin", nil).After(firstGetPermission).Times(1)

	team := "showks-members"
	protectionReq := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   false,
			Contexts: nil,
		},
		EnforceAdmins: false,
		Restrictions: &github.BranchRestrictionsRequest{
			Teams: []string{team},
		},
	}
	protection := &github.Protection{Restrictions: &github.BranchRestrictions{Teams: []*github.Team{{Name: &team}}}}
	c.EXPECT().UpdateBranchProtection(org, repoName, "master", protectionReq).Return(protection, nil).Times(2)

	active := true
	var hookID int64 = 12345
	hook := &github.Hook{
		Config: map[string]interface{}{
			"url":          "https://example.com",
			"content_type": "json",
		},
		Events: []string{"push"},
		Active: &active,
	}
	c.EXPECT().CreateHook(org, repoName, hook).Return(hook, nil).Times(1)

	firstListHook := c.EXPECT().ListHook(org, repoName).Return([]*github.Hook{}, nil).Times(1)
	hooks := []*github.Hook{
		{
			Config: map[string]interface{}{
				"url":          "https://example.com",
				"content_type": "json",
			},
			Events: []string{"push"},
			Active: &active,
			ID:     &hookID,
		},
	}
	c.EXPECT().ListHook(org, repoName).Return(hooks, nil).After(firstListHook).Times(1)
	c.EXPECT().UpdateHook(org, repoName, hookID, hook).Return(hook, nil).Times(1)
	return c
}

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &showksv1beta1.GitHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: showksv1beta1.GitHubSpec{
			Repository: showksv1beta1.RepositorySpec{
				Org:  org,
				Name: repoName,
				Collaborators: []showksv1beta1.CollaboratorSpec{
					{
						Name:       "alice",
						Permission: "admin",
					},
				},
			},
			BranchProtections: []showksv1beta1.BranchProtectionSpec{
				{
					BranchName: "master",
					RequiredStatusChecks: showksv1beta1.RequiredStatusChecksSpec{
						Strict:   false,
						Contexts: nil,
					},
					EnforceAdmin: false,
					Restrictions: showksv1beta1.RestrictionsSpec{
						Teams: []string{"showks-members"},
					},
				},
			},
			Webhooks: []showksv1beta1.WebhookSpec{
				{
					Active: true,
					Config: showksv1beta1.WebhookConfigSpec{
						Url:         "https://example.com",
						ContentType: "json",
					},
					Events: []string{"push"},
				},
			},
		},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ghClient := newGitHubClientMock(mockCtrl)
	recFn, requests := SetupTestReconcile(newReconciler(mgr, ghClient))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the GitHub object and expect the Reconcile and Deployment to be created
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

	//deploy := &appsv1.Deployment{}
	//g.Eventually(func() error { return c.Get(context.TODO(), depKey, deploy) }, timeout).
	//	Should(gomega.Succeed())

	repo := &showksv1beta1.GitHub{}
	// Delete the Deployment and expect Reconcile to be called for Deployment deletion
	g.Eventually(func() error { return c.Get(context.TODO(), repoKey, repo) }, timeout).
		Should(gomega.Succeed())
	g.Expect(repo.Status.ID).To(gomega.Equal(repoID))

	g.Expect(c.Delete(context.TODO(), repo)).NotTo(gomega.HaveOccurred())
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
	// Manually delete Deployment since GC isn't enabled in the test control plane
	g.Eventually(func() error { return c.Delete(context.TODO(), repo) }, timeout).
		Should(gomega.MatchError("githubs.showks.cloudnativedays.jp \"foo\" not found"))

}
