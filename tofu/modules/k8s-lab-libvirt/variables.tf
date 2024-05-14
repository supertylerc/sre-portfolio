variable "prefix" {
  type        = string
  description = "Prefix used for uniquely identifying  a stack."
}

variable "nodes" {
  type = list(object({
    kind   = string # One of: [control-plane, node]
    num    = number # Number of kind to create
    cpu    = number # Number of CPUs
    memory = number # Amount of RAM, in GB
    disk   = number # Size of disk, in GB
  }))
  description = "A list of node configurations."
}

variable "join_token" {
  type        = string
  sensitive   = true
  description = "Kubernetes join token created with 'kubeadm token generate' to simplify cluster joins"
}

variable "libvirt_cidr" {
  type        = string
  description = "CIDR for the libvirt routed network (must have a route on your network pointing at the libvirt node's IP)"
}

variable "control_plane_num" {
  type        = number
  description = "Which number in the subnet to use for the control plane"
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

variable "users" {
  type = list(object({
    name                = string
    hashed_passwd       = string
    lock_passwd         = optional(bool, false)
    groups              = optional(string, "users, admin")
    shell               = optional(string, "/bin/bash")
    ssh_authorized_keys = optional(list(string))
  }))
  description = "Users to create with cloud-init"
}
