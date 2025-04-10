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

package fuzz

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	StructuredResources = []struct {
		Name string
		Data []byte
	}{
		{
			Name: "Secret",
			Data: []byte(`{
				"apiVersion": "v1",
				"data": {
					"test": "MWYyZDFlMmU2N2Rm"
				},
				"kind": "Secret",
				"metadata": {
					"name": "test-secret-2"
				},
				"type": "Opaque"
			}`),
		},
		{
			Name: "ConfigMap",
			Data: []byte(`{
				"apiVersion": "v1",
				"data": {
					"game.properties": "enemy.types=aliens,monsters\nplayer.maximum-lives=5\n",
					"player_initial_lives": "3",
					"ui_properties_file_name": "user-interface.properties",
					"user-interface.properties": "color.good=purple\ncolor.bad=yellow\nallow.textmode=true\n"
				},
				"kind": "ConfigMap",
				"metadata": {
					"name": "game-demo"
				}
			}`),
		},
		{
			Name: "Pod",
			Data: []byte(`{
				"apiVersion": "v1",
				"kind": "Pod",
				"metadata": {
					"name": "test-pod-1"
				},
				"spec": {
					"containers": [
						{
							"image": "nginx:1.14.2",
							"name": "test-container"
						}
					]
				}
			}`),
		},
	}
)

func GenerateResourceDistributionResource(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.ResourceDistribution) error {
	isStructured, err := cf.GetBool()
	if err != nil {
		return err
	}

	if isStructured {
		choice, err := cf.GetInt()
		if err != nil {
			return err
		}
		ud.Spec.Resource = runtime.RawExtension{
			Raw: StructuredResources[choice%len(StructuredResources)].Data,
		}
		return nil
	}

	raw := runtime.RawExtension{}
	if err := cf.GenerateStruct(&raw); err != nil {
		return err
	}
	ud.Spec.Resource = raw
	return nil
}

func GenerateResourceDistributionTargets(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.ResourceDistribution) error {
	isStructured, err := cf.GetBool()
	if err != nil {
		return err
	}

	if isStructured {
		targets := appsv1alpha1.ResourceDistributionTargets{}

		includedNamespacesSlice := make([]appsv1alpha1.ResourceDistributionNamespace, 0)
		if err := cf.CreateSlice(&includedNamespacesSlice); err != nil {
			return err
		}

		excludedNamespacesSlice := make([]appsv1alpha1.ResourceDistributionNamespace, 0)
		if err := cf.CreateSlice(&excludedNamespacesSlice); err != nil {
			return err
		}
		targets.IncludedNamespaces.List = includedNamespacesSlice
		targets.ExcludedNamespaces.List = excludedNamespacesSlice

		for i := range targets.IncludedNamespaces.List {
			if valid, err := cf.GetBool(); valid && err == nil {
				targets.IncludedNamespaces.List[i].Name = GenerateValidNamespaceName(cf)
			} else {
				targets.IncludedNamespaces.List[i].Name = GenerateInvalidNamespaceName(cf)
			}
		}

		for i := range targets.ExcludedNamespaces.List {
			if valid, err := cf.GetBool(); valid && err == nil {
				targets.ExcludedNamespaces.List[i].Name = GenerateValidNamespaceName(cf)
			} else {
				targets.ExcludedNamespaces.List[i].Name = GenerateInvalidNamespaceName(cf)
			}
		}

		labels := make(map[string]string)
		if err := cf.FuzzMap(&labels); err != nil {
			return err
		}

		labelSelectorSlice := make([]metav1.LabelSelectorRequirement, 0)
		if err := cf.CreateSlice(&labelSelectorSlice); err != nil {
			return err
		}

		targets.NamespaceLabelSelector.MatchLabels = labels
		targets.NamespaceLabelSelector.MatchExpressions = labelSelectorSlice
		ud.Spec.Targets = targets
		return nil
	}

	targets := appsv1alpha1.ResourceDistributionTargets{}
	if err := cf.GenerateStruct(&targets); err != nil {
		return err
	}
	ud.Spec.Targets = targets
	return nil
}

func GenerateResourceObject(cf *fuzz.ConsumeFuzzer) (runtime.Object, error) {
	choice, err := cf.GetInt()
	if err != nil {
		return nil, err
	}

	var obj runtime.Object
	switch choice % 2 {
	case 0:
		obj = &corev1.ConfigMap{}
	case 1:
		obj = &corev1.Secret{}
	}

	if err := cf.GenerateStruct(obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func GenerateValidNamespaceName(cf *fuzz.ConsumeFuzzer) string {
	base := "test-"
	// Generate a valid namespace name (DNS-1123 compliant)
	if name, err := cf.GetStringFrom(base, 63-len(base)); err == nil {
		return base + name
	}
	return base
}

func GenerateInvalidNamespaceName(cf *fuzz.ConsumeFuzzer) string {
	invalidChars := []rune{'$', '_', ' ', '💣'}
	name, err := cf.GetString()
	if err != nil || name == "" {
		return "_invalid"
	}

	runes := []rune(name)
	choice, err := cf.GetInt()
	if err != nil {
		return "_invalid"
	}

	// Make sure the first character is illegal
	runes[0] = invalidChars[choice%(len(invalidChars))]

	if len(runes) > 253 {
		return string(runes[:253])
	}
	return string(runes)
}
