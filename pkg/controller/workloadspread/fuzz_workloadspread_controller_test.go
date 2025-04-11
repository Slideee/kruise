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
	"context"
	"encoding/json"
	"k8s.io/klog/v2"
	"strconv"
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	fuzzutils "github.com/openkruise/kruise/pkg/util/fuzz"
	wsutil "github.com/openkruise/kruise/pkg/util/workloadspread"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	fakeScheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(fakeScheme)
	_ = appsv1alpha1.AddToScheme(fakeScheme)
}

func FuzzPatchFavoriteSubsetMetadataToPod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cf := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := cf.GenerateStruct(pod); err != nil {
			return
		}

		// Cleanup deletion timestamp when no finalizers exist
		if pod.GetDeletionTimestamp() != nil && len(pod.GetFinalizers()) == 0 {
			pod.SetDeletionTimestamp(nil)
		}

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := cf.GenerateStruct(ws); err != nil {
			return
		}
		if ws.GetAnnotations() == nil {
			ws.Annotations = make(map[string]string)
		}

		ignore, err := cf.GetBool()
		if err != nil {
			return
		}
		ws.GetAnnotations()[IgnorePatchExistingPodsAnnotation] = strconv.FormatBool(ignore)

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := cf.GenerateStruct(subset); err != nil {
			return
		}

		if err := fuzzutils.GenerateWorkloadSpreadSubsetPatch(cf, subset); err != nil {
			return
		}

		r := &ReconcileWorkloadSpread{
			Client: fake.NewClientBuilder().
				WithScheme(fakeScheme).
				WithObjects(pod.DeepCopy()).
				Build(),
		}

		if err := r.patchFavoriteSubsetMetadataToPod(pod, ws, subset); err != nil {
			return
		}

		updatedPod := &corev1.Pod{}
		if err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(pod), updatedPod); err != nil {
			t.Errorf("Failed to fetch updated Pod: %v", err)
			return
		}

		// Calculate expected annotations/labels inline
		expectedAnnotations := make(map[string]string)
		expectedLabels := make(map[string]string)

		// Copy original pod's metadata
		if originalAnnotations := pod.GetAnnotations(); originalAnnotations != nil {
			for k, v := range originalAnnotations {
				expectedAnnotations[k] = v
			}
		}
		if originalLabels := pod.GetLabels(); originalLabels != nil {
			for k, v := range originalLabels {
				expectedLabels[k] = v
			}
		}

		// Apply subset patch if applicable
		if subset.Patch.Raw != nil && ws.Annotations[IgnorePatchExistingPodsAnnotation] != "true" {
			var patchField map[string]interface{}
			if err := json.Unmarshal(subset.Patch.Raw, &patchField); err == nil {
				if metadata, ok := patchField["metadata"].(map[string]interface{}); ok {
					if labels, ok := metadata["labels"].(map[string]interface{}); ok {
						for k, v := range labels {
							if vs, ok := v.(string); ok {
								expectedLabels[k] = vs
							}
						}
					}
					if annotations, ok := metadata["annotations"].(map[string]interface{}); ok {
						for k, v := range annotations {
							if vs, ok := v.(string); ok {
								expectedAnnotations[k] = vs
							}
						}
					}
				}
			}
		}

		injectWS := &wsutil.InjectWorkloadSpread{
			Name:   ws.Name,
			Subset: subset.Name,
		}
		injectWSBytes, _ := json.Marshal(injectWS)
		expectedAnnotations[wsutil.MatchedWorkloadSpreadSubsetAnnotations] = string(injectWSBytes)

		// Check annotations
		actualAnnotations := updatedPod.GetAnnotations()
		for k, v := range expectedAnnotations {
			actualValue, ok := actualAnnotations[k]
			if !ok {
				t.Errorf("Missing annotation %s (expected: %s)", k, v)
			} else if actualValue != v {
				t.Errorf("Annotation mismatch for %s: expected %s, got %s", k, v, actualValue)
			}
			klog.Infof("FuzzPatchFavoriteSubsetMetadataToPod 1")
		}

		// Check labels
		actualLabels := updatedPod.GetLabels()
		for k, v := range expectedLabels {
			actualValue, ok := actualLabels[k]
			if !ok {
				t.Errorf("Missing label %s (expected: %s)", k, v)
			} else if actualValue != v {
				t.Errorf("Label mismatch for %s: expected %s, got %s", k, v, actualValue)
			}
			klog.Infof("FuzzPatchFavoriteSubsetMetadataToPod 2")
		}
	})
}

func FuzzPodPreferredScore(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cf := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := cf.GenerateStruct(pod); err != nil {
			return
		}

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := cf.GenerateStruct(subset); err != nil {
			return
		}

		if err := fuzzutils.GenerateWorkloadSpreadSubsetPatch(cf, subset); err != nil {
			return
		}

		_ = podPreferredScore(subset, pod)
	})
}
