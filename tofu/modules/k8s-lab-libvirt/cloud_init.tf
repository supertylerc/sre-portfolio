resource "libvirt_cloudinit_disk" "control_plane" {
  name      = "${var.prefix}_control-plane_cloud-init.iso"
  user_data = join("\n", ["#cloud-config", yamlencode(merge(local.user_data, { runcmd = local.control_plane_run_cmds }))])
}

resource "libvirt_cloudinit_disk" "node" {
  name      = "${var.prefix}_node_cloud-init.iso"
  user_data = join("\n", ["#cloud-config", yamlencode(merge(local.user_data, { runcmd = local.node_run_cmds }))])
}

locals {
  common_run_cmds = [
    "hostnamectl set-hostname $(uuidgen)",
    "modprobe br_netfilter",
    "modprobe overlay",
    "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg", # Add Docker’s official GPG key
    "curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.30/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg",
    "echo \"deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null", # set up the stable docker repository
    "echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list",
    "apt-get update -y",
    "apt-get install -y containerd.io cri-tools kubelet kubeadm kubectl",
    "apt-mark hold kubelet kubeadm kubectl",
    "sysctl --system", # Reload settings from all system configuration files to take iptables configuration
    "systemctl restart containerd",
    "systemctl enable containerd",
  ]

  node_run_cmds = concat(local.common_run_cmds, [
    "until nc -zvw5 ${cidrhost(var.libvirt_cidr, var.control_plane_num)} 6443; do sleep 0.5; done",
    "kubeadm join ${cidrhost(var.libvirt_cidr, var.control_plane_num)}:6443 --token ${var.join_token} --discovery-token-unsafe-skip-ca-verification"
  ])
  control_plane_run_cmds = concat(local.common_run_cmds, [
    "snap install helm --classic",
    "helm repo add cilium https://helm.cilium.io/",
    "kubeadm init --token=${var.join_token} --skip-phases=addon/kube-proxy --cri-socket unix:///var/run/containerd/containerd.sock --v=5",
    "helm upgrade --kubeconfig /etc/kubernetes/admin.conf --install --namespace kube-system cilium cilium/cilium -f /tmp/values-cilium.yaml --set k8sServiceHost=$(ip --json -4 a | jq -r '.[] | select(.ifname!=\"lo\") | .addr_info[0].local')",
    "mkdir -p /home/supertylerc/.kube",
    "cp -i /etc/kubernetes/admin.conf /home/supertylerc/.kube/config",
    "chown supertylerc:supertylerc /home/supertylerc/.kube",
    "chown supertylerc:supertylerc /home/supertylerc/.kube/config"
  ])

  user_data = {
    users = [{
      name                = "supertylerc"
      hashed_passwd       = "$6$rounds=4096$nDzAXP.111c5arXO$58.b5gsevh.JHRNKzuw2BMd7P78VC19GYFOlzAOGIwOZUwWgxhaKq0HyWnq8GaKwfeDcF4cXIU.M4fyp7169U."
      lock_passwd         = false
      groups              = "users, admin"
      shell               = "/bin/bash"
      ssh_authorized_keys = [file("~/.ssh/id_ed25519.pub")]
    }]

    package_update  = true
    package_upgrade = true
    packages = [
      "apt-transport-https",
      "ca-certificates",
      "curl",
      "gnupg",
      "lsb-release",
    ]
    write_files = [
      {
        path = "/etc/modules-load.d/k8s.conf"
        content = join("\n", [
          "br_netfilter",
          "overlay"
        ])
      },
      {
        path = "/etc/sysctl.d/k8s.conf"
        content = join("\n", [
          "net.ipv4.ip_forward=1",
          "net.ipv6.conf.all.forwarding=1",
          "net.bridge.bridge-nf-call-ip6tables=1",
          "net.bridge.bridge-nf-call-iptables=1"
        ])
      },
      {
        path    = "/etc/containerd/config.toml"
        content = file("${path.module}/containerd.config.toml")
      },
      {
        path = "/etc/crictl.yaml"
        content = join("\n", [
          "runtime-endpoint: unix:///run/containerd/containerd.sock",
          "image-endpoint: unix:///run/containerd/containerd.sock",
          "timeout: 2",
          "debug: true # <- if you don't want to see debug info you can set this to false",
          "pull-image-on-create: false",
        ])
      },
      {
        path    = "/tmp/values-cilium.yaml"
        content = file("${path.module}/cilium.values.yaml")
      }
    ]
    groups = ["docker"]
    power_state = {
      mode = "reboot"
    }
  }
}
