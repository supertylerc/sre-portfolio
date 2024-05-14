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
func nodeRunCmds() []string {
	return []string{
		"until nc -zvw5 ${cidrhost(var.libvirt_cidr, var.control_plane_num)} 6443; do sleep 0.5; done",
		"kubeadm join ${cidrhost(var.libvirt_cidr, var.control_plane_num)}:6443 --token ${var.join_token} --discovery-token-unsafe-skip-ca-verification",
	}
}

func controlPlaneRunCmds() []string {
	return []string{
		"snap install helm --classic",
		"helm repo add cilium https://helm.cilium.io/",
		"helm repo add argo https://argoproj.github.io/argo-helm",
		"helm repo update",
		"kubeadm init --token=${var.join_token} --skip-phases=addon/kube-proxy --cri-socket unix:///var/run/containerd/containerd.sock --v=5",
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
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=${var.cloudflare_token} --from-literal=email=${var.cloudflare_email} -n cert-manager",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=${var.cloudflare_token} --from-literal=email=${var.cloudflare_email} -n external-dns",
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

func CloudRunCmds(nodeKind string) []string {
	switch nodeKind {
	case "control-plane":
		return append(commonRunCmds(), controlPlaneRunCmds()...)
	case "node":
		return append(commonRunCmds(), nodeRunCmds()...)
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
