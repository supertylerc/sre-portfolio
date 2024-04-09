source "digitalocean" "base" {
  image            = var.base_os_slug
  region           = var.region
  size             = var.size
  snapshot_regions = var.snapshot_regions
  ssh_username     = var.ssh_username
  ipv6             = true
  state_timeout    = "15m"
  ssh_timeout      = "15m"
  api_token        = var.api_token
}

# Consul
build {
  name = "consul"
  source "digitalocean.base" {
    droplet_name  = "consul-packer"
    snapshot_name = "consul-${var.base_os_slug}"
    tags          = concat(var.tags, ["consul"])
  }
  provisioner "shell" {
    environment_vars = [
      "CONSUL_VERSION=${var.consul_version}",
      "CONSUL_DC=${var.consul_dc}",
    ]
    scripts = [
      "scripts/common.sh",
      "scripts/consul.sh",
    ]
  }
}

# Nomad Server
build {
  name = "nomad-server"
  source "digitalocean.base" {
    droplet_name  = "nomad-server-packer"
    snapshot_name = "nomad-server-${var.base_os_slug}"
    tags          = concat(var.tags, ["nomad-server"])
  }
  provisioner "shell" {
    environment_vars = [
      "CONSUL_VERSION=${var.consul_version}",
      "CONSUL_DC=${var.consul_dc}",
      "NOMAD_VERSION=${var.nomad_version}",
      "NOMAD_PODMAN_VERSION=${var.nomad_podman_version}",
      "CNI_VERSION=${var.cni_version}",
    ]
    scripts = [
      "scripts/common.sh",
      "scripts/nomad-common.sh",
      "scripts/nomad-server.sh",
      "scripts/nomad-client.sh",
    ]
  }
}

# Nomad Client
build {
  name = "nomad-client"
  source "digitalocean.base" {
    droplet_name  = "nomad-client-packer"
    snapshot_name = "nomad-client-${var.base_os_slug}"
    tags          = concat(var.tags, ["nomad-client"])
  }
  provisioner "shell" {
    environment_vars = [
      "CONSUL_VERSION=${var.consul_version}",
      "CONSUL_DC=${var.consul_dc}",
      "NOMAD_VERSION=${var.nomad_version}",
      "NOMAD_PODMAN_VERSION=${var.nomad_podman_version}",
      "CNI_VERSION=${var.cni_version}",
    ]
    scripts = [
      "scripts/common.sh",
      "scripts/nomad-common.sh",
      "scripts/nomad-client.sh",
    ]
  }
}
