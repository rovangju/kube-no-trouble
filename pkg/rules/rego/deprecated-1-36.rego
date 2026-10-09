package deprecated136

# Source: https://kubernetes.io/blog/2026/05/14/kubernetes-v1-36-deprecation-and-removal-of-service-externalips/
main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		# Namespace does not have to be defined in case of local manifests
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated API fields since 1.36",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "Service"
	r.apiVersion == "v1"
	spec := object.get(r, "spec", {})
	count(object.get(spec, "externalIPs", [])) > 0
	api := {
		"old": "v1",
		"new": "Use a LoadBalancer or NodePort Service, or Gateway API, instead of spec.externalIPs",
		"since": "1.36",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
