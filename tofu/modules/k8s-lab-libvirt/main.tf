resource "libvirt_domain" "k8s_lab" {
  for_each = {
    for i in flatten([
      for n in var.nodes : [
        for c in range(n.count) : {
          num    = c
          kind   = n.kind
          cpu    = n.cpu
          memory = n.memory
          disk   = n.disk
        }
      ]
    ]) : "${i.kind}-${i.num}" => i
  }
  name = each.key
  cpu {
    mode = "host-passthrough"
  }
  vcpu       = each.value.cpu
  memory     = each.value.memory * 1024
  autostart  = true
  qemu_agent = false
  cloudinit  = each.value.kind == "control-plane" ? libvirt_cloudinit_disk.control_plane.id : libvirt_cloudinit_disk.node.id
  graphics {
    type = "vnc"
  }
  video {
    type = "none"
  }

  console {
    type        = "pty"
    target_port = "0"
    target_type = "serial"
  }

  disk {
    volume_id = libvirt_volume.k8s_lab_vm[each.key].id
  }

  network_interface {
    hostname   = each.key
    addresses  = each.value.kind == "control-plane" ? [cidrhost(var.libvirt_cidr, var.control_plane_num)] : null
    network_id = libvirt_network.k8s_lab.id
  }
}
