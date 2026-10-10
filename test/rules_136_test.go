package test

import (
	"testing"

	"github.com/rovangju/kube-no-trouble/pkg/judge"
	"github.com/rovangju/kube-no-trouble/pkg/printer"
	"github.com/rovangju/kube-no-trouble/pkg/rules"
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
			"metadata":   map[string]interface{}{"name": "load-balancer-service"},
			"spec": map[string]interface{}{
				"type":           "LoadBalancer",
				"loadBalancerIP": "192.0.2.20",
			},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "ignored-load-balancer-ip"},
			"spec":       map[string]interface{}{"loadBalancerIP": "192.0.2.21"},
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
			"metadata":   map[string]interface{}{"name": "headless-external-service"},
			"spec": map[string]interface{}{
				"clusterIP":   "None",
				"externalIPs": []interface{}{"192.0.2.11"},
			},
		},
		{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]interface{}{"name": "external-name-service"},
			"spec": map[string]interface{}{
				"type":        "ExternalName",
				"externalIPs": []interface{}{"192.0.2.12"},
			},
		},
		{
			"apiVersion": "certificates.k8s.io/v1alpha1",
			"kind":       "ClusterTrustBundle",
			"metadata":   map[string]interface{}{"name": "old-trust-bundle"},
		},
		{
			"apiVersion": "certificates.k8s.io/v1beta1",
			"kind":       "ClusterTrustBundle",
			"metadata":   map[string]interface{}{"name": "beta-trust-bundle"},
		},
		{
			"apiVersion": "certificates.k8s.io/v1beta1",
			"kind":       "PodCertificateRequest",
			"metadata":   map[string]interface{}{"name": "pod-certificate-request"},
			"spec": map[string]interface{}{
				"pkixPublicKey":     "cHVibGljLWtleQ==",
				"proofOfPossession": "cHJvb2Y=",
			},
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
	if len(results) != 10 {
		t.Fatalf("expected 10 deprecated resources or fields, got %d", len(results))
	}

	expected := map[string]struct {
		since       string
		replaceWith string
	}{
		"load-balancer-service/Deprecated API fields since 1.24":   {"1.24.0", "implementation-specific LoadBalancer annotations"},
		"node/Deprecated API fields since 1.31":                    {"1.31.0", "kube-proxy version is no longer reported in Node status"},
		"example/Deprecated APIs since 1.33":                       {"1.33.0", "discovery.k8s.io/v1 EndpointSlice"},
		"topology-service/Deprecated APIs since 1.33":              {"1.33.0", "spec.trafficDistribution"},
		"old-trust-bundle/Deprecated APIs since 1.34":              {"1.34.0", "certificates.k8s.io/v1beta1 (v1 is available in Kubernetes 1.37+)"},
		"prefer-close-service/Deprecated API values since 1.35":    {"1.35.0", "spec.trafficDistribution: PreferSameZone"},
		"external-service/Deprecated API fields since 1.36":        {"1.36.0", "Use an administrator-managed LoadBalancer Service or Gateway API instead of spec.externalIPs"},
		"pod-certificate-request/Deprecated API fields since 1.36": {"1.36.0", "spec.stubPKCS10Request"},
		"beta-trust-bundle/Deprecated APIs since 1.37":             {"1.37.0", "certificates.k8s.io/v1"},
		"pod-certificate-request/Deprecated APIs since 1.37":       {"1.37.0", "certificates.k8s.io/v1"},
	}
	for _, result := range results {
		key := result.Name + "/" + result.RuleSet
		want, ok := expected[key]
		if !ok {
			t.Errorf("unexpected deprecated resource or field: %s", key)
			continue
		}
		if result.Since.String() != want.since || result.ReplaceWith != want.replaceWith {
			t.Errorf("unexpected finding for %s: %+v", key, result)
		}
	}

	for _, target := range []struct {
		version  string
		expected int
	}{
		{"1.23", 0},
		{"1.24", 1},
		{"1.30", 1},
		{"1.31", 2},
		{"1.32", 2},
		{"1.33", 4},
		{"1.34", 5},
		{"1.35", 6},
		{"1.36", 8},
		{"1.37", 10},
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
