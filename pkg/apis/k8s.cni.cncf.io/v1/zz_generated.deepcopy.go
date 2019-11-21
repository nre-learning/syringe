// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinition) DeepCopyInto(out *NetworkAttachmentDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkAttachmentDefinition.
func (in *NetworkAttachmentDefinition) DeepCopy() *NetworkAttachmentDefinition {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkAttachmentDefinitionList) DeepCopyInto(out *NetworkAttachmentDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NetworkAttachmentDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkAttachmentDefinitionList.
func (in *NetworkAttachmentDefinitionList) DeepCopy() *NetworkAttachmentDefinitionList {
	if in == nil {
		return nil
	}
	out := new(NetworkAttachmentDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkAttachmentDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkList) DeepCopyInto(out *NetworkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NetworkAttachmentDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkList.
func (in *NetworkList) DeepCopy() *NetworkList {
	if in == nil {
		return nil
	}
	out := new(NetworkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetworkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkNames) DeepCopyInto(out *NetworkNames) {
	*out = *in
	if in.ShortNames != nil {
		in, out := &in.ShortNames, &out.ShortNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkNames.
func (in *NetworkNames) DeepCopy() *NetworkNames {
	if in == nil {
		return nil
	}
	out := new(NetworkNames)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkSpec) DeepCopyInto(out *NetworkSpec) {
	*out = *in
	in.Names.DeepCopyInto(&out.Names)
	out.Validation = in.Validation
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkSpec.
func (in *NetworkSpec) DeepCopy() *NetworkSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidation) DeepCopyInto(out *NetworkValidation) {
	*out = *in
	out.OpenAPIV3Schema = in.OpenAPIV3Schema
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidation.
func (in *NetworkValidation) DeepCopy() *NetworkValidation {
	if in == nil {
		return nil
	}
	out := new(NetworkValidation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidationConfig) DeepCopyInto(out *NetworkValidationConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidationConfig.
func (in *NetworkValidationConfig) DeepCopy() *NetworkValidationConfig {
	if in == nil {
		return nil
	}
	out := new(NetworkValidationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidationProperties) DeepCopyInto(out *NetworkValidationProperties) {
	*out = *in
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidationProperties.
func (in *NetworkValidationProperties) DeepCopy() *NetworkValidationProperties {
	if in == nil {
		return nil
	}
	out := new(NetworkValidationProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidationSchema) DeepCopyInto(out *NetworkValidationSchema) {
	*out = *in
	out.Properties = in.Properties
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidationSchema.
func (in *NetworkValidationSchema) DeepCopy() *NetworkValidationSchema {
	if in == nil {
		return nil
	}
	out := new(NetworkValidationSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidationSpec) DeepCopyInto(out *NetworkValidationSpec) {
	*out = *in
	out.Properties = in.Properties
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidationSpec.
func (in *NetworkValidationSpec) DeepCopy() *NetworkValidationSpec {
	if in == nil {
		return nil
	}
	out := new(NetworkValidationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkValidationSpecProperties) DeepCopyInto(out *NetworkValidationSpecProperties) {
	*out = *in
	out.Config = in.Config
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkValidationSpecProperties.
func (in *NetworkValidationSpecProperties) DeepCopy() *NetworkValidationSpecProperties {
	if in == nil {
		return nil
	}
	out := new(NetworkValidationSpecProperties)
	in.DeepCopyInto(out)
	return out
}
