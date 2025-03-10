package workloadspread

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

// generateRawSeedInputs generates a set of raw JSON byte slices representing different seed inputs.
// These seeds are used to simulate various configurations for fuzz testing.
func generateRawSeedInputs() [][]byte {
	return [][]byte{
		[]byte(``),
		[]byte(`{}`),
		[]byte(`{"metadata":{"labels":{"topology.application.deploy/zone":"true"},"annotations":{"topology.application.deploy/zone":"true"}}}`),
		[]byte(`{"spec":{"containers":[{"name":"main","resources":{"limits":{"cpu":"2","memory":"800Mi"}}}]}}`),
		[]byte(`{"spec":{"containers":[{"name":"main","command":["echo","hello"]}]}}`),
		[]byte(`{"spec":{"containers":[{"name":"main","env":[{"name":"K8S_AZ_NAME","value":"zone-a"}]}]}}`),
		[]byte(`{"spec":{"tolerations":[{"key":"key","operator":"Equal","value":"value","effect":"NoSchedule"}]}}`),
		[]byte(`{"spec":{"nodeSelector":{"topology.application.deploy/zone":"zone-a"}}}`),
		[]byte(`{"spec":{"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"topology.application.deploy/zone","operator":"In","values":["zone-a"]}]}]}}}}}`),
		[]byte(`{"spec":{"volumes":[{"name":"config","configMap":{"name":"config"}}]}}`),
	}
}

func FuzzMatchesSubsetWithHeaders(f *testing.F) {
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

func FuzzMatchesSubsetWithNative(f *testing.F) {
	for _, seed := range generateRawSeedInputs() {
		f.Add(seed, 0)
	}

	f.Fuzz(func(t *testing.T, data []byte, missingReplicas int) {
		pod := matchPodDemo.DeepCopy()
		node := matchNodeDemo.DeepCopy()
		subset := matchSubsetDemo.DeepCopy()
		subset.Patch = runtime.RawExtension{Raw: data}
		_, _, _ = matchesSubset(pod, node, subset, missingReplicas)
	})
}

func FuzzPodPreferredScoreWithHeaders(f *testing.F) {
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

func FuzzPodPreferredScoreWithNative(f *testing.F) {
	for _, seed := range generateRawSeedInputs() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		pod := matchPodDemo.DeepCopy()
		subset := matchSubsetDemo.DeepCopy()
		subset.Patch = runtime.RawExtension{Raw: data}
		_ = podPreferredScore(subset, pod)
	})
}

func FuzzMatchesSubsetRequiredAndTolerationWithHeaders(f *testing.F) {
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
