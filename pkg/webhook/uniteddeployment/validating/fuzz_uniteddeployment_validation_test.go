package validating

import (
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	fuzzutils "github.com/openkruise/kruise/pkg/util/fuzz"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func FuzzValidateUnitedDeploymentSpec(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		cf := fuzz.NewConsumer(data)
		ud := &appsv1alpha1.UnitedDeployment{}
		if err := cf.GenerateStruct(ud); err != nil {
			return
		}

		if err := fuzzutils.GenerateReplicas(cf, ud); err != nil {
			return
		}

		if err := fuzzutils.GenerateSelector(cf, ud); err != nil {
			return
		}

		if err := fuzzutils.GenerateTemplate(cf, ud); err != nil {
			return
		}

		if err := fuzzutils.GenerateSubset(cf, ud); err != nil {
			return
		}

		if err := fuzzutils.GenerateUpdateStrategy(cf, ud); err != nil {
			return
		}

		_ = validateUnitedDeploymentSpec(&ud.Spec, field.NewPath("spec"))
	})
}
