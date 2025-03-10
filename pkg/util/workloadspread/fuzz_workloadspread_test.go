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
	"encoding/json"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func FuzzInjectWorkloadSpreadIntoPod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := f.GenerateStruct(pod); err != nil {
			return
		}

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}

		subsetName, err := f.GetString()
		if err != nil {
			return
		}

		generatedUID, err := f.GetString()
		if err != nil {
			return
		}

		inject, err := injectWorkloadSpreadIntoPod(ws, pod, subsetName, generatedUID)
		if err != nil {
			return
		}

		if inject {
			validateInjection(t, ws.GetName(), subsetName, generatedUID, pod)
		}
	})
}

func validateInjection(t *testing.T, wsName, subsetName, generatedUID string, pod *corev1.Pod) {
	injectWS := &InjectWorkloadSpread{
		Name:   wsName,
		Subset: subsetName,
		UID:    generatedUID,
	}
	expectedBy, _ := json.Marshal(injectWS)
	actual := pod.GetAnnotations()[MatchedWorkloadSpreadSubsetAnnotations]

	if actual != string(expectedBy) {
		t.Errorf("expected %s, got %s", string(expectedBy), actual)
	}
}

func FuzzGetReplicasFromObject(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)
		path, err := f.GetString()
		if err != nil {
			return
		}

		obj := &unstructured.Unstructured{}
		err = f.GenerateStruct(obj)
		if err != nil {
			return
		}

		_, _ = GetReplicasFromObject(obj, path)
	})
}
