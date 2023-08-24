package controllers

import (
	"encoding/json"
	sriovnetworkv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/apimachinery/pkg/util/runtime"

	dptypes "github.com/k8snetworkplumbingwg/sriov-network-device-plugin/pkg/types"
)

func TestNodeSelectorMerge(t *testing.T) {
	table := []struct {
		tname    string
		policies []sriovnetworkv1.SriovNetworkNodePolicy
		expected []corev1.NodeSelectorTerm
	}{
		{
			tname: "testoneselector",
			policies: []sriovnetworkv1.SriovNetworkNodePolicy{
				{
					Spec: sriovnetworkv1.SriovNetworkNodePolicySpec{
						NodeSelector: map[string]string{
							"foo": "bar",
						},
					},
				},
				{
					Spec: sriovnetworkv1.SriovNetworkNodePolicySpec{
						NodeSelector: map[string]string{
							"bb": "cc",
						},
					},
				},
			},
			expected: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "foo",
							Values:   []string{"bar"},
						},
					},
				},
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "bb",
							Values:   []string{"cc"},
						},
					},
				},
			},
		},
		{
			tname: "testtwoselectors",
			policies: []sriovnetworkv1.SriovNetworkNodePolicy{
				{
					Spec: sriovnetworkv1.SriovNetworkNodePolicySpec{
						NodeSelector: map[string]string{
							"foo":  "bar",
							"foo1": "bar1",
						},
					},
				},
				{
					Spec: sriovnetworkv1.SriovNetworkNodePolicySpec{
						NodeSelector: map[string]string{
							"bb":  "cc",
							"bb1": "cc1",
							"bb2": "cc2",
						},
					},
				},
			},
			expected: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "foo",
							Values:   []string{"bar"},
						},
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "foo1",
							Values:   []string{"bar1"},
						},
					},
				},
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "bb",
							Values:   []string{"cc"},
						},
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "bb1",
							Values:   []string{"cc1"},
						},
						{
							Operator: corev1.NodeSelectorOpIn,
							Key:      "bb2",
							Values:   []string{"cc2"},
						},
					},
				},
			},
		},
	}

	for _, tc := range table {
		t.Run(tc.tname, func(t *testing.T) {
			selectors := nodeSelectorTermsForPolicyList(tc.policies)
			if !cmp.Equal(selectors, tc.expected) {
				t.Error(tc.tname, "Selectors not as expected", cmp.Diff(selectors, tc.expected))
			}
		})
	}
}

func mustMarshallSelector(t *testing.T, input *dptypes.NetDeviceSelectors) *json.RawMessage {
	out, err := json.Marshal(input)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	ret := json.RawMessage(out)
	return &ret
}
