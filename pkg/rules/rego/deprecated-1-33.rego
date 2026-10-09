package deprecated133

# Source: https://kubernetes.io/blog/2025/04/24/endpoints-deprecation/
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
		"RuleSet": "Deprecated APIs since 1.33",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "Endpoints"
	r.apiVersion == "v1"
	api := {
		"old": "v1",
		"new": "discovery.k8s.io/v1 EndpointSlice",
		"since": "1.33",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
