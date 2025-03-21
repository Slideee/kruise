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
	"strconv"
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	fakeScheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(fakeScheme)
	_ = appsv1alpha1.AddToScheme(fakeScheme)
}

// FuzzPatchFavoriteSubsetMetadataToPod tests the metadata patching logic for Pods in WorkloadSpread
// This fuzzer validates both valid and invalid patch scenarios
func FuzzPatchFavoriteSubsetMetadataToPod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cf := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := cf.GenerateStruct(pod); err != nil {
			return
		}

		// Cleanup deletion timestamp when no finalizers exist.
		// fake client will refuse to create obj with metadata.deletionTimestamp but no finalizers
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

		// Randomly decide whether to ignore existing pods
		ignore, err := cf.GetBool()
		if err != nil {
			return
		}
		ws.GetAnnotations()[IgnorePatchExistingPodsAnnotation] = strconv.FormatBool(ignore)

		subset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := cf.GenerateStruct(subset); err != nil {
			return
		}

		// Generate valid or invalid JSON patch randomly
		if err := generatePatch(cf, subset); err != nil {
			return
		}

		r := &ReconcileWorkloadSpread{
			Client: fake.NewClientBuilder().WithScheme(fakeScheme).WithObjects(pod).Build(),
		}

		_ = r.patchFavoriteSubsetMetadataToPod(pod, ws, subset)
	})
}

// FuzzPodPreferredScore tests the scoring logic for Pod scheduling preference
// This fuzzer validates score calculation under various subset configurations
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

		// Generate valid or invalid JSON patch randomly
		if err := generatePatch(cf, subset); err != nil {
			return
		}

		_ = podPreferredScore(subset, pod)
	})
}

// generatePatch creates either valid labeled JSON patches or random byte payloads
// to test both successful merges and error handling scenarios.
func generatePatch(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.WorkloadSpreadSubset) error {
	// 50% chance to generate structured label patch
	if isStructured, _ := cf.GetBool(); isStructured {
		labels := make(map[string]string)
		if err := cf.FuzzMap(&labels); err != nil {
			return err
		}

		patch := map[string]interface{}{
			"metadata": map[string]interface{}{"labels": labels},
		}

		raw, err := json.Marshal(patch)
		if err != nil {
			return err
		}
		subset.Patch.Raw = raw
		return nil
	}

	raw, err := cf.GetBytes()
	if err != nil {
		return err
	}
	subset.Patch.Raw = raw
	return nil
}
