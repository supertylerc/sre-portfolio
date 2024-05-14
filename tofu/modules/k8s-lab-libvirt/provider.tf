provider "go" {
  go = file("${path.module}/lib/lib.go")
}

provider "libvirt" {
  uri = "qemu:///system"
}
