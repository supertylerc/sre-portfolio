module "lab" {
  source = "../modules/k8s-lab-libvirt"

  prefix            = local.prefix
  nodes             = local.nodes
  libvirt_cidr      = local.libvirt_cidr
  control_plane_num = local.control_plane_num
  join_token        = var.join_token
  cloudflare_token  = var.cloudflare_token
  cloudflare_email  = var.cloudflare_email
  users             = local.users
  argocd_domain     = local.argocd_domain
  argocd_apps       = local.argocd_apps
}

locals {
  prefix = "k8s_lab"
  nodes = [
    {
      kind   = "control-plane"
      num    = 1
      cpu    = 2
      memory = 8
      disk   = 40
    },
    {
      kind   = "node"
      num    = 6
      cpu    = 2
      memory = 8
      disk   = 40
    }
  ]
  libvirt_cidr      = "10.0.42.0/24"
  control_plane_num = 10
  users = [{
    name                = "supertylerc"
    hashed_passwd       = var.hashed_password
    ssh_authorized_keys = [file("~/.ssh/id_ed25519.pub")]
  }]
  argocd_domain = "argocd.local.tylerc.me"
  argocd_apps = [
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/gateway-api.yaml",
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/argocd.yaml",
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/cilium.yaml",
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/cert-manager.yaml",
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/external-dns.yaml",
    "https://raw.githubusercontent.com/supertylerc/sre-portfolio/main/argo/applications/longhorn.yaml",
  ]

}

variable "join_token" {
  type        = string
  sensitive   = true
  description = "kubeadm join token to simplify cluster provisioning (generate with 'kubeadm token generate')"
}

variable "cloudflare_token" {
  type        = string
  sensitive   = true
  description = "Cloudflare API token for cert-manager and LE ClusterIssuer to use"
}

variable "cloudflare_email" {
  type        = string
  description = "E-mail address used for Cloudflare API token"
}

variable "hashed_password" {
  type        = string
  sensitive   = true
  description = "Hashed password for user"
}
