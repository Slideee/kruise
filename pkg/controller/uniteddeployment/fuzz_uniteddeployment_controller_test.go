package uniteddeployment

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	"testing"
)

func FuzzGetNextUpdate(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ud := &appsv1alpha1.UnitedDeployment{}
		if err := f.GenerateStruct(ud); err != nil {
			return
		}

		nextReplicas := make(map[string]int32)
		if err := f.FuzzMap(&nextReplicas); err != nil {
			return
		}

		nextPartitions := make(map[string]int32)
		if err := f.FuzzMap(&nextPartitions); err != nil {
			return
		}

		_ = getNextUpdate(ud, &nextReplicas, &nextPartitions)
	})
}

func FuzzCalcNextPartitions(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ud := &appsv1alpha1.UnitedDeployment{}
		if err := f.GenerateStruct(ud); err != nil {
			return
		}

		nextReplicas := make(map[string]int32)
		if err := f.FuzzMap(&nextReplicas); err != nil {
			return
		}

		_ = calcNextPartitions(ud, &nextReplicas)
	})
}
