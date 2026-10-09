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
			"metadata":   map[string]interface{}{"name": "regular-service"},
			"spec":       map[string]interface{}{"ports": []interface{}{}},
		},
	}

	results, err := regoJudge.Eval(resources)
	if err != nil {
		t.Fatalf("failed to evaluate resources: %s", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 deprecated resources, got %d", len(results))
	}

	for _, result := range results {
		switch result.Name {
		case "example":
			if result.Since.String() != "1.33.0" || result.ReplaceWith != "discovery.k8s.io/v1 EndpointSlice" {
				t.Errorf("unexpected Endpoints finding: %+v", result)
			}
		case "external-service":
			if result.Since.String() != "1.36.0" || result.ReplaceWith != "Use a LoadBalancer or NodePort Service, or Gateway API, instead of spec.externalIPs" {
				t.Errorf("unexpected Service finding: %+v", result)
			}
		default:
			t.Errorf("unexpected deprecated resource: %s", result.Name)
		}
	}

	for _, target := range []struct {
		version  string
		expected int
	}{
		{"1.32", 0},
		{"1.35", 1},
		{"1.36", 2},
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
