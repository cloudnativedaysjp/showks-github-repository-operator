// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionSpec) DeepCopyInto(out *BranchProtectionSpec) {
	*out = *in
	in.RequiredStatusChecks.DeepCopyInto(&out.RequiredStatusChecks)
	out.RequiredPullRequestReviews = in.RequiredPullRequestReviews
	in.Restrictions.DeepCopyInto(&out.Restrictions)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionSpec.
func (in *BranchProtectionSpec) DeepCopy() *BranchProtectionSpec {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollaboratorSpec) DeepCopyInto(out *CollaboratorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollaboratorSpec.
func (in *CollaboratorSpec) DeepCopy() *CollaboratorSpec {
	if in == nil {
		return nil
	}
	out := new(CollaboratorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitHub) DeepCopyInto(out *GitHub) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitHub.
func (in *GitHub) DeepCopy() *GitHub {
	if in == nil {
		return nil
	}
	out := new(GitHub)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GitHub) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitHubList) DeepCopyInto(out *GitHubList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GitHub, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitHubList.
func (in *GitHubList) DeepCopy() *GitHubList {
	if in == nil {
		return nil
	}
	out := new(GitHubList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GitHubList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitHubRepositorySpec) DeepCopyInto(out *GitHubRepositorySpec) {
	*out = *in
	if in.BranchProtections != nil {
		in, out := &in.BranchProtections, &out.BranchProtections
		*out = make([]BranchProtectionSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Webhooks != nil {
		in, out := &in.Webhooks, &out.Webhooks
		*out = make([]WebhookSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Collaborators != nil {
		in, out := &in.Collaborators, &out.Collaborators
		*out = make([]CollaboratorSpec, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitHubRepositorySpec.
func (in *GitHubRepositorySpec) DeepCopy() *GitHubRepositorySpec {
	if in == nil {
		return nil
	}
	out := new(GitHubRepositorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitHubStatus) DeepCopyInto(out *GitHubStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitHubStatus.
func (in *GitHubStatus) DeepCopy() *GitHubStatus {
	if in == nil {
		return nil
	}
	out := new(GitHubStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequiredPullRequestReviewsSpec) DeepCopyInto(out *RequiredPullRequestReviewsSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequiredPullRequestReviewsSpec.
func (in *RequiredPullRequestReviewsSpec) DeepCopy() *RequiredPullRequestReviewsSpec {
	if in == nil {
		return nil
	}
	out := new(RequiredPullRequestReviewsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequiredStatusChecksSpec) DeepCopyInto(out *RequiredStatusChecksSpec) {
	*out = *in
	if in.Contexts != nil {
		in, out := &in.Contexts, &out.Contexts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequiredStatusChecksSpec.
func (in *RequiredStatusChecksSpec) DeepCopy() *RequiredStatusChecksSpec {
	if in == nil {
		return nil
	}
	out := new(RequiredStatusChecksSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestrictionsSpec) DeepCopyInto(out *RestrictionsSpec) {
	*out = *in
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Teams != nil {
		in, out := &in.Teams, &out.Teams
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestrictionsSpec.
func (in *RestrictionsSpec) DeepCopy() *RestrictionsSpec {
	if in == nil {
		return nil
	}
	out := new(RestrictionsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookConfigSpec) DeepCopyInto(out *WebhookConfigSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookConfigSpec.
func (in *WebhookConfigSpec) DeepCopy() *WebhookConfigSpec {
	if in == nil {
		return nil
	}
	out := new(WebhookConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookSpec) DeepCopyInto(out *WebhookSpec) {
	*out = *in
	out.Config = in.Config
	if in.Events != nil {
		in, out := &in.Events, &out.Events
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookSpec.
func (in *WebhookSpec) DeepCopy() *WebhookSpec {
	if in == nil {
		return nil
	}
	out := new(WebhookSpec)
	in.DeepCopyInto(out)
	return out
}
