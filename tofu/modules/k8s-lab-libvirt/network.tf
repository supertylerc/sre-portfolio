resource "libvirt_network" "k8s_lab" {
  name      = var.prefix
  mode      = "route"
  domain    = "k8s_lab"
  addresses = [var.libvirt_cidr]
  dhcp {
    enabled = true
  }
  dns {
    enabled    = true
    local_only = false
  }
  autostart = true
}
