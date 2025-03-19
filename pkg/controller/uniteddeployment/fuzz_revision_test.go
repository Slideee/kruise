package uniteddeployment

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	apps "k8s.io/api/apps/v1"
	"testing"
)

func FuzzGetUnitedDeploymentPatch(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		ud := &appsv1alpha1.UnitedDeployment{}
		if err := f.GenerateStruct(ud); err != nil {
			return
		}
		_, _ = getUnitedDeploymentPatch(ud)
	})
}

func FuzzNextRevision(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		nums, err := f.GetInt()
		if err != nil {
			return
		}

		revisions := make([]*apps.ControllerRevision, 0, nums%30)
		for i := 0; i < nums%30; i++ {
			cr := &apps.ControllerRevision{}
			if err := f.GenerateStruct(cr); err != nil {
				return
			}
			revisions = append(revisions, cr)
		}
		_ = nextRevision(revisions)
	})
}
