package adapter

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	"testing"
)

func FuzzGetSubsetPrefix(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		f := fuzz.NewConsumer(data)

		controllerName, err := f.GetString()
		if err != nil {
			return
		}

		subsetName, err := f.GetString()
		if err != nil {
			return
		}

		_ = getSubsetPrefix(controllerName, subsetName)
	})
}
