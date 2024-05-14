resource "libvirt_cloudinit_disk" "control_plane" {
  name      = "${var.prefix}_control-plane_cloud-init.iso"
  user_data = join("\n", ["#cloud-config", yamlencode(merge(local.user_data, { runcmd = local.control_plane_run_cmds }))])
}

resource "libvirt_cloudinit_disk" "node" {
  name      = "${var.prefix}_node_cloud-init.iso"
  user_data = join("\n", ["#cloud-config", yamlencode(merge(local.user_data, { runcmd = local.node_run_cmds }))])
}

locals {
  node_run_cmds = provider::go::cloudruncmds("node", {
    control_plane_ip = cidrhost(var.libvirt_cidr, var.control_plane_num)
    join_token       = var.join_token
  })
  control_plane_run_cmds = provider::go::cloudruncmds("control-plane", {
    join_token       = var.join_token
    cloudflare_token = var.cloudflare_token
    cloudflare_email = var.cloudflare_email
  })
  user_data = {
    users = var.users
    apt = {
      sources = {
        "docker.list" = {
          source = "deb [arch=amd64] https://download.docker.com/linux/ubuntu $RELEASE stable"
          keyid  = "9DC858229FC7DD38854AE2D88D81803C0EBFCD88"
        }
        "kubernetes.list" = {
          source = "deb https://pkgs.k8s.io/core:/stable:/v1.30/deb/ /"
          keyid  = "DE15B14486CD377B9E876E1A234654DA9A296436"
        }
      }
    }
    package_update  = true
    package_upgrade = true
    packages = [
      "apt-transport-https",
      "ca-certificates",
      "curl",
      "gnupg",
      "lsb-release",
      "containerd.io",
      "cri-tools",
      "kubelet",
      "kubeadm",
      "kubectl",
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
          "net.bridge.bridge-nf-call-iptables=1",
          "net.ipv4.conf.all.proxy_arp=1",
          "net.ipv4.conf.ens2.proxy_arp=1"
        ])
      },
      {
        path    = "/etc/containerd/config.toml"
        content = file("${path.module}/configs/containerd.config.toml")
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
        content = file("${path.module}/configs/cilium.values.yaml")
      },
      {
        path    = "/tmp/values-argocd.yaml"
        content = templatefile("${path.module}/configs/argocd.values.yaml.tftpl", { argocd_domain = var.argocd_domain })
      },
    ]
    groups = ["docker"]
    power_state = {
      mode = "reboot"
    }
  }
}
