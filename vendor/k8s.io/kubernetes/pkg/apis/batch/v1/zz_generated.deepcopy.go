// +build !ignore_autogenerated

/*
Copyright 2017 The Kubernetes Authors.

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	api_v1 "k8s.io/kubernetes/pkg/api/v1"
	reflect "reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1_Job, InType: reflect.TypeOf(&Job{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1_JobCondition, InType: reflect.TypeOf(&JobCondition{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1_JobList, InType: reflect.TypeOf(&JobList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1_JobSpec, InType: reflect.TypeOf(&JobSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_v1_JobStatus, InType: reflect.TypeOf(&JobStatus{})},
	)
}

func DeepCopy_v1_Job(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*Job)
		out := out.(*Job)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*meta_v1.ObjectMeta)
		}
		if err := DeepCopy_v1_JobSpec(&in.Spec, &out.Spec, c); err != nil {
			return err
		}
		if err := DeepCopy_v1_JobStatus(&in.Status, &out.Status, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_v1_JobCondition(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*JobCondition)
		out := out.(*JobCondition)
		*out = *in
		out.LastProbeTime = in.LastProbeTime.DeepCopy()
		out.LastTransitionTime = in.LastTransitionTime.DeepCopy()
		return nil
	}
}

func DeepCopy_v1_JobList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*JobList)
		out := out.(*JobList)
		*out = *in
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]Job, len(*in))
			for i := range *in {
				if err := DeepCopy_v1_Job(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func DeepCopy_v1_JobSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*JobSpec)
		out := out.(*JobSpec)
		*out = *in
		if in.Parallelism != nil {
			in, out := &in.Parallelism, &out.Parallelism
			*out = new(int32)
			**out = **in
		}
		if in.Completions != nil {
			in, out := &in.Completions, &out.Completions
			*out = new(int32)
			**out = **in
		}
		if in.ActiveDeadlineSeconds != nil {
			in, out := &in.ActiveDeadlineSeconds, &out.ActiveDeadlineSeconds
			*out = new(int64)
			**out = **in
		}
		if in.Selector != nil {
			in, out := &in.Selector, &out.Selector
			if newVal, err := c.DeepCopy(*in); err != nil {
				return err
			} else {
				*out = newVal.(*meta_v1.LabelSelector)
			}
		}
		if in.ManualSelector != nil {
			in, out := &in.ManualSelector, &out.ManualSelector
			*out = new(bool)
			**out = **in
		}
		if err := api_v1.DeepCopy_v1_PodTemplateSpec(&in.Template, &out.Template, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_v1_JobStatus(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*JobStatus)
		out := out.(*JobStatus)
		*out = *in
		if in.Conditions != nil {
			in, out := &in.Conditions, &out.Conditions
			*out = make([]JobCondition, len(*in))
			for i := range *in {
				if err := DeepCopy_v1_JobCondition(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		if in.StartTime != nil {
			in, out := &in.StartTime, &out.StartTime
			*out = new(meta_v1.Time)
			**out = (*in).DeepCopy()
		}
		if in.CompletionTime != nil {
			in, out := &in.CompletionTime, &out.CompletionTime
			*out = new(meta_v1.Time)
			**out = (*in).DeepCopy()
		}
		return nil
	}
}
