provider "go" {
  go = file("${path.module}/lib/expand_nodes.go")
}

provider "libvirt" {
  uri = "qemu:///system"
}
