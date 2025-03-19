package uniteddeployment

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"k8s.io/apimachinery/pkg/util/intstr"

	"testing"
)

func FuzzParseSubsetReplicas(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		udReplicas, err := f.GetInt()
		if err != nil {
			return
		}

		subsetReplicas := &intstr.IntOrString{}
		if err = f.GenerateStruct(subsetReplicas); err != nil {
			return
		}

		_, _ = ParseSubsetReplicas(int32(udReplicas), *subsetReplicas)
	})
}
