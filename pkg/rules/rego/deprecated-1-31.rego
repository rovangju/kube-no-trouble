package deprecated131

# Source: https://github.com/kubernetes/kubernetes/blob/v1.33.0/CHANGELOG/CHANGELOG-1.33.md
main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated API fields since 1.31",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	r.kind == "Node"
	r.apiVersion == "v1"
	status := object.get(r, "status", {})
	node_info := object.get(status, "nodeInfo", {})
	object.get(node_info, "kubeProxyVersion", null) != null
	api := {
		"old": "v1",
		"new": "kube-proxy version is no longer reported in Node status",
		"since": "1.31",
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
