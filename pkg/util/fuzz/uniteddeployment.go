package fuzz

import (
	"encoding/json"
	"fmt"
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type SubsetFunc = func(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error

func GenerateReplicas(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment) error {
	if rep, err := cf.GetInt(); err == nil {
		r := int32(rep)
		ud.Spec.Replicas = &r
	} else {
		r := int32(5)
		ud.Spec.Replicas = &r
	}
	return nil
}

func GenerateSubset(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment, fns ...SubsetFunc) error {
	num, err := cf.GetInt()
	if err != nil {
		return err
	}

	if len(fns) == 0 {
		fns = []SubsetFunc{
			GenerateSubsetName,
			GenerateSubSetReplicas,
			GeneratePatch,
			GenerateNodeSelectorTerm,
			GenerateTolerations,
		}
	}

	nSubsets := (num % 5) + 1
	subsets := make([]appsv1alpha1.Subset, nSubsets)

	for i := 0; i < nSubsets; i++ {
		for _, fn := range fns {
			if err := fn(cf, &subsets[i]); err != nil {
				return err
			}
		}
	}

	ud.Spec.Topology.Subsets = subsets
	return nil
}

func GenerateScheduleStrategy(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment) error {
	choice, err := cf.GetInt()
	if err != nil {
		return err
	}

	switch choice % 3 {
	case 0:
		ud.Spec.Topology.ScheduleStrategy.Type = appsv1alpha1.AdaptiveUnitedDeploymentScheduleStrategyType
	case 1:
		ud.Spec.Topology.ScheduleStrategy.Type = appsv1alpha1.FixedUnitedDeploymentScheduleStrategyType
	case 2:
		str, err := cf.GetString()
		if err != nil {
			return err
		}
		ud.Spec.Topology.ScheduleStrategy.Type = appsv1alpha1.UnitedDeploymentScheduleStrategyType(str)
	}
	return nil
}

func GenerateSelector(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment) error {
	if set, err := cf.GetBool(); set && err == nil {
		selector := &metav1.LabelSelector{}
		if nonEmpty, err := cf.GetBool(); nonEmpty && err == nil {
			labelsMap := make(map[string]string)
			if err := cf.FuzzMap(&labelsMap); err != nil {
				return err
			}
			selector.MatchLabels = labelsMap
		} else {
			selector.MatchLabels = map[string]string{}
			selector.MatchExpressions = []metav1.LabelSelectorRequirement{}
		}
		ud.Spec.Selector = selector
	}
	return nil
}

func GenerateTemplate(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment) error {
	var tmpl appsv1alpha1.SubsetTemplate

	choice, err := cf.GetInt()
	if err != nil {
		return err
	}

	switch choice % 5 {
	case 0:
		s := &appsv1alpha1.DeploymentTemplateSpec{}
		if err = cf.GenerateStruct(s); err != nil {
			return err
		}
		tmpl.DeploymentTemplate = s
	case 1:
		s := &appsv1alpha1.AdvancedStatefulSetTemplateSpec{}
		if err = cf.GenerateStruct(s); err != nil {
			return err
		}
		tmpl.AdvancedStatefulSetTemplate = s
	case 2:
		s := &appsv1alpha1.StatefulSetTemplateSpec{}
		if err = cf.GenerateStruct(s); err != nil {
			return err
		}
		tmpl.StatefulSetTemplate = s
	case 3:
		s := &appsv1alpha1.CloneSetTemplateSpec{}
		if err = cf.GenerateStruct(s); err != nil {
			return err
		}
		tmpl.CloneSetTemplate = s
	case 4:
		if err = cf.GenerateStruct(&tmpl); err != nil {
			return err
		}
	}

	ud.Spec.Template = tmpl
	return nil
}
func GenerateUpdateStrategy(cf *fuzz.ConsumeFuzzer, ud *appsv1alpha1.UnitedDeployment) error {
	setParts, err := cf.GetBool()
	if err != nil || !setParts {
		return err
	}

	np, err := cf.GetInt()
	if err != nil {
		return err
	}
	numParts := np % 3
	partitions := make(map[string]int32)

	for j := 0; j < numParts; j++ {
		var key string

		useExisting, err := cf.GetBool()
		if err == nil && useExisting {
			if len(ud.Spec.Topology.Subsets) > 0 {
				key = ud.Spec.Topology.Subsets[j%len(ud.Spec.Topology.Subsets)].Name
			} else {
				key, err = cf.GetString()
				if err != nil {
					return err
				}
			}
		} else {
			key, err = cf.GetString()
			if err != nil {
				return err
			}
		}
		partitions[key] = int32(j)
	}

	ud.Spec.UpdateStrategy.ManualUpdate = &appsv1alpha1.ManualUpdate{
		Partitions: partitions,
	}
	return nil
}

func GenerateSubsetName(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error {
	name, err := cf.GetString()
	if err != nil {
		return err
	}
	subset.Name = name
	return nil
}

func GenerateSubSetReplicas(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error {
	if setMin, err := cf.GetBool(); setMin && err == nil {
		minVal, err := GenerateSubsetReplicas(cf)
		if err != nil {
			return err
		}
		subset.MinReplicas = &minVal
	}

	if setMax, err := cf.GetBool(); setMax && err == nil {
		maxVal, err := GenerateSubsetReplicas(cf)
		if err != nil {
			return err
		}
		subset.MaxReplicas = &maxVal
	}

	if setReplicas, err := cf.GetBool(); setReplicas && err == nil {
		replicas, err := GenerateSubsetReplicas(cf)
		if err != nil {
			return err
		}
		subset.Replicas = &replicas
	}

	return nil
}

func GeneratePatch(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error {
	isStructured, err := cf.GetBool()
	if err != nil {
		return err
	}

	var raw []byte
	if isStructured {
		labels := make(map[string]string)
		if err := cf.FuzzMap(&labels); err != nil {
			return err
		}
		patch := map[string]interface{}{
			"metadata": map[string]interface{}{"labels": labels},
		}
		raw, err = json.Marshal(patch)
		if err != nil {
			return err
		}
	} else {
		raw, err = cf.GetBytes()
		if err != nil {
			return err
		}
	}
	subset.Patch.Raw = raw
	return nil
}

func GenerateNodeSelectorTerm(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error {
	isStructured, err := cf.GetBool()
	if err != nil {
		return err
	}

	var term corev1.NodeSelectorTerm
	if isStructured {
		term = corev1.NodeSelectorTerm{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      "key",
					Operator: "In",
					Values:   []string{"value"},
				},
			},
		}
	} else {
		if err := cf.GenerateStruct(&term); err != nil {
			return err
		}
	}
	subset.NodeSelectorTerm = term
	return nil
}

func GenerateTolerations(cf *fuzz.ConsumeFuzzer, subset *appsv1alpha1.Subset) error {
	isStructured, err := cf.GetBool()
	if err != nil {
		return err
	}

	var tolerations []corev1.Toleration
	if isStructured {
		tolerations = []corev1.Toleration{
			{
				Key:      "key",
				Operator: "In",
				Value:    "value",
			},
		}
	} else {
		toleration := corev1.Toleration{}
		if err := cf.GenerateStruct(&toleration); err != nil {
			return err
		}
		tolerations = []corev1.Toleration{toleration}
	}
	subset.Tolerations = tolerations
	return nil
}

func GenerateSubsetReplicas(cf *fuzz.ConsumeFuzzer) (intstr.IntOrString, error) {
	// First, decide if we are generating an integer or a string variant.
	isInt, err := cf.GetBool()
	if err != nil {
		return intstr.IntOrString{}, err
	}

	if isInt {
		intVal, err := cf.GetInt()
		if err != nil {
			return intstr.IntOrString{}, err
		}
		return intstr.FromInt32(int32(intVal)), nil
	}

	// For the string variant, decide whether to append a '%' suffix.
	hasSuffix, err := cf.GetBool()
	if err != nil {
		return intstr.IntOrString{}, err
	}

	if hasSuffix {
		percent, err := cf.GetInt()
		if err != nil {
			return intstr.IntOrString{}, err
		}
		// Ensure the percentage is within a reasonable range using modulo.
		return intstr.FromString(fmt.Sprintf("%d%%", percent%1000)), nil
	}

	strVal, err := cf.GetString()
	if err != nil {
		return intstr.IntOrString{}, err
	}
	return intstr.FromString(strVal), nil
}
