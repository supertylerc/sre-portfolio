resource "libvirt_volume" "k8s_lab_base" {
  name   = var.prefix
  source = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
  format = "qcow2"
}

resource "libvirt_volume" "k8s_lab_vm" {
  for_each = local.nodes

  name           = each.key
  base_volume_id = libvirt_volume.k8s_lab_base.id
  format         = "qcow2"
  size           = each.value.disk * 1024 * 1024 * 1024
}
