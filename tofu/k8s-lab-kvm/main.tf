module "lab" {
  source = "../modules/k8s-lab-libvirt"

  prefix            = local.prefix
  nodes             = local.nodes
  libvirt_cidr      = local.libvirt_cidr
  control_plane_num = local.control_plane_num
  join_token        = var.join_token
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
}

variable "join_token" {
  type        = string
  sensitive   = true
  description = "kubeadm join token to simplify cluster provisioning (generate with 'kubeadm token generate')"
}
