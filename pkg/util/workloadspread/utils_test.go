package workloadspread

import (
	"testing"

	appsv1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsPodSelected(t *testing.T) {
	commonFilter := &appsv1alpha1.TargetFilter{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "selected",
			},
		},
	}
	cases := []struct {
		name     string
		filter   *appsv1alpha1.TargetFilter
		labels   map[string]string
		selected bool
		wantErr  bool
	}{
		{
			name:   "selected",
			filter: commonFilter,
			labels: map[string]string{
				"app": "selected",
			},
			selected: true,
		},
		{
			name:   "not selected",
			filter: commonFilter,
			labels: map[string]string{
				"app": "not-selected",
			},
			selected: false,
		},
		{
			name:   "selector is nil",
			filter: nil,
			labels: map[string]string{
				"app": "selected",
			},
			selected: true,
		},
		{
			name: "selector is invalid",
			filter: &appsv1alpha1.TargetFilter{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "selected",
					},
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      "app",
							Operator: "Invalid",
							Values:   []string{"selected"},
						},
					},
				},
			},
			selected: false,
			wantErr:  true,
		},
	}
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			selected, err := IsPodSelected(cs.filter, cs.labels)
			if selected != cs.selected || (err != nil) != cs.wantErr {
				t.Fatalf("got unexpected result, actual: [selected=%v,err=%v] expected: [selected=%v,wantErr=%v]",
					selected, err, cs.selected, cs.wantErr)
			}
		})
	}
}
