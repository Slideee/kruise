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
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
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
		f := fuzz.NewConsumer(data)

		pod := &corev1.Pod{}
		if err := f.GenerateStruct(pod); err != nil {
			return
		}
		pod.ObjectMeta.SetDeletionTimestamp(nil)

		ws := &appsv1alpha1.WorkloadSpread{}
		if err := f.GenerateStruct(ws); err != nil {
			return
		}
		ws.ObjectMeta.SetDeletionTimestamp(nil)

		favoriteSubset := &appsv1alpha1.WorkloadSpreadSubset{}
		if err := f.GenerateStruct(favoriteSubset); err != nil {
			return
		}

		cl := fake.NewClientBuilder().WithScheme(fakeScheme).WithObjects(
			pod,
			ws,
		).Build()

		r := &ReconcileWorkloadSpread{
			Client: cl,
		}

		_ = r.patchFavoriteSubsetMetadataToPod(pod, ws, favoriteSubset)
	})
}
