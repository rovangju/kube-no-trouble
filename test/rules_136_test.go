package test

import (
	"testing"

	"github.com/doitintl/kube-no-trouble/pkg/judge"
	"github.com/doitintl/kube-no-trouble/pkg/printer"
	"github.com/doitintl/kube-no-trouble/pkg/rules"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestRego136(t *testing.T) {
	testCases := []resourceFixtureTestCase{
		{"Service with deprecated externalIPs", []string{"../fixtures/service-external-ips-v1.yaml"}, []string{"Service"}},
	}

	testResourcesUsingFixtures(t, testCases)
}

func TestRegoLatestDeprecations(t *testing.T) {
	loadedRules, err := rules.FetchRegoRules([]schema.GroupVersionKind{})
	if err != nil {
		t.Fatalf("failed to load rules: %s", err)
	}

	regoJudge, err := judge.NewRegoJudge(&judge.RegoOpts{}, loadedRules)
	if err != nil {
		t.Fatalf("failed to create judge: %s", err)
	}

	resources := []map[string]interface{}{
		{
			"apiVersion": "v1",
			"kind":       "Node",
			"metadata":   map[string]interface{}{"name": "node"},
			"status": map[string]interface{}{
				"nodeInfo": map[string]interface{}{"kubeProxyVersion": "v1.30.0"},
			},
		},
		{
			"apiVersion": "v1",
			"kind":       "Endpoints",
			"metadata":   map[string]interface{}{"name": "example"},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "external-service"},
			"spec":       map[string]interface{}{"externalIPs": []interface{}{"192.0.2.10"}},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":        "topology-service",
				"annotations": map[string]interface{}{"service.kubernetes.io/topology-mode": "Auto"},
			},
			"spec": map[string]interface{}{},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "prefer-close-service"},
			"spec":       map[string]interface{}{"trafficDistribution": "PreferClose"},
		},
		{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   map[string]interface{}{"name": "resizing-pod"},
			"status":     map[string]interface{}{"resize": "Proposed"},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "empty-external-ips-service"},
			"spec":       map[string]interface{}{"externalIPs": []interface{}{}},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "service-without-spec"},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "regular-service"},
			"spec":       map[string]interface{}{"ports": []interface{}{}},
		},
	}

	results, err := regoJudge.Eval(resources)
	if err != nil {
		t.Fatalf("failed to evaluate resources: %s", err)
	}
	if len(results) != 6 {
		t.Fatalf("expected 6 deprecated resources, got %d", len(results))
	}

	expected := map[string]struct {
		since       string
		replaceWith string
	}{
		"node":                 {"1.31.0", "kube-proxy version is no longer reported in Node status"},
		"example":              {"1.33.0", "discovery.k8s.io/v1 EndpointSlice"},
		"topology-service":     {"1.33.0", "spec.trafficDistribution"},
		"resizing-pod":         {"1.33.0", "status.conditions (PodResizePending and PodResizeInProgress)"},
		"prefer-close-service": {"1.35.0", "spec.trafficDistribution: PreferSameZone"},
		"external-service":     {"1.36.0", "Use a LoadBalancer or NodePort Service, or Gateway API, instead of spec.externalIPs"},
	}
	for _, result := range results {
		want, ok := expected[result.Name]
		if !ok {
			t.Errorf("unexpected deprecated resource: %s", result.Name)
			continue
		}
		if result.Since.String() != want.since || result.ReplaceWith != want.replaceWith {
			t.Errorf("unexpected finding for %s: %+v", result.Name, result)
		}
	}

	for _, target := range []struct {
		version  string
		expected int
	}{
		{"1.30", 0},
		{"1.31", 1},
		{"1.32", 1},
		{"1.33", 4},
		{"1.34", 4},
		{"1.35", 5},
		{"1.36", 6},
	} {
		targetVersion, err := judge.NewVersion(target.version)
		if err != nil {
			t.Fatalf("failed to parse target version: %s", err)
		}
		filtered, err := printer.FilterNonRelevantResults(results, targetVersion)
		if err != nil {
			t.Fatalf("failed to filter results: %s", err)
		}
		if len(filtered) != target.expected {
			t.Errorf("target %s: expected %d findings, got %d", target.version, target.expected, len(filtered))
		}
	}
}
