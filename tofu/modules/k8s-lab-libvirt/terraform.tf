terraform {
  required_version = ">1.6"
  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = ">0.7"
    }
  }
}
