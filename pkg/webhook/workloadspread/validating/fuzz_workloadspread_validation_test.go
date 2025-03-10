/*
Copyright 2021 The Kruise Authors.

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

package validating

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

var (
	fakeScheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(fakeScheme)
	_ = appsv1alpha1.AddToScheme(fakeScheme)
}

func FuzzValidatingWorkloadSpreadFn(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}
		ws.ObjectMeta.SetDeletionTimestamp(nil)

		h := &WorkloadSpreadCreateUpdateHandler{
			Client:  fake.NewClientBuilder().WithObjects(ws).WithScheme(fakeScheme).Build(),
			Decoder: admission.NewDecoder(fakeScheme),
		}

		_ = h.validatingWorkloadSpreadFn(ws)
	})
}

func FuzzValidateWorkloadSpreadSubsets(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)
		f.AllowUnexportedFields()

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}

		elementNums, err := f.GetInt()
		if err != nil {
			return
		}

		subsets := make([]appsv1alpha1.WorkloadSpreadSubset, 0)
		for i := 0; i < elementNums%30; i++ {
			var element appsv1alpha1.WorkloadSpreadSubset
			if err := f.GenerateStruct(&element); err != nil {
				return
			}
			subsets = append(subsets, element)
		}

		fieldPath := &field.Path{}
		if err := f.GenerateStruct(fieldPath); err != nil {
			return
		}

		targetType, err := f.GetInt()
		if err != nil {
			return
		}

		if targetType/2 == 0 {
			cloneSet := &appsv1alpha1.CloneSet{}
			if err := f.GenerateStruct(cloneSet); err != nil {
				return
			}
			_ = validateWorkloadSpreadSubsets(ws, subsets, cloneSet, fieldPath)
		} else {
			sts := &appsv1.StatefulSet{}
			if err := f.GenerateStruct(sts); err != nil {
				return
			}
			_ = validateWorkloadSpreadSubsets(ws, subsets, sts, fieldPath)
		}
	})
}

func FuzzValidateWorkloadSpreadConflict(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)
		f.AllowUnexportedFields()

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}

		elementNums, err := f.GetInt()
		if err != nil {
			return
		}

		others := make([]appsv1alpha1.WorkloadSpread, 0)
		for i := 0; i < elementNums%30; i++ {
			var element appsv1alpha1.WorkloadSpread
			if err := f.GenerateStruct(&element); err != nil {
				return
			}
			others = append(others, element)
		}

		fieldPath := &field.Path{}
		if err := f.GenerateStruct(fieldPath); err != nil {
			return
		}

		_ = validateWorkloadSpreadConflict(ws, others, fieldPath)
	})
}

func FuzzValidateWorkloadSpreadTargetRefUpdate(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)
		f.AllowUnexportedFields()

		targetRef := &appsv1alpha1.TargetReference{}
		if err := f.GenerateStruct(targetRef); err != nil {
			return
		}

		oldTargetRef := &appsv1alpha1.TargetReference{}
		if err := f.GenerateStruct(oldTargetRef); err != nil {
			return
		}

		fieldPath := &field.Path{}
		if err := f.GenerateStruct(fieldPath); err != nil {
			return
		}

		_ = validateWorkloadSpreadTargetRefUpdate(targetRef, oldTargetRef, fieldPath)
	})
}
