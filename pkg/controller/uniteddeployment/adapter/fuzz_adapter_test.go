package adapter

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	"testing"
)

func FuzzTestApplySubsetTemplate(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ud := &appsv1alpha1.UnitedDeployment{}
		if err := f.GenerateStruct(ud); err != nil {
			return
		}

		request, err := f.GetInt()
		if err != nil {
			return
		}

		subsetName, err := f.GetString()
		if err != nil {
			return
		}

		revision, err := f.GetString()
		if err != nil {
			return
		}

		replicas, err := f.GetInt()
		if err != nil {
			return
		}

		partition, err := f.GetInt()
		if err != nil {
			return
		}

		switch request % 4 {
		case 0:
			adapter := &CloneSetAdapter{}
			_ = adapter.ApplySubsetTemplate(ud, subsetName, revision, int32(replicas), int32(partition), adapter.NewResourceObject())
		case 1:
			adapter := &StatefulSetAdapter{}
			_ = adapter.ApplySubsetTemplate(ud, subsetName, revision, int32(replicas), int32(partition), adapter.NewResourceObject())
		case 2:
			adapter := &AdvancedStatefulSetAdapter{}
			_ = adapter.ApplySubsetTemplate(ud, subsetName, revision, int32(replicas), int32(partition), adapter.NewResourceObject())
		case 3:
			adapter := &DeploymentAdapter{}
			_ = adapter.ApplySubsetTemplate(ud, subsetName, revision, int32(replicas), int32(partition), adapter.NewResourceObject())
		}
	})
}
