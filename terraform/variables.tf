variable "labels" {
  description = "Map of labels and values"
  type        = map(string)
  default = {
    name        = "lab"
    environment = "dev"
    managed_by  = "OpenTofu"
  }
}

variable "region" {
  description = "DigitalOcean Region in which to deploy the stack"
  type        = string
  default     = "tor1"
}

variable "spaces_region" {
  description = "DigitalOcean Region in which to deploy Spaces (not availble in all regions)"
  type        = string
  default     = "sfo3"
}

variable "spaces_name" {
  description = "Name of Space (globally unique)"
  type        = string
  default     = "supertylerc-tofu-lab"
}

variable "droplet_count" {
  description = "Number of droplets to create"
  type        = number
  default     = 1
}

variable "image" {
  description = "DigitalOcean Image Slug for Droplets"
  type        = string
  default     = "ubuntu-22-04-x64"
}

variable "ssh_key_name" {
  description = "Name in DigitalOcean for the SSH Key used to access Droplets"
  type        = string
  default     = "tyler-laptop"
}
variable "ssh_public_key_path" {
  description = "Path on local disk to public SSH key"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "outbound_rules" {
  description = "List of objects that represent the configuration of each outbound rule."
  type = list(object({
    protocol              = string
    port_range            = string
    destination_addresses = list(string)
  }))
  default = [
    {
      protocol              = "tcp"
      port_range            = "1-65535"
      destination_addresses = ["0.0.0.0/0"]
    },
    {
      protocol              = "udp"
      port_range            = "1-65535"
      destination_addresses = ["0.0.0.0/0"]
    }
  ]
}

variable "laptop_ports" {
  type        = list(string)
  description = "List of ports to allow from local laptop"
  default     = ["22"]
}
