package lib

import (
	"fmt"
)

type NodeRunCmdVars struct {
	JoinToken      string `tf:"join_token"`
	ControlPlaneIP string `tf:"control_plane_ip"`
}
type ControlPlaneRunCmdVars struct {
	JoinToken       string   `tf:"join_token"`
	ArgoCDApps      []string `tf:"argocd_apps"`
	CloudflareToken string   `tf:"cloudflare_token"`
	CloudflareEmail string   `tf:"cloudflare_email"`
	CNI             string   `tf:"cni"`
}

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
func NodeRunCmds(vars NodeRunCmdVars) []string {
	return append(commonRunCmds(), []string{
		fmt.Sprintf("until nc -zvw5 %s 6443; do sleep 0.5; done", vars.ControlPlaneIP),
		fmt.Sprintf("kubeadm join %s:6443 --token %s --discovery-token-unsafe-skip-ca-verification", vars.ControlPlaneIP, vars.JoinToken),
	}...)
}

func argoCDApps(apps []string) []string {
	cmds := []string
	for _, v := range apps {
		cmds = append(cmds, fmt.Sprintf("kubectl --kubeconfig /etc/kubernetes/admin.conf create -f %s", v))
	}
	return cmds
}

func ControlPlaneRunCmds(vars ControlPlaneRunCmdVars) []string {
	cmds := append(commonRunCmds(), []string{
		"snap install helm --classic",
		"helm repo add argo https://argoproj.github.io/argo-helm",
		"helm repo update",
		fmt.Sprintf("kubeadm init --token=%s --skip-phases=addon/kube-proxy --cri-socket unix:///var/run/containerd/containerd.sock --v=5", vars.JoinToken),
		"curl -Lsk -o /tmp/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64",
		"install /tmp/argocd /usr/bin",
	}...)
	var cniCmds []string
	switch vars.CNI {
	case "cilium":
		cniCmds = []string{
			"helm repo add cilium https://helm.cilium.io/",
			"helm repo update",
			"helm upgrade --wait --kubeconfig /etc/kubernetes/admin.conf --install --namespace kube-system cilium cilium/cilium -f /tmp/values-cilium.yaml --set k8sServiceHost=$(ip --json -4 a | jq -r '.[] | select(.ifname!=\"lo\") | .addr_info[0].local')",
		}
	}
	cmds = append(cmds, cniCmds...)
	cmds = append(cmds, []string{
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace argocd",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace gateway",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace cert-manager",
		"kubectl --kubeconfig /etc/kubernetes/admin.conf create namespace external-dns",
		"helm upgrade --wait --kubeconfig /etc/kubernetes/admin.conf --install argocd argo/argo-cd -f /tmp/values-argocd.yaml -n argocd",
	}...)
	cmds = append(cmds, argoCDApps(vars.ArgoCDApps)...)
	cmds = append(cmds, []string{
		fmt.Sprintf("kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=%s --from-literal=email=%s -n cert-manager", vars.CloudflareToken, vars.CloudflareEmail),
		fmt.Sprintf("kubectl --kubeconfig /etc/kubernetes/admin.conf create secret generic cloudflare-api-token --from-literal=token=%s --from-literal=email=%s -n external-dns", vars.CloudflareToken, vars.CloudflareEmail),
		"mkdir -p /home/supertylerc/.kube",
		"cp -i /etc/kubernetes/admin.conf /home/supertylerc/.kube/config",
		"chown supertylerc:supertylerc /home/supertylerc/.kube",
		"chown supertylerc:supertylerc /home/supertylerc/.kube/config",
		"while kubectl --kubeconfig /etc/kubernetes/admin.conf get application -A | grep -v 'Synced.*Healthy' | grep -v NAME; do sleep 0.5; done",
		"while kubectl --kubeconfig /etc/kubernetes/admin.conf get -A pods -o custom-columns=NAMESPACE:metadata.namespace,POD:metadata.name,PodIP:status.podIP,READY-true:status.containerStatuses[*].ready | grep -v true; do sleep 0.5; done",
	}...)
	switch vars.CNI {
	case "cilium":
		cmds = append(cmds, []string{
			"kubectl --kubeconfig /etc/kubernetes/admin.conf rollout restart deploy/cilium-operator -n kube-system",
			"kubectl --kubeconfig /etc/kubernetes/admin.conf rollout restart ds/cilium -n kube-system",
			"while kubectl --kubeconfig /etc/kubernetes/admin.conf get -A pods -o custom-columns=NAMESPACE:metadata.namespace,POD:metadata.name,PodIP:status.podIP,READY-true:status.containerStatuses[*].ready | grep -v true; do sleep 0.5; done",
		}...)
	}
	return cmds
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
