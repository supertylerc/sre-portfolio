package lib

import (
	"fmt"
)

func commonRunCmds() []string {
	return []string{
		"hostnamectl set-hostname $(uuidgen)",
		"modprobe br_netfilter",
		"modprobe overlay",
		"apt-mark hold kubelet kubeadm kubectl",
		"sysctl --system",
		"systemctl restart containerd",
		"systemctl enable containerd",
	}
}
func nodeRunCmds(vars map[string]string) []string {
	return []string{
		fmt.Sprintf("until nc -zvw5 %s 6443; do sleep 0.5; done", vars["control_plane_ip"]),
		fmt.Sprintf("kubeadm join %s:6443 --token %s --discovery-token-unsafe-skip-ca-verification", vars["control_plane_ip"], vars["join_token"]),
	}
}

func controlPlaneRunCmds(vars map[string]string) []string {
	return []string{
		"snap install helm --classic",
		"helm repo add cilium https://helm.cilium.io/",
		"helm repo add argo https://argoproj.github.io/argo-helm",
		"helm repo update",
		fmt.Sprintf("kubeadm init --token=%s --skip-phases=addon/kube-proxy --cri-socket unix:///var/run/containerd/containerd.sock --v=5", vars["join_token"]),
		"curl -Lsk -o /tmp/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64",
		"install /tmp/argocd /usr/bin",
		"helm upgrade --wait --kubeconfig /etc/kubernetes/admin.conf --install --namespace kube-system cilium cilium/cilium -f /tmp/values-cilium.yaml --set k8sServiceHost=$(ip --json -4 a | jq -r '.[] | select(.ifname!=\"lo\") | .addr_info[0].local')",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace argocd",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace gateway",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace cert-manager",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace external-dns",
		"helm upgrade --wait --kubeconfig /etc/kubernetes/admin.conf --install argocd argo/argo-cd -f /tmp/values-argocd.yaml -n argocd",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/gateway-api.yaml",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/argocd.yaml",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/cilium.yaml",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/cert-manager.yaml",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create -f https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/external-dns.yaml",
		fmt.Sprintf("kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=%s --from-literal=email=%s -n cert-manager", vars["cloudflare_token"], vars["cloudflare_email"]),
		fmt.Sprintf("kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=%s --from-literal=email=%s -n external-dns", vars["cloudflare_token"], vars["cloudflare_email"]),
		"mkdir -p /home/supertylerc/.kube",
		"cp -i /etc/kubernetes/admin.conf /home/supertylerc/.kube/config",
		"chown supertylerc:supertylerc /home/supertylerc/.kube",
		"chown supertylerc:supertylerc /home/supertylerc/.kube/config",
		"while kubectl --kubeconfig /etc/kubernetes/admin.conf get application -A | grep -v 'Synced.*Healthy' | grep -v NAME; do sleep 0.5; done",
		"while kubectl --kubeconfig /etc/kubernetes/admin.conf get -A pods -o custom-columns=NAMESPACE:metadata.namespace,POD:metadata.name,PodIP:status.podIP,READY-true:status.containerStatuses[*].ready | grep -v true; do sleep 0.5; done",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf rollout restart deploy/cilium-operator -n kube-system",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf rollout restart ds/cilium -n kube-system",
		"while kubectl --kubeconfig /etc/kubernetes/admin.conf get -A pods -o custom-columns=NAMESPACE:metadata.namespace,POD:metadata.name,PodIP:status.podIP,READY-true:status.containerStatuses[*].ready | grep -v true; do sleep 0.5; done",
	}
}

func CloudRunCmds(nodeKind string, vars map[string]string) []string {
	switch nodeKind {
	case "control-plane":
		return append(commonRunCmds(), controlPlaneRunCmds(vars)...)
	case "node":
		return append(commonRunCmds(), nodeRunCmds(vars)...)
	default:
		return []string{}
	}
}

type Node struct {
	Kind   string
	Num    int
	CPU    int `tf:"cpu"`
	Memory int
	Disk   int
}

func ExpandNodes(nodes []Node) map[string]Node {
	nodeMap := make(map[string]Node)
	for _, n := range nodes {
		for i := 0; i < n.Num; i++ {
			k := fmt.Sprintf("%s-%d", n.Kind, i)
			nodeMap[k] = Node{
				Kind:   n.Kind,
				CPU:    n.CPU,
				Memory: n.Memory,
				Disk:   n.Disk,
				Num:    i,
			}
		}
	}
	return nodeMap
}
