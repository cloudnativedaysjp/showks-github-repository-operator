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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitHubRepositorySpec defines the desired state of GitHub
type GitHubRepositorySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Org               string                 `json:"org"`
	Name              string                 `json:"name"`
	BranchProtections []BranchProtectionSpec `json:"branchProtections,omitempty"`
	Webhooks          []WebhookSpec          `json:"webhooks,omitempty"`
	Collaborators     []CollaboratorSpec     `json:"collaborators,omitempty"`
}

type CollaboratorSpec struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

type BranchProtectionSpec struct {
	BranchName                 string                         `json:"branchName,omitempty"`
	RequiredStatusChecks       RequiredStatusChecksSpec       `json:"requiredStatusChecks,omitempty"`
	EnforceAdmin               bool                           `json:"enforceAdmin,omitempty"`
	RequiredPullRequestReviews RequiredPullRequestReviewsSpec `json:"requiredPullRequestReviews,omitempty"`
	Restrictions               RestrictionsSpec               `json:"restrictions,omitempty"`
}

type RequiredPullRequestReviewsSpec struct {
}

type RequiredStatusChecksSpec struct {
	Strict   bool     `json:"strict,omitempty"`
	Contexts []string `json:"contexts,omitempty"`
}

type RestrictionsSpec struct {
	Users []string `json:"users,omitempty"`
	Teams []string `json:"teams,omitempty"`
}

type WebhookSpec struct {
	Config WebhookConfigSpec `json:"config,omitempty"`
	Events []string          `json:"events,omitempty"`
	Active bool              `json:"active,omitempty"`
}

type WebhookConfigSpec struct {
	Url         string `json:"url"`
	ContentType string `json:"contentType"`
}

// GitHubStatus defines the observed state of GitHub
type GitHubStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ID int64
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GitHub is the Schema for the githubs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type GitHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubRepositorySpec `json:"spec,omitempty"`
	Status GitHubStatus         `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GitHubList contains a list of GitHub
type GitHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHub{}, &GitHubList{})
}
