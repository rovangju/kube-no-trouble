package deprecated135

# Source: https://github.com/kubernetes/kubernetes/blob/v1.35.0/CHANGELOG/CHANGELOG-1.35.md
main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated API values since 1.35",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "Service"
	r.apiVersion == "v1"
	spec := object.get(r, "spec", {})
	object.get(spec, "trafficDistribution", "") == "PreferClose"
	api := {
		"old": "v1",
		"new": "spec.trafficDistribution: PreferSameZone",
		"since": "1.35",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
