variable "base_os_slug" {
  description = "Image slug to use"
  type        = string
  default     = "fedora-38-x64"
}

variable "region" {
  description = "Region in which to build the image"
  type        = string
  default     = "tor1"
}

variable "size" {
  description = "Size slug for CPU/RAM"
  type        = string
  default     = "s-1vcpu-1gb"
}

variable "snapshot_regions" {
  description = "Region(s) in which the image will be available"
  type        = list(string)
  default = [
    "tor1",
  ]
}

variable "ssh_username" {
  description = "Username when SSHing for provisioning"
  type        = string
  default     = "root"
}

variable "tags" {
  description = "Tags to apply to the DigitalOcean instance"
  type        = list(string)
  default = [
    "dev",
    "packer",
    "fedora",
  ]
}

variable "nomad_version" {
  description = "Version of Nomad to install"
  type        = string
  default     = "1.5.5"
}

variable "consul_version" {
  description = "Version of Consul to install"
  type        = string
  default     = "1.15.2"
}

variable "consul_dc" {
  description = "DC to set when building the image (for development)"
  type        = string
  default     = "dc1"
}

variable "nomad_podman_version" {
  description = "Version of the Nomad Podman plugin to install"
  type        = string
  default     = "0.4.2"
}

variable "cni_version" {
  description = "Version of the reference CNI plugins to install"
  type        = string
  default     = "1.3.0"
}

variable "api_token" {
  description = "The DigitalOcean API Token"
  type        = string
}
