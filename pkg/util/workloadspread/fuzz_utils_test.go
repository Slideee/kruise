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

package workloadspread

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	"testing"
)

func FuzzNestedField(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		paths := make([]string, 0)
		if err := f.CreateSlice(&paths); err != nil {
			return
		}

		request, err := f.GetInt()
		if err != nil {
			return
		}

		if request/2 == 0 {
			m := make(map[string]any)
			if err := f.FuzzMap(&m); err != nil {
				return
			}
			_, _, _ = NestedField[any](m, paths...)
		} else {
			m := make([]any, 0)
			if err := f.CreateSlice(&m); err != nil {
				return
			}
			_, _, _ = NestedField[any](m, paths...)
		}
	})
}

func FuzzIsPodSelected(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		filter := &appsv1alpha1.TargetFilter{}
		if err := f.GenerateStruct(filter); err != nil {
			return
		}

		labels := make(map[string]string)
		if err := f.FuzzMap(&labels); err != nil {
			return
		}
		_, _ = IsPodSelected(filter, labels)
	})
}

func FuzzHasPercentSubset(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}
		_ = hasPercentSubset(ws)
	})
}
