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
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func FuzzPodPreferredScore(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := f.GenerateStruct(pod); err != nil {
			return
		}

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := f.GenerateStruct(subset); err != nil {
			return
		}
		_ = podPreferredScore(subset, pod)
	})
}

func FuzzMatchesSubsetRequiredAndToleration(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := f.GenerateStruct(pod); err != nil {
			return
		}

		node := &corev1.Node{}
		if err := f.GenerateStruct(node); err != nil {
			return
		}

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := f.GenerateStruct(subset); err != nil {
			return
		}

		_, _ = matchesSubsetRequiredAndToleration(pod, node, subset)
	})
}

func FuzzMatchSubset(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := f.GenerateStruct(pod); err != nil {
			return
		}

		node := &corev1.Node{}
		if err := f.GenerateStruct(node); err != nil {
			return
		}

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := f.GenerateStruct(subset); err != nil {
			return
		}

		missingReplicas, err := f.GetInt()
		if err != nil {
			return
		}

		_, _, _ = matchesSubset(pod, node, subset, missingReplicas)
	})
}
