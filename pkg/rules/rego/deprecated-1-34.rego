package deprecated134

# Source: https://github.com/kubernetes/kubernetes/blob/v1.37.0/staging/src/k8s.io/api/certificates/v1alpha1/zz_generated.prerelease-lifecycle.go
main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated APIs since 1.34",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "ClusterTrustBundle"
	r.apiVersion == "certificates.k8s.io/v1alpha1"
	api := {
		"old": "certificates.k8s.io/v1alpha1",
		"new": "certificates.k8s.io/v1beta1 (v1 is available in Kubernetes 1.37+)",
		"since": "1.34",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
