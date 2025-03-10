package workloadspread

import (
	"encoding/json"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func FuzzInjectWorkloadSpreadIntoPodWithHeaders(f *testing.F) {
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

func FuzzInjectWorkloadSpreadIntoPodWithNative(f *testing.F) {
	f.Add([]byte(`{"metadata":{"labels":{"subset":"subset-a"},"annotations":{"subset":"subset-a"}}}`), "subset-a", "uid")

	f.Fuzz(func(t *testing.T, data []byte, subsetName, generatedUID string) {
		pod := podDemo.DeepCopy()
		ws := workloadSpreadDemo.DeepCopy()

		ws.Spec.Subsets = []appsv1alpha1.WorkloadSpreadSubset{
			{
				Name: subsetName,
				Patch: runtime.RawExtension{
					Raw: data,
				},
			},
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

func FuzzGetReplicasFromObjectWithHeaders(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		replicasPath, err := f.GetString()
		if err != nil {
			return
		}

		var obj unstructured.Unstructured
		if err := f.GenerateStruct(&obj); err != nil {
			return
		}

		_, _ = GetReplicasFromObject(&obj, replicasPath)
	})
}

func FuzzGetReplicasFromObjectWithNative(f *testing.F) {
	initial, _ := json.Marshal(&unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"replicas":     int64(5),
				"replicaSlice": []int64{1, 2},
				"stringField":  "5",
			},
		},
	})

	testCases := []struct {
		data []byte
		path string
	}{
		{initial, "spec.replicas"},
		{initial, "spec.replicaSlice.1"},
		{initial, "spec.stringField"},
		{nil, ""},
	}

	for _, tc := range testCases {
		f.Add(tc.data, tc.path)
	}

	f.Fuzz(func(t *testing.T, data []byte, path string) {
		var object unstructured.Unstructured

		err := json.Unmarshal(data, &object.Object)
		if err != nil {
			return
		}

		_, _ = GetReplicasFromObject(&object, path)
	})
}
