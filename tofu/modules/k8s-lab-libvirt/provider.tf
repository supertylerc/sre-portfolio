provider "go" {
  go = file("${path.module}/lib.go")
}

provider "libvirt" {
  uri = "qemu:///system"
}
