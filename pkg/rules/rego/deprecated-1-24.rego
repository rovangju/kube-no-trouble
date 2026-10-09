package deprecated124

# Source: https://github.com/kubernetes/kubernetes/blob/v1.24.0/staging/src/k8s.io/api/core/v1/types.go
main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated API fields since 1.24",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "Service"
	r.apiVersion == "v1"
	spec := object.get(r, "spec", {})
	object.get(spec, "type", "ClusterIP") == "LoadBalancer"
	object.get(spec, "loadBalancerIP", "") != ""
	api := {
		"old": "v1",
		"new": "implementation-specific LoadBalancer annotations",
		"since": "1.24",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
