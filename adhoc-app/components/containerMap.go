package components

func getContainerIdxMap() map[string]int {
	containerIdxMap := map[string]int{
		"kube-apiserver":          0,
		"etcd":                    1,
		"prometheus":              2,
		"olm-operator":            3,
		"openshift-apiserver":     4,
		"kube-multus":             5,
		"authentication-operator": 6,
		"registry-server":         7,
		"oauth-apiserver":         8,
		"etcd-health-monitor":     9,
	}
	return containerIdxMap
}
